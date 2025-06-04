// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0

package etcd

import (
	"context"
	"fmt"
	"strings"

	clientv3 "go.etcd.io/etcd/client/v3"
	"spheric.cloud/spheric/actuo/codec"
	"spheric.cloud/spheric/actuo/list"
	"spheric.cloud/spheric/actuo/runtime"
	"spheric.cloud/spheric/actuo/storage/store"
	"spheric.cloud/spheric/actuo/watch"
	"spheric.cloud/spheric/utils/generic"
)

func notFound(key string) clientv3.Cmp {
	return clientv3.Compare(clientv3.ModRevision(key), "=", 0)
}

type Simple[Object interface {
	runtime.Object
	*ObjectVal
}, ObjectVal any] = Etcd[string, Object, *list.List[Object, ObjectVal]]

func NewSimple[Object interface {
	runtime.Object
	*ObjectVal
}, ObjectVal any](
	client *clientv3.Client,
	codec codec.Codec[Object],
	factory store.Factory[Object, *list.List[Object, ObjectVal]],
	versioner store.Versioner[Object, *list.List[Object, ObjectVal]],
) *Simple[Object, ObjectVal] {
	return New(client, generic.Identity, codec, factory, versioner)
}

type Etcd[Key, Object any, ObjectList runtime.List[Object]] struct {
	client *clientv3.Client

	keyFunc   func(Key) string
	codec     codec.Codec[Object]
	factory   store.Factory[Object, ObjectList]
	versioner store.Versioner[Object, ObjectList]
}

func _[Key, Object any, ObjectList runtime.List[Object]]() store.Store[Key, Object] {
	return generic.Stub[*Etcd[Key, Object, ObjectList]]()
}

func New[Key, Object runtime.Object, ObjectList runtime.List[Object]](
	client *clientv3.Client,
	keyFunc func(Key) string,
	codec codec.Codec[Object],
	factory store.Factory[Object, ObjectList],
	versioner store.Versioner[Object, ObjectList],
) *Etcd[Key, Object, ObjectList] {
	return &Etcd[Key, Object, ObjectList]{
		client:    client,
		keyFunc:   keyFunc,
		codec:     codec,
		factory:   factory,
		versioner: versioner,
	}
}

func (e *Etcd[Key, ObjectVal, ObjectList]) prepareKey(key Key) (string, error) {
	preparedKey := e.keyFunc(key)
	return preparedKey, nil
}

func (e *Etcd[Key, ObjectVal, ObjectList]) Create(ctx context.Context, k Key, obj ObjectVal) (ObjectVal, error) {
	preparedKey, err := e.prepareKey(k)
	if err != nil {
		return generic.Zero[ObjectVal](), err
	}

	if err := e.versioner.PrepareObjectForStorage(obj); err != nil {
		return generic.Zero[ObjectVal](), err
	}

	data, err := codec.Encode(e.codec, obj)
	if err != nil {
		return generic.Zero[ObjectVal](), err
	}

	txnResp, err := e.client.KV.Txn(ctx).If(
		notFound(preparedKey),
	).Then(
		clientv3.OpPut(preparedKey, string(data)),
	).Commit()
	if err != nil {
		return generic.Zero[ObjectVal](), err
	}

	if !txnResp.Succeeded {
		return generic.Zero[ObjectVal](), fmt.Errorf("%v %w", k, store.ErrAlreadyExists)
	}

	putResp := txnResp.Responses[0].GetResponsePut()
	return e.decodeNew(data, putResp.Header.Revision)
}

func (e *Etcd[Key, ObjectVal, ObjectList]) Get(ctx context.Context, key Key) (ObjectVal, error) {
	preparedKey, err := e.prepareKey(key)
	if err != nil {
		return generic.Zero[ObjectVal](), err
	}

	getResp, err := e.client.Get(ctx, preparedKey)
	if err != nil {
		return generic.Zero[ObjectVal](), err
	}

	if len(getResp.Kvs) == 0 {
		return generic.Zero[ObjectVal](), fmt.Errorf("%v %w", key, store.ErrNotFound)
	}

	kv := getResp.Kvs[0]
	return e.decodeNew(kv.Value, kv.ModRevision)
}

