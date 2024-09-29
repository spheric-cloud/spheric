// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package instance

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/apiserver/pkg/registry/generic"
	apisrvstorage "k8s.io/apiserver/pkg/storage"
	"k8s.io/apiserver/pkg/storage/names"
	"k8s.io/client-go/tools/cache"
	"sigs.k8s.io/structured-merge-diff/v4/fieldpath"
	"spheric.cloud/spheric/internal/api"
	"spheric.cloud/spheric/internal/apis/core"
	"spheric.cloud/spheric/internal/apis/core/validation"
	"spheric.cloud/spheric/internal/spherelet/client"
)

func GetAttrs(obj runtime.Object) (labels.Set, fields.Set, error) {
	instance, ok := obj.(*core.Instance)
	if !ok {
		return nil, nil, fmt.Errorf("given object is not a Instance")
	}
	return instance.Labels, SelectableFields(instance), nil
}

func MatchInstance(label labels.Selector, field fields.Selector) apisrvstorage.SelectionPredicate {
	return apisrvstorage.SelectionPredicate{
		Label:       label,
		Field:       field,
		GetAttrs:    GetAttrs,
		IndexFields: []string{core.InstanceFleetRefNameField},
	}
}

func instanceFleetRefName(instance *core.Instance) string {
	if fleetRef := instance.Spec.FleetRef; fleetRef != nil {
		return fleetRef.Name
	}
	return ""
}

func instanceInstanceTypeRefName(instance *core.Instance) string {
	return instance.Spec.InstanceTypeRef.Name
}

func SelectableFields(instance *core.Instance) fields.Set {
	fieldsSet := make(fields.Set)
	fieldsSet[core.InstanceFleetRefNameField] = instanceFleetRefName(instance)
	fieldsSet[core.InstanceInstanceTypeRefNameField] = instanceInstanceTypeRefName(instance)
	return generic.AddObjectMetaFieldsSet(fieldsSet, &instance.ObjectMeta, true)
}

func FleetRefNameIndexFunc(obj any) ([]string, error) {
	instance, ok := obj.(*core.Instance)
	if !ok {
		return nil, fmt.Errorf("not a instance")
	}
	return []string{instanceFleetRefName(instance)}, nil
}

func FleetRefNameTriggerFunc(obj runtime.Object) string {
	return instanceFleetRefName(obj.(*core.Instance))
}

func Indexers() *cache.Indexers {
	return &cache.Indexers{
		apisrvstorage.FieldIndex(core.InstanceFleetRefNameField): FleetRefNameIndexFunc,
	}
}

type instanceStrategy struct {
	runtime.ObjectTyper
	names.NameGenerator
}

var Strategy = instanceStrategy{api.Scheme, names.SimpleNameGenerator}

func (instanceStrategy) NamespaceScoped() bool {
	return true
}

func (instanceStrategy) PrepareForCreate(ctx context.Context, obj runtime.Object) {
}

func (instanceStrategy) PrepareForUpdate(ctx context.Context, obj, old runtime.Object) {
}

func (instanceStrategy) Validate(ctx context.Context, obj runtime.Object) field.ErrorList {
	instance := obj.(*core.Instance)
	return validation.ValidateInstance(instance)
}

func (instanceStrategy) WarningsOnCreate(ctx context.Context, obj runtime.Object) []string {
	return nil
}

func (instanceStrategy) AllowCreateOnUpdate() bool {
	return false
}

func (instanceStrategy) AllowUnconditionalUpdate() bool {
	return false
}

func (instanceStrategy) Canonicalize(obj runtime.Object) {
}

func (instanceStrategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	oldInstance := old.(*core.Instance)
	newInstance := obj.(*core.Instance)
	return validation.ValidateInstanceUpdate(newInstance, oldInstance)
}

func (instanceStrategy) WarningsOnUpdate(ctx context.Context, obj, old runtime.Object) []string {
	return nil
}

type instanceStatusStrategy struct {
	instanceStrategy
}

var StatusStrategy = instanceStatusStrategy{Strategy}

func (instanceStatusStrategy) GetResetFields() map[fieldpath.APIVersion]*fieldpath.Set {
	return map[fieldpath.APIVersion]*fieldpath.Set{
		"core.spheric.cloud/v1alpha1": fieldpath.NewSet(
			fieldpath.MakePathOrDie("spec"),
		),
	}
}

func (instanceStatusStrategy) PrepareForUpdate(ctx context.Context, obj, old runtime.Object) {
	newInstance := obj.(*core.Instance)
	oldInstance := old.(*core.Instance)
	newInstance.Spec = oldInstance.Spec
}

func (instanceStatusStrategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	newInstance := obj.(*core.Instance)
	oldInstance := old.(*core.Instance)
	return validation.ValidateInstanceUpdate(newInstance, oldInstance)
}

func (instanceStatusStrategy) WarningsOnUpdate(cxt context.Context, obj, old runtime.Object) []string {
	return nil
}

type ResourceGetter interface {
	Get(ctx context.Context, name string, options *metav1.GetOptions) (runtime.Object, error)
}

func ExecLocation(
	ctx context.Context,
	getter ResourceGetter,
	connInfo client.ConnectionInfoGetter,
	name string,
	opts *core.InstanceExecOptions,
) (*url.URL, http.RoundTripper, error) {
	instance, err := getInstance(ctx, getter, name)
	if err != nil {
		return nil, nil, err
	}

	fleetRef := instance.Spec.FleetRef
	if fleetRef == nil {
		return nil, nil, apierrors.NewBadRequest(fmt.Sprintf("instance %s has no instance pool assigned", name))
	}

	fleetName := fleetRef.Name
	fleetInfo, err := connInfo.GetConnectionInfo(ctx, fleetName)
	if err != nil {
		return nil, nil, err
	}

	loc := &url.URL{
		Scheme: fleetInfo.Scheme,
		Host:   net.JoinHostPort(fleetInfo.Hostname, fleetInfo.Port),
		Path:   fmt.Sprintf("/apis/core.spheric.cloud/namespaces/%s/instances/%s/exec", instance.Namespace, instance.Name),
	}
	transport := fleetInfo.Transport
	if opts.InsecureSkipTLSVerifyBackend {
		transport = fleetInfo.InsecureSkipTLSVerifyTransport
	}

	return loc, transport, nil
}

func getInstance(ctx context.Context, getter ResourceGetter, name string) (*core.Instance, error) {
	obj, err := getter.Get(ctx, name, &metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	instance, ok := obj.(*core.Instance)
	if !ok {
		return nil, fmt.Errorf("unexpected object type %T", obj)
	}
	return instance, nil
}
