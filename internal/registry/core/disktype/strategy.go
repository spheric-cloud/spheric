// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package disktype

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
	"spheric.cloud/spheric/internal/api"
	"spheric.cloud/spheric/internal/apis/core"
)

func GetAttrs(obj runtime.Object) (labels.Set, fields.Set, error) {
	diskType, ok := obj.(*core.DiskType)
	if !ok {
		return nil, nil, fmt.Errorf("given object is not a DiskType")
	}
	return diskType.Labels, SelectableFields(diskType), nil
}

func MatchDiskType(label labels.Selector, field fields.Selector) apisrvstorage.SelectionPredicate {
	return apisrvstorage.SelectionPredicate{
		Label:       label,
		Field:       field,
		GetAttrs:    GetAttrs,
		IndexFields: []string{},
	}
}

func SelectableFields(diskType *core.DiskType) fields.Set {
	fieldsSet := make(fields.Set)
	return generic.AddObjectMetaFieldsSet(fieldsSet, &diskType.ObjectMeta, true)
}

func Indexers() *cache.Indexers {
	return &cache.Indexers{}
}

type diskTypeStrategy struct {
	runtime.ObjectTyper
	names.NameGenerator
}

var Strategy = diskTypeStrategy{api.Scheme, names.SimpleNameGenerator}

func (diskTypeStrategy) NamespaceScoped() bool {
	return false
}

func (diskTypeStrategy) PrepareForCreate(ctx context.Context, obj runtime.Object) {
}

func (diskTypeStrategy) PrepareForUpdate(ctx context.Context, obj, old runtime.Object) {
}

func (diskTypeStrategy) Validate(ctx context.Context, obj runtime.Object) field.ErrorList {
	_ = obj.(*core.DiskType)
	return nil
}

func (diskTypeStrategy) WarningsOnCreate(ctx context.Context, obj runtime.Object) []string {
	return nil
}

func (diskTypeStrategy) AllowCreateOnUpdate() bool {
	return false
}

func (diskTypeStrategy) AllowUnconditionalUpdate() bool {
	return false
}

func (diskTypeStrategy) Canonicalize(obj runtime.Object) {
}

func (diskTypeStrategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	newDiskType, oldDiskType := obj.(*core.DiskType), old.(*core.DiskType)
	_, _ = newDiskType, oldDiskType
	return nil
}

func (diskTypeStrategy) WarningsOnUpdate(ctx context.Context, obj, old runtime.Object) []string {
	return nil
}