func (e *Etcd[Key, ObjectVal, ObjectList]) getState(getResp *clientv3.GetResponse, key string, ignoreNotFound bool) (ObjectVal, int64, error) {
	if len(getResp.Kvs) == 0 {
		if !ignoreNotFound {
			return generic.Zero[ObjectVal](), 0, fmt.Errorf("%v %w", key, store.ErrNotFound)
		}
		obj := e.factory.New()
		return obj, 0, nil
	}

	kv := getResp.Kvs[0]
	obj := e.factory.New()
	if err := codec.Decode(e.codec, kv.Value, obj); err != nil {
		return generic.Zero[ObjectVal](), 0, err
	}

	return obj, kv.ModRevision, nil
}

func (e *Etcd[Key, ObjectVal, ObjectList]) getCurrentState(ctx context.Context, key string, ignoreNotFound bool) (ObjectVal, int64, error) {
	res, err := e.client.Get(ctx, key)
	if err != nil {
		return generic.Zero[ObjectVal](), 0, err
	}

	return e.getState(res, key, ignoreNotFound)
}

func (e *Etcd[Key, ObjectVal, ObjectList]) decode(data []byte, into ObjectVal, revision int64) error {
	if err := codec.Decode(e.codec, data, into); err != nil {
		return err
	}
	if err := e.versioner.UpdateObject(into, uint64(revision)); err != nil {
		return err
	}
	return nil
}

func (e *Etcd[Key, ObjectVal, ObjectList]) decodeNew(data []byte, revision int64) (ObjectVal, error) {
	into := e.factory.New()
	if err := e.decode(data, into, revision); err != nil {
		return generic.Zero[ObjectVal](), err
	}
	return into, nil
}

func (e *Etcd[Key, ObjectVal, ObjectList]) Update(ctx context.Context, k Key, ignoreNotFound bool, update store.Update[ObjectVal]) (ObjectVal, error) {
	preparedKey, err := e.prepareKey(k)
	if err != nil {
		return generic.Zero[ObjectVal](), err
	}

	obj, rev, err := e.getCurrentState(ctx, preparedKey, ignoreNotFound)
	if err != nil {
		return generic.Zero[ObjectVal](), err
	}

	for {
		updated, err := update(ctx, obj)
		if err != nil {
			return generic.Zero[ObjectVal](), err
		}

		if err := e.versioner.PrepareObjectForStorage(updated); err != nil {
			return generic.Zero[ObjectVal](), err
		}

		data, err := codec.Encode(e.codec, updated)
		if err != nil {
			return generic.Zero[ObjectVal](), err
		}

		txnResp, err := e.client.Txn(ctx).If(
			clientv3.Compare(clientv3.ModRevision(preparedKey), "=", rev),
		).Then(
			clientv3.OpPut(preparedKey, string(data)),
		).Else(
			clientv3.OpGet(preparedKey),
		).Commit()
		if err != nil {
			return generic.Zero[ObjectVal](), err
		}
		if !txnResp.Succeeded {
			getResp := (*clientv3.GetResponse)(txnResp.Responses[0].GetResponseRange())
			obj, rev, err = e.getState(getResp, preparedKey, ignoreNotFound)
			if err != nil {
				return generic.Zero[ObjectVal](), err
			}
			continue
		}

		putResp := txnResp.Responses[0].GetResponsePut()
		return e.decodeNew(data, putResp.Header.Revision)
	}
}

func (e *Etcd[Key, ObjectVal, ObjectList]) Delete(ctx context.Context, k Key, del store.Delete[ObjectVal]) (ObjectVal, error) {
	preparedKey, err := e.prepareKey(k)
	if err != nil {
		return generic.Zero[ObjectVal](), err
	}

	obj, rev, err := e.getCurrentState(ctx, preparedKey, false)
	if err != nil {
		return generic.Zero[ObjectVal](), err
	}

	for {
		if err := del(ctx, obj); err != nil {
			return generic.Zero[ObjectVal](), err
		}

		txnResp, err := e.client.KV.Txn(ctx).If(
			clientv3.Compare(clientv3.ModRevision(preparedKey), "=", rev),
		).Then(
			clientv3.OpDelete(preparedKey),
		).Else(
			clientv3.OpGet(preparedKey),
		).Commit()
		if err != nil {
			return generic.Zero[ObjectVal](), err
		}
		if !txnResp.Succeeded {
			getResp := (*clientv3.GetResponse)(txnResp.Responses[0].GetResponseRange())
			obj, rev, err = e.getState(getResp, preparedKey, false)
			if err != nil {
				return generic.Zero[ObjectVal](), err
			}
			continue
		}

		return obj, nil
	}
}

