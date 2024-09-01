// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package network

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/api/equality"

	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/apiserver/pkg/registry/generic"
	apisrvstorage "k8s.io/apiserver/pkg/storage"
	"k8s.io/apiserver/pkg/storage/names"
	"sigs.k8s.io/structured-merge-diff/v4/fieldpath"
	"spheric.cloud/spheric/internal/api"
	"spheric.cloud/spheric/internal/apis/core"
	"spheric.cloud/spheric/internal/apis/core/validation"
)

func GetAttrs(obj runtime.Object) (labels.Set, fields.Set, error) {
	network, ok := obj.(*core.Network)
	if !ok {
		return nil, nil, fmt.Errorf("given object is not a Network")
	}
	return network.Labels, SelectableFields(network), nil
}

func MatchNetwork(label labels.Selector, field fields.Selector) apisrvstorage.SelectionPredicate {
	return apisrvstorage.SelectionPredicate{
		Label:    label,
		Field:    field,
		GetAttrs: GetAttrs,
	}
}

func SelectableFields(network *core.Network) fields.Set {
	return generic.ObjectMetaFieldsSet(&network.ObjectMeta, true)
}

type networkStrategy struct {
	runtime.ObjectTyper
	names.NameGenerator
}

var Strategy = networkStrategy{api.Scheme, names.SimpleNameGenerator}

func (networkStrategy) NamespaceScoped() bool {
	return true
}

func (networkStrategy) PrepareForCreate(ctx context.Context, obj runtime.Object) {
	machinePool := obj.(*core.Network)
	machinePool.Status = core.NetworkStatus{}
	machinePool.Generation = 1
}

func (networkStrategy) PrepareForUpdate(ctx context.Context, obj, old runtime.Object) {
	newNetwork, oldNetwork := obj.(*core.Network), old.(*core.Network)
	newNetwork.Status = oldNetwork.Status

	if !equality.Semantic.DeepEqual(newNetwork.Spec, oldNetwork.Spec) {
		newNetwork.Generation = oldNetwork.Generation + 1
	}
}

func (networkStrategy) Validate(ctx context.Context, obj runtime.Object) field.ErrorList {
	network := obj.(*core.Network)
	return validation.ValidateNetwork(network)
}

func (networkStrategy) WarningsOnCreate(ctx context.Context, obj runtime.Object) []string {
	return nil
}

func (networkStrategy) AllowCreateOnUpdate() bool {
	return false
}

func (networkStrategy) AllowUnconditionalUpdate() bool {
	return false
}

func (networkStrategy) Canonicalize(obj runtime.Object) {
}

func (networkStrategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	newNetwork, oldNetwork := obj.(*core.Network), old.(*core.Network)
	return validation.ValidateNetworkUpdate(newNetwork, oldNetwork)
}

func (networkStrategy) WarningsOnUpdate(ctx context.Context, obj, old runtime.Object) []string {
	return nil
}

type networkStatusStrategy struct {
	networkStrategy
}

var StatusStrategy = networkStatusStrategy{Strategy}

func (networkStatusStrategy) GetResetFields() map[fieldpath.APIVersion]*fieldpath.Set {
	return map[fieldpath.APIVersion]*fieldpath.Set{
		"core.spheric.cloud/v1alpha1": fieldpath.NewSet(
			fieldpath.MakePathOrDie("spec"),
		),
	}
}

func (networkStatusStrategy) PrepareForUpdate(ctx context.Context, obj, old runtime.Object) {
	newNetwork, oldNetwork := obj.(*core.Network), old.(*core.Network)
	newNetwork.Spec = oldNetwork.Spec
}

func (networkStatusStrategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	newNetwork := obj.(*core.Network)
	oldNetwork := old.(*core.Network)
	return validation.ValidateNetworkUpdate(newNetwork, oldNetwork)
}

func (networkStatusStrategy) WarningsOnUpdate(cxt context.Context, obj, old runtime.Object) []string {
	return nil
}
