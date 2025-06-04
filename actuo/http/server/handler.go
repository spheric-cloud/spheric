// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0

package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"golang.org/x/net/websocket"
	"k8s.io/apimachinery/pkg/util/httpstream"
	"k8s.io/apimachinery/pkg/util/httpstream/wsstream"
	"spheric.cloud/spheric/actuo/meta"
	"spheric.cloud/spheric/actuo/resource"
	"spheric.cloud/spheric/actuo/types"
	"spheric.cloud/spheric/actuo/watch"
)

type RequestKeyer[Key any] interface {
	RequestKey(ctx context.Context, req *http.Request) (context.Context, Key)
	Request(ctx context.Context, req *http.Request) context.Context
}

type namespacedNameRequestKeyer struct {
	namespaceParam string
	nameParam      string
}

func (n *namespacedNameRequestKeyer) RequestKey(ctx context.Context, req *http.Request) (context.Context, types.NamespacedName) {
	namespace := req.PathValue(n.namespaceParam)
	name := req.PathValue(n.nameParam)
	return types.WithNamespace(ctx, namespace), types.NamespacedName{Namespace: namespace, Name: name}
}

func (n *namespacedNameRequestKeyer) Request(ctx context.Context, req *http.Request) context.Context {
	namespace := req.PathValue(n.namespaceParam)
	return types.WithNamespace(ctx, namespace)
}

const (
	DefaultNamespaceParam = "namespace"
	DefaultNameParam      = "name"

	DefaultNamespacePath = "/namespaces/{" + DefaultNamespaceParam + "}"
	DefaultItemPath      = "/{" + DefaultNameParam + "}"
)

func NewNamespacedNameRequestKeyer(namespaceParam, nameParam string) RequestKeyer[types.NamespacedName] {
	return &namespacedNameRequestKeyer{
		namespaceParam: namespaceParam,
		nameParam:      nameParam,
	}
}

var DefaultNamespacedNameRequestKeyer = NewNamespacedNameRequestKeyer(DefaultNamespaceParam, DefaultNameParam)

type ObjectFactory[Object any] interface {
	New() Object
}

type objectValFactory[Object interface{ *ObjectVal }, ObjectVal any] struct{}

func (objectValFactory[Object, ObjectVal]) New() Object {
	return Object(new(ObjectVal))
}

func ObjectValFactory[Object interface{ *ObjectVal }, ObjectVal any]() ObjectFactory[Object] {
	return objectValFactory[Object, ObjectVal]{}
}

type Middleware func(http.Handler) http.Handler
type PatternMiddleware func(pattern string) string

type RegisterCreateHandlerOptions struct {
	Middlewares    []Middleware
	PathMiddleware []PatternMiddleware
}

func applyMiddlewares(handler http.Handler, middlewares []Middleware) http.Handler {
	for _, m := range middlewares {
		handler = m(handler)
	}
	return handler
}

func applyPathMiddlewares(pattern string, middlewares []PatternMiddleware) string {
	for _, m := range middlewares {
		pattern = m(pattern)
	}
	return pattern
}

func (o *RegisterCreateHandlerOptions) ApplyOptions(opts []RegisterCreateHandlerOption) *RegisterCreateHandlerOptions {
	for _, opt := range opts {
		opt.ApplyToRegisterCreateHandler(o)
	}
	return o
}

type RegisterCreateHandlerOption interface {
	ApplyToRegisterCreateHandler(*RegisterCreateHandlerOptions)
}

func RegisterCreateHandler(mux *http.ServeMux, prefix string, handler http.Handler, opts ...RegisterCreateHandlerOption) {
	o := (&RegisterCreateHandlerOptions{}).ApplyOptions(opts)
	path := applyPathMiddlewares(prefix, o.PathMiddleware)
	handler = applyMiddlewares(handler, o.Middlewares)
	mux.Handle(fmt.Sprintf("POST %s", path), handler)
}

