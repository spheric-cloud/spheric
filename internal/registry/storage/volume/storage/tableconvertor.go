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
		{Name: "VolumePoolRef", Type: "string", Description: "The volume pool this volume is hosted on"},
		{Name: "Image", Type: "string", Description: "The image the volume should be populated from"},
		{Name: "VolumeClass", Type: "string", Description: "The volume class of this volume"},
		{Name: "State", Type: "string", Description: "The state of the volume on the underlying volume pool"},
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
		volume := obj.(*storage.Volume)

		cells = append(cells, name)
		if volumePoolRef := volume.Spec.VolumePoolRef; volumePoolRef != nil {
			cells = append(cells, volumePoolRef.Name)
		} else {
			cells = append(cells, "<none>")
		}
		if image := volume.Spec.Image; image != "" {
			cells = append(cells, image)
		} else {
			cells = append(cells, "<none>")
		}
		if volumeClassRef := volume.Spec.VolumeClassRef; volumeClassRef != nil {
			cells = append(cells, volumeClassRef.Name)
		} else {
			cells = append(cells, "<none>")
		}
		if state := volume.Status.State; state != "" {
			cells = append(cells, state)
		} else {
			cells = append(cells, "<unknown>")
		}
		cells = append(cells, age)

		return cells, nil
	})
	return tab, err
}
