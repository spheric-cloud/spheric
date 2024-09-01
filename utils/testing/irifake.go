// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and Spheric contributors
// SPDX-License-Identifier: Apache-2.0

package testing

import (
	"google.golang.org/protobuf/proto"
	"k8s.io/apimachinery/pkg/types"
	iri "spheric.cloud/spheric/iri-api/apis/runtime/v1alpha1"
	sphereletv1alpha1 "spheric.cloud/spheric/spherelet/api/v1alpha1"
	"spheric.cloud/spheric/spherelet/iri/remote/fake"
)

func NewFakeInstanceWithUID(uid types.UID) *fake.FakeInstance {
	return &fake.FakeInstance{
		Instance: iri.Instance{
			Metadata: &iri.ObjectMetadata{
				Labels: map[string]string{
					sphereletv1alpha1.InstanceUIDLabel: string(uid),
				},
			},
		},
	}
}

func getUID(inst *fake.FakeInstance) string {
	uid := inst.GetMetadata().GetLabels()[sphereletv1alpha1.InstanceUIDLabel]
	return uid
}

func GetInstanceByUID(srv *fake.FakeRuntimeService, inst *fake.FakeInstance) func() error {
	uid := getUID(inst)

	return func() error {
		found, err := srv.GetFirstInstanceByLabel(sphereletv1alpha1.InstanceUIDLabel, uid)
		if err != nil {
			return err
		}

		proto.Reset(inst)
		proto.Merge(inst, found)
		return nil
	}
}

func GetInstance(srv *fake.FakeRuntimeService, inst *fake.FakeInstance) func() error {
	id := inst.GetMetadata().GetId()

	return func() error {
		found, err := srv.GetInstance(id)
		if err != nil {
			return err
		}

		proto.Reset(inst)
		proto.Merge(inst, found)
		return nil
	}
}

func Instance(srv *fake.FakeRuntimeService, inst *fake.FakeInstance) func() (*fake.FakeInstance, error) {
	id := inst.GetMetadata().GetId()

	return func() (*fake.FakeInstance, error) {
		found, err := srv.GetInstance(id)
		if err != nil {
			return nil, err
		}

		proto.Reset(inst)
		proto.Merge(inst, found)
		return inst, nil
	}
}

func UpdateInstance(srv *fake.FakeRuntimeService, inst *fake.FakeInstance, update func()) func() error {
	id := inst.GetMetadata().GetId()
	return func() error {
		found, err := srv.GetInstance(id)
		if err != nil {
			return err
		}

		proto.Reset(inst)
		proto.Merge(inst, found)

		update()

		return srv.UpdateInstance(inst)
	}
}
