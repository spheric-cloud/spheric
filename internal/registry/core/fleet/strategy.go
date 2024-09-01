// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package fleet

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
	fleet, ok := obj.(*core.Fleet)
	if !ok {
		return nil, nil, fmt.Errorf("given object is not a FleetRef")
	}
	return fleet.Labels, SelectableFields(fleet), nil
}

func MatchFleet(label labels.Selector, field fields.Selector) apisrvstorage.SelectionPredicate {
	return apisrvstorage.SelectionPredicate{
		Label:    label,
		Field:    field,
		GetAttrs: GetAttrs,
	}
}

func SelectableFields(fleet *core.Fleet) fields.Set {
	return generic.ObjectMetaFieldsSet(&fleet.ObjectMeta, false)
}

type fleetStrategy struct {
	runtime.ObjectTyper
	names.NameGenerator
}

var Strategy = fleetStrategy{api.Scheme, names.SimpleNameGenerator}

func (fleetStrategy) NamespaceScoped() bool {
	return false
}

func (fleetStrategy) GetResetFields() map[fieldpath.APIVersion]*fieldpath.Set {
	return map[fieldpath.APIVersion]*fieldpath.Set{
		"core.spheric.cloud/v1alpha1": fieldpath.NewSet(
			fieldpath.MakePathOrDie("status"),
		),
	}
}

func (fleetStrategy) PrepareForCreate(ctx context.Context, obj runtime.Object) {
	fleet := obj.(*core.Fleet)
	fleet.Status = core.FleetStatus{}
	fleet.Generation = 1
}

func (fleetStrategy) PrepareForUpdate(ctx context.Context, obj, old runtime.Object) {
	newFleet := obj.(*core.Fleet)
	oldFleet := old.(*core.Fleet)
	newFleet.Status = oldFleet.Status

	if !equality.Semantic.DeepEqual(newFleet.Spec, oldFleet.Spec) {
		newFleet.Generation = oldFleet.Generation + 1
	}
}

func (fleetStrategy) Validate(ctx context.Context, obj runtime.Object) field.ErrorList {
	fleet := obj.(*core.Fleet)
	return validation.ValidateFleet(fleet)
}

func (fleetStrategy) WarningsOnCreate(ctx context.Context, obj runtime.Object) []string {
	return nil
}

func (fleetStrategy) AllowCreateOnUpdate() bool {
	return false
}

func (fleetStrategy) AllowUnconditionalUpdate() bool {
	return false
}

func (fleetStrategy) Canonicalize(obj runtime.Object) {
}

func (fleetStrategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	newFleet := obj.(*core.Fleet)
	oldFleet := old.(*core.Fleet)
	return validation.ValidateFleetUpdate(newFleet, oldFleet)
}

func (fleetStrategy) WarningsOnUpdate(ctx context.Context, obj, old runtime.Object) []string {
	return nil
}

type fleetStatusStrategy struct {
	fleetStrategy
}

var StatusStrategy = fleetStatusStrategy{Strategy}

func (fleetStatusStrategy) GetResetFields() map[fieldpath.APIVersion]*fieldpath.Set {
	return map[fieldpath.APIVersion]*fieldpath.Set{
		"core.spheric.cloud/v1alpha1": fieldpath.NewSet(
			fieldpath.MakePathOrDie("spec"),
		),
	}
}

func (fleetStatusStrategy) PrepareForUpdate(ctx context.Context, obj, old runtime.Object) {
	newFleet := obj.(*core.Fleet)
	oldFleet := old.(*core.Fleet)
	newFleet.Spec = oldFleet.Spec
}

func (fleetStatusStrategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	newFleet := obj.(*core.Fleet)
	oldFleet := old.(*core.Fleet)
	return validation.ValidateFleetUpdate(newFleet, oldFleet)
}

func (fleetStatusStrategy) WarningsOnUpdate(cxt context.Context, obj, old runtime.Object) []string {
	return nil
}
