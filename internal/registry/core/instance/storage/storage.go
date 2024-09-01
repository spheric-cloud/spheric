// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package storage

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"spheric.cloud/spheric/internal/registry/core/instance"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/proxy"
	"k8s.io/apiserver/pkg/registry/generic"
	genericregistry "k8s.io/apiserver/pkg/registry/generic/registry"
	"k8s.io/apiserver/pkg/registry/rest"
	"k8s.io/apiserver/pkg/storage"
	"sigs.k8s.io/structured-merge-diff/v4/fieldpath"
	"spheric.cloud/spheric/internal/apis/core"
	"spheric.cloud/spheric/internal/spherelet/client"
)

type InstanceStorage struct {
	Instance *REST
	Status   *StatusREST
	Exec     *ExecREST
}

type REST struct {
	*genericregistry.Store
}

func NewStorage(optsGetter generic.RESTOptionsGetter, k client.ConnectionInfoGetter) (InstanceStorage, error) {
	store := &genericregistry.Store{
		NewFunc: func() runtime.Object {
			return &core.Instance{}
		},
		NewListFunc: func() runtime.Object {
			return &core.InstanceList{}
		},
		PredicateFunc:             instance.MatchInstance,
		DefaultQualifiedResource:  core.Resource("instances"),
		SingularQualifiedResource: core.Resource("instance"),

		CreateStrategy: instance.Strategy,
		UpdateStrategy: instance.Strategy,
		DeleteStrategy: instance.Strategy,

		TableConvertor: newTableConvertor(),
	}

	options := &generic.StoreOptions{
		RESTOptions: optsGetter,
		AttrFunc:    instance.GetAttrs,
		TriggerFunc: map[string]storage.IndexerFunc{
			core.InstanceFleetRefNameField: instance.FleetRefNameTriggerFunc,
		},
		Indexers: instance.Indexers(),
	}
	if err := store.CompleteWithOptions(options); err != nil {
		return InstanceStorage{}, err
	}

	statusStore := *store
	statusStore.UpdateStrategy = instance.StatusStrategy
	statusStore.ResetFieldsStrategy = instance.StatusStrategy

	return InstanceStorage{
		Instance: &REST{store},
		Status:   &StatusREST{&statusStore},
		Exec:     &ExecREST{store, k},
	}, nil
}

type StatusREST struct {
	store *genericregistry.Store
}

func (r *StatusREST) New() runtime.Object {
	return &core.Instance{}
}

func (r *StatusREST) Get(ctx context.Context, name string, options *metav1.GetOptions) (runtime.Object, error) {
	return r.store.Get(ctx, name, options)
}

func (r *StatusREST) Update(ctx context.Context, name string, objInfo rest.UpdatedObjectInfo, createValidation rest.ValidateObjectFunc, updateValidation rest.ValidateObjectUpdateFunc, forceAllowCreate bool, options *metav1.UpdateOptions) (runtime.Object, bool, error) {
	return r.store.Update(ctx, name, objInfo, createValidation, updateValidation, false, options)
}

func (r *StatusREST) GetResetFields() map[fieldpath.APIVersion]*fieldpath.Set {
	return r.store.GetResetFields()
}

func (r *StatusREST) Destroy() {}

// Support both GET and POST methods. We must support GET for browsers that want to use WebSockets.
var upgradeableMethods = []string{"GET", "POST"}

type ExecREST struct {
	Store        *genericregistry.Store
	InstanceConn client.ConnectionInfoGetter
}

func (r *ExecREST) New() runtime.Object {
	return &core.InstanceExecOptions{}
}

func (r *ExecREST) Connect(ctx context.Context, name string, opts runtime.Object, responder rest.Responder) (http.Handler, error) {
	execOpts, ok := opts.(*core.InstanceExecOptions)
	if !ok {
		return nil, fmt.Errorf("invalid options objects: %#v", opts)
	}

	location, transport, err := instance.ExecLocation(ctx, r.Store, r.InstanceConn, name, execOpts)
	if err != nil {
		return nil, err
	}

	return newThrottledUpgradeAwareProxyHandler(location, transport, false, true, responder), nil
}

func newThrottledUpgradeAwareProxyHandler(location *url.URL, transport http.RoundTripper, wrapTransport, upgradeRequired bool, responder rest.Responder) *proxy.UpgradeAwareHandler {
	handler := proxy.NewUpgradeAwareHandler(location, transport, wrapTransport, upgradeRequired, proxy.NewErrorResponder(responder))
	return handler
}

func (r *ExecREST) NewConnectOptions() (runtime.Object, bool, string) {
	return &core.InstanceExecOptions{}, false, ""
}

func (r *ExecREST) ConnectMethods() []string {
	return upgradeableMethods
}

func (r *ExecREST) Destroy() {}
