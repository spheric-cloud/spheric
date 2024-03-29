// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package storage

import (
	"context"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/api/meta/table"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"spheric.cloud/spheric/internal/apis/storage"
)

type convertor struct{}

var (
	objectMetaSwaggerDoc = metav1.ObjectMeta{}.SwaggerDoc()

	headers = []metav1.TableColumnDefinition{
		{Name: "Name", Type: "string", Format: "name", Description: objectMetaSwaggerDoc["name"]},
		{Name: "BucketPoolRef", Type: "string", Description: "The bucket pool this bucket is hosted on"},
		{Name: "BucketClass", Type: "string", Description: "The bucket class of this bucket"},
		{Name: "State", Type: "string", Description: "The state of the bucket on the underlying bucket pool"},
		{Name: "Age", Type: "string", Format: "date", Description: objectMetaSwaggerDoc["creationTimestamp"]},
	}
)

func newTableConvertor() *convertor {
	return &convertor{}
}

func (c *convertor) ConvertToTable(ctx context.Context, obj runtime.Object, tableOptions runtime.Object) (*metav1.Table, error) {
	tab := &metav1.Table{
		ColumnDefinitions: headers,
	}

	if m, err := meta.ListAccessor(obj); err == nil {
		tab.ResourceVersion = m.GetResourceVersion()
		tab.Continue = m.GetContinue()
	} else {
		if m, err := meta.CommonAccessor(obj); err == nil {
			tab.ResourceVersion = m.GetResourceVersion()
		}
	}

	var err error
	tab.Rows, err = table.MetaToTableRow(obj, func(obj runtime.Object, m metav1.Object, name, age string) (cells []interface{}, err error) {
		bucket := obj.(*storage.Bucket)

		cells = append(cells, name)
		if bucketPoolRef := bucket.Spec.BucketPoolRef; bucketPoolRef != nil {
			cells = append(cells, bucketPoolRef.Name)
		} else {
			cells = append(cells, "<none>")
		}
		if bucketClassRef := bucket.Spec.BucketClassRef; bucketClassRef != nil {
			cells = append(cells, bucketClassRef.Name)
		} else {
			cells = append(cells, "<none>")
		}
		if state := bucket.Status.State; state != "" {
			cells = append(cells, state)
		} else {
			cells = append(cells, "<unknown>")
		}
		cells = append(cells, age)

		return cells, nil
	})
	return tab, err
}
