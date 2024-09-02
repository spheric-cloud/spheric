// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"context"
	"fmt"

	utilreconcile "spheric.cloud/spheric/utils/reconcile"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"spheric.cloud/spheric/utils/annotations"
)

func PatchEnsureNoReconcileAnnotation(ctx context.Context, c client.Client, obj client.Object) (modified bool, err error) {
	if !annotations.HasReconcileAnnotation(obj) {
		return false, nil
	}

	if err := PatchRemoveReconcileAnnotation(ctx, c, obj); err != nil {
		return false, err
	}
	return true, nil
}

func PatchAddReconcileAnnotation(ctx context.Context, c client.Client, obj client.Object) error {
	base := obj.DeepCopyObject().(client.Object)

	annotations.SetReconcileAnnotation(obj)

	if err := c.Patch(ctx, obj, client.MergeFrom(base)); err != nil {
		return fmt.Errorf("error adding reconcile annotation: %w", err)
	}
	return nil
}

func PatchRemoveReconcileAnnotation(ctx context.Context, c client.Client, obj client.Object) error {
	base := obj.DeepCopyObject().(client.Object)

	annotations.RemoveReconcileAnnotation(obj)

	if err := c.Patch(ctx, obj, client.MergeFrom(base)); err != nil {
		return fmt.Errorf("error removing reconcile annotation: %w", err)
	}
	return nil
}

type Object[O any] interface {
	client.Object
	*O
}

func ReconcileRequestsFromObjectStructSlice[O Object[OStruct], S ~[]OStruct, OStruct any](objs S) []reconcile.Request {
	res := make([]reconcile.Request, len(objs))
	for i := range objs {
		obj := O(&objs[i])
		res[i] = utilreconcile.RequestFromObject(obj)
	}
	return res
}
