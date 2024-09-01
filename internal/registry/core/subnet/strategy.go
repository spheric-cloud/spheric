// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package subnet

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
	subnet, ok := obj.(*core.Subnet)
	if !ok {
		return nil, nil, fmt.Errorf("given object is not a Subnet")
	}
	return subnet.Labels, SelectableFields(subnet), nil
}

func MatchSubnet(label labels.Selector, field fields.Selector) apisrvstorage.SelectionPredicate {
	return apisrvstorage.SelectionPredicate{
		Label:    label,
		Field:    field,
		GetAttrs: GetAttrs,
	}
}

func SelectableFields(subnet *core.Subnet) fields.Set {
	return generic.ObjectMetaFieldsSet(&subnet.ObjectMeta, true)
}

type subnetStrategy struct {
	runtime.ObjectTyper
	names.NameGenerator
}

var Strategy = subnetStrategy{api.Scheme, names.SimpleNameGenerator}

func (subnetStrategy) NamespaceScoped() bool {
	return true
}

func (subnetStrategy) PrepareForCreate(ctx context.Context, obj runtime.Object) {
	machinePool := obj.(*core.Subnet)
	machinePool.Status = core.SubnetStatus{}
	machinePool.Generation = 1
}

func (subnetStrategy) PrepareForUpdate(ctx context.Context, obj, old runtime.Object) {
	newSubnet, oldSubnet := obj.(*core.Subnet), old.(*core.Subnet)
	newSubnet.Status = oldSubnet.Status

	if !equality.Semantic.DeepEqual(newSubnet.Spec, oldSubnet.Spec) {
		newSubnet.Generation = oldSubnet.Generation + 1
	}
}

func (subnetStrategy) Validate(ctx context.Context, obj runtime.Object) field.ErrorList {
	subnet := obj.(*core.Subnet)
	return validation.ValidateSubnet(subnet)
}

func (subnetStrategy) WarningsOnCreate(ctx context.Context, obj runtime.Object) []string {
	return nil
}

func (subnetStrategy) AllowCreateOnUpdate() bool {
	return false
}

func (subnetStrategy) AllowUnconditionalUpdate() bool {
	return false
}

func (subnetStrategy) Canonicalize(obj runtime.Object) {
}

func (subnetStrategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	newSubnet, oldSubnet := obj.(*core.Subnet), old.(*core.Subnet)
	return validation.ValidateSubnetUpdate(newSubnet, oldSubnet)
}

func (subnetStrategy) WarningsOnUpdate(ctx context.Context, obj, old runtime.Object) []string {
	return nil
}

type subnetStatusStrategy struct {
	subnetStrategy
}

var StatusStrategy = subnetStatusStrategy{Strategy}

func (subnetStatusStrategy) GetResetFields() map[fieldpath.APIVersion]*fieldpath.Set {
	return map[fieldpath.APIVersion]*fieldpath.Set{
		"core.spheric.cloud/v1alpha1": fieldpath.NewSet(
			fieldpath.MakePathOrDie("spec"),
		),
	}
}

func (subnetStatusStrategy) PrepareForUpdate(ctx context.Context, obj, old runtime.Object) {
	newSubnet, oldSubnet := obj.(*core.Subnet), old.(*core.Subnet)
	newSubnet.Spec = oldSubnet.Spec
}

func (subnetStatusStrategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	newSubnet := obj.(*core.Subnet)
	oldSubnet := old.(*core.Subnet)
	return validation.ValidateSubnetUpdate(newSubnet, oldSubnet)
}

func (subnetStatusStrategy) WarningsOnUpdate(cxt context.Context, obj, old runtime.Object) []string {
	return nil
}
