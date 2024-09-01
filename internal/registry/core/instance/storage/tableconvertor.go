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
	"spheric.cloud/spheric/internal/apis/core"
)

type convertor struct{}

var (
	objectMetaSwaggerDoc = metav1.ObjectMeta{}.SwaggerDoc()

	headers = []metav1.TableColumnDefinition{
		{Name: "Name", Type: "string", Format: "name", Description: objectMetaSwaggerDoc["name"]},
		{Name: "InstanceClassRef", Type: "string", Description: "The instance class of this instance"},
		{Name: "Image", Type: "string", Description: "The image the instance shall use"},
		{Name: "FleetRef", Type: "string", Description: "The instance pool the instance is running on"},
		{Name: "State", Type: "string", Description: "The current state of the instance"},
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
		instance := obj.(*core.Instance)

		cells = append(cells, name)
		cells = append(cells, instance.Spec.InstanceTypeRef.Name)
		if image := instance.Spec.Image; image != "" {
			cells = append(cells, image)
		} else {
			cells = append(cells, "<none>")
		}
		if fleetRef := instance.Spec.FleetRef; fleetRef != nil {
			cells = append(cells, fleetRef.Name)
		} else {
			cells = append(cells, "<none>")
		}
		if state := instance.Status.State; state != "" {
			cells = append(cells, state)
		} else {
			cells = append(cells, "<unknown>")
		}
		cells = append(cells, age)

		return cells, nil
	})
	return tab, err
}
