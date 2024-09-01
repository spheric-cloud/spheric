// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package instancetype

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/api/equality"

	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/storage"
	"k8s.io/apiserver/pkg/storage/names"
	"spheric.cloud/spheric/internal/api"
	"spheric.cloud/spheric/internal/apis/core"
	"spheric.cloud/spheric/internal/apis/core/validation"
)

func GetAttrs(obj runtime.Object) (labels.Set, fields.Set, error) {
	instanceType, ok := obj.(*core.InstanceType)
	if !ok {
		return nil, nil, fmt.Errorf("given object is not a InstanceTypeRef")
	}
	return instanceType.Labels, SelectableFields(instanceType), nil
}

func MatchInstanceType(label labels.Selector, field fields.Selector) storage.SelectionPredicate {
	return storage.SelectionPredicate{
		Label:    label,
		Field:    field,
		GetAttrs: GetAttrs,
	}
}

func SelectableFields(machine *core.InstanceType) fields.Set {
	return generic.ObjectMetaFieldsSet(&machine.ObjectMeta, false)
}

type instanceTypeStrategy struct {
	runtime.ObjectTyper
	names.NameGenerator
}

var Strategy = instanceTypeStrategy{api.Scheme, names.SimpleNameGenerator}

func (instanceTypeStrategy) NamespaceScoped() bool {
	return false
}

func (instanceTypeStrategy) PrepareForCreate(ctx context.Context, obj runtime.Object) {
	instanceType := obj.(*core.InstanceType)
	instanceType.Generation = 1
}

func (instanceTypeStrategy) PrepareForUpdate(ctx context.Context, obj, old runtime.Object) {
	newInstanceType := obj.(*core.InstanceType)
	oldInstanceType := old.(*core.InstanceType)

	if !equality.Semantic.DeepEqual(newInstanceType.Capabilities, oldInstanceType.Capabilities) {
		newInstanceType.Generation = oldInstanceType.Generation + 1
	}
}

func (instanceTypeStrategy) Validate(ctx context.Context, obj runtime.Object) field.ErrorList {
	instanceType := obj.(*core.InstanceType)
	return validation.ValidateInstanceType(instanceType)
}

func (instanceTypeStrategy) WarningsOnCreate(ctx context.Context, obj runtime.Object) []string {
	return nil
}

func (instanceTypeStrategy) AllowCreateOnUpdate() bool {
	return false
}

func (instanceTypeStrategy) AllowUnconditionalUpdate() bool {
	return false
}

func (instanceTypeStrategy) Canonicalize(obj runtime.Object) {
}

func (instanceTypeStrategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	newInstanceType := obj.(*core.InstanceType)
	oldInstanceType := old.(*core.InstanceType)
	return validation.ValidateInstanceTypeUpdate(newInstanceType, oldInstanceType)
}

func (instanceTypeStrategy) WarningsOnUpdate(ctx context.Context, obj, old runtime.Object) []string {
	return nil
}