func (e *Etcd[Key, ObjectVal, ObjectList]) List(ctx context.Context, k Key) (runtime.List[ObjectVal], error) {
	preparedKey, err := e.prepareKey(k)
	if err != nil {
		return nil, err
	}

	if !strings.HasSuffix(preparedKey, "/") {
		preparedKey += "/"
	}

	var options []clientv3.OpOption
	rangeEnd := clientv3.GetPrefixRangeEnd(preparedKey)
	options = append(options, clientv3.WithRange(rangeEnd))

	getResp, err := e.client.Get(ctx, preparedKey, options...)
	if err != nil {
		return nil, err
	}

	res := e.factory.NewList(len(getResp.Kvs))
	for i, kv := range getResp.Kvs {
		if err := e.decode(kv.Value, res.Item(i), kv.ModRevision); err != nil {
			return nil, err
		}
	}
	if err := e.versioner.UpdateList(res, uint64(getResp.Header.Revision)); err != nil {
		return nil, err
	}
	return res, nil
}

func (e *Etcd[Key, ObjectVal, ObjectList]) Watch(ctx context.Context, k Key) (watch.Watch[ObjectVal], error) {
	preparedKey, err := e.prepareKey(k)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(ctx)
	w := &watcher[Key, ObjectVal, ObjectList]{
		Etcd:        e,
		preparedKey: preparedKey,
		out:         make(chan watch.Event[ObjectVal]),
		cancel:      cancel,
	}
	go w.run(ctx)
	return w, nil
}

type watcher[Key, Object any, ObjectList runtime.List[Object]] struct {
	*Etcd[Key, Object, ObjectList]
	preparedKey string
	out         chan watch.Event[Object]
	cancel      context.CancelFunc
}

func (w *watcher[Key, Object, ObjectList]) run(ctx context.Context) {
	defer w.cancel()
	defer close(w.out)

	opts := []clientv3.OpOption{
		clientv3.WithRev(1),
		clientv3.WithPrevKV(),
		clientv3.WithPrefix(),
	}
	wch := w.client.Watch(ctx, w.preparedKey, opts...)
	for wres := range wch {
		if err := wres.Err(); err != nil {
			// TODO: Enhance / fix
			return
		}

		if wres.IsProgressNotify() {
			// TODO: Handle
			continue
		}

		for _, e := range wres.Events {
			evt, err := w.toEvent(e)
			if err != nil {
				// TODO: Handle
				continue
			}
			if evt == nil {
				continue
			}

			select {
			case <-ctx.Done():
				return
			case w.out <- *evt:
			}
		}
	}
}

func (w *watcher[Key, Object, ObjectList]) toEvent(evt *clientv3.Event) (*watch.Event[Object], error) {
	switch evt.Type {
	case clientv3.EventTypePut:
		obj, err := w.decodeNew(evt.Kv.Value, evt.Kv.ModRevision)
		if err != nil {
			return nil, err
		}

		typ := watch.EventTypeUpdated
		if evt.IsCreate() {
			typ = watch.EventTypeCreated
		}

		return &watch.Event[Object]{
			Type:   typ,
			Object: obj,
		}, nil
	case clientv3.EventTypeDelete:
		obj, err := w.decodeNew(evt.PrevKv.Value, evt.PrevKv.ModRevision)
		if err != nil {
			return nil, err
		}

		return &watch.Event[Object]{
			Type:   watch.EventTypeDeleted,
			Object: obj,
		}, nil
	default:
		return nil, nil
	}
}

func (w *watcher[Key, Object, ObjectList]) Events() <-chan watch.Event[Object] {
	return w.out
}

func (w *watcher[Key, Object, ObjectList]) Stop() {
	w.cancel()
}