func CreateHandler[Object any](objectFactory ObjectFactory[Object], creater resource.Creater[Object]) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		obj := objectFactory.New()
		if err := json.NewDecoder(req.Body).Decode(obj); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		created, err := creater.Create(ctx, obj)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(created)
	})
}

func GetHandler[Key, Object any](keyer RequestKeyer[Key], getter resource.Getter[Key, Object]) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		ctx, key := keyer.RequestKey(ctx, req)

		obj, err := getter.Get(ctx, key)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(obj)
	})
}

func ListHandler[Key, Object any](
	keyer RequestKeyer[Key],
	lister resource.Lister[Object],
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		ctx = keyer.Request(ctx, req)

		list, err := lister.List(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(list)
	})
}

func UpdateHandler[Key, Object any](
	requestKeyer RequestKeyer[Key],
	objectFactory ObjectFactory[Object],
	updater resource.Updater[Key, Object],
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		ctx, key := requestKeyer.RequestKey(ctx, req)

		obj := objectFactory.New()
		updated, err := updater.Update(ctx, key, obj)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(updated)
	})
}

func DeleteHandler[Key, Object any](keyer RequestKeyer[Key], deleter resource.Deleter[Key, Object]) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		ctx, key := keyer.RequestKey(ctx, req)

		deleted, err := deleter.Delete(ctx, key)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(deleted)
	})
}

func WatchHandler[Key, Object any](keyer RequestKeyer[Key], watcher resource.Watcher[Object]) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		ctx = keyer.Request(ctx, req)
		req = req.WithContext(ctx)

		wt, err := watcher.Watch(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if IsWebSocketRequest(req) {
			serveWebsocketWatch(w, req, wt)
		} else {
			serveHTTPWatch(w, req, wt)
		}
	})
}

// IsWebSocketRequest returns true if the incoming request contains connection upgrade headers
// for WebSockets.
func IsWebSocketRequest(req *http.Request) bool {
	if !strings.EqualFold(req.Header.Get("Upgrade"), "websocket") {
		return false
	}
	return httpstream.IsUpgradeRequest(req)
}

func writeWatchEvent[Object any](w io.Writer, evt watch.Event[Object]) error {
	data, err := json.Marshal(&meta.WatchEvent{
		Type:   string(evt.Type),
		Object: evt.Object,
	})
	if err != nil {
		return err
	}

	if _, err := w.Write(data); err != nil {
		return err
	}
	return nil
}

func serveWebsocketWatch[Object any](
	w http.ResponseWriter,
	req *http.Request,
	wt watch.Watch[Object],
) {
	w.Header().Set("Content-Type", "application/json")
	websocket.Handler(func(conn *websocket.Conn) {
		defer wt.Stop()

		ctx, cancel := signalWebsocketClosed(context.Background(), conn)
		defer cancel()

		for {
			select {
			case <-ctx.Done():
				return
			case evt, ok := <-wt.Events():
				if !ok {
					return
				}

				if err := writeWatchEvent(conn, evt); err != nil {
					return
				}
			}
		}
	}).ServeHTTP(w, req)
}

func signalWebsocketClosed(ctx context.Context, conn *websocket.Conn) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	go func() {
		defer cancel()
		wsstream.IgnoreReceives(conn, 0)
	}()
	return ctx, cancel
}

func serveHTTPWatch[Object any](
	w http.ResponseWriter,
	req *http.Request,
	wt watch.Watch[Object],
) {
	defer wt.Stop()

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Transfer-Encoding", "chunked")
	w.WriteHeader(http.StatusOK)
	flusher.Flush()

	for {
		select {
		case <-req.Context().Done():
			return
		case evt, ok := <-wt.Events():
			if !ok {
				return
			}

			if err := writeWatchEvent(w, evt); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if len(wt.Events()) == 0 {
				flusher.Flush()
			}
		}
	}
}
