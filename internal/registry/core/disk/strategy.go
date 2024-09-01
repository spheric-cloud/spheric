// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package disk

import (
	"context"
	"fmt"

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
)

func GetAttrs(obj runtime.Object) (labels.Set, fields.Set, error) {
	disk, ok := obj.(*core.Disk)
	if !ok {
		return nil, nil, fmt.Errorf("given object is not a Disk")
	}
	return disk.Labels, SelectableFields(disk), nil
}

func MatchDisk(label labels.Selector, field fields.Selector) apisrvstorage.SelectionPredicate {
	return apisrvstorage.SelectionPredicate{
		Label:       label,
		Field:       field,
		GetAttrs:    GetAttrs,
		IndexFields: []string{},
	}
}

func SelectableFields(disk *core.Disk) fields.Set {
	fieldsSet := make(fields.Set)
	return generic.AddObjectMetaFieldsSet(fieldsSet, &disk.ObjectMeta, true)
}

func Indexers() *cache.Indexers {
	return &cache.Indexers{}
}

type diskStrategy struct {
	runtime.ObjectTyper
	names.NameGenerator
}

var Strategy = diskStrategy{api.Scheme, names.SimpleNameGenerator}

func (diskStrategy) NamespaceScoped() bool {
	return true
}

func (diskStrategy) PrepareForCreate(ctx context.Context, obj runtime.Object) {
}

func (diskStrategy) PrepareForUpdate(ctx context.Context, obj, old runtime.Object) {
}

func (diskStrategy) Validate(ctx context.Context, obj runtime.Object) field.ErrorList {
	disk := obj.(*core.Disk)
	return validation.ValidateDisk(disk)
}

func (diskStrategy) WarningsOnCreate(ctx context.Context, obj runtime.Object) []string {
	return nil
}

func (diskStrategy) AllowCreateOnUpdate() bool {
	return false
}

func (diskStrategy) AllowUnconditionalUpdate() bool {
	return false
}

func (diskStrategy) Canonicalize(obj runtime.Object) {
}

func (diskStrategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	newDisk, oldDisk := obj.(*core.Disk), old.(*core.Disk)
	return validation.ValidateDiskUpdate(newDisk, oldDisk)
}

func (diskStrategy) WarningsOnUpdate(ctx context.Context, obj, old runtime.Object) []string {
	return nil
}

type diskStatusStrategy struct {
	diskStrategy
}

var StatusStrategy = diskStatusStrategy{Strategy}

func (diskStatusStrategy) GetResetFields() map[fieldpath.APIVersion]*fieldpath.Set {
	return map[fieldpath.APIVersion]*fieldpath.Set{
		"core.spheric.cloud/v1alpha1": fieldpath.NewSet(
			fieldpath.MakePathOrDie("spec"),
		),
	}
}

func (diskStatusStrategy) PrepareForUpdate(ctx context.Context, obj, old runtime.Object) {
	newDisk := obj.(*core.Disk)
	oldDisk := old.(*core.Disk)
	newDisk.Spec = oldDisk.Spec
}

func (diskStatusStrategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	return nil
}

func (diskStatusStrategy) WarningsOnUpdate(cxt context.Context, obj, old runtime.Object) []string {
	return nil
}
