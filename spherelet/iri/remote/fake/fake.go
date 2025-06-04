// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package fake

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strconv"
	"sync"
	"time"

	"google.golang.org/protobuf/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/apimachinery/pkg/labels"
	iri "spheric.cloud/spheric/iri-api/apis/runtime/v1alpha1"
)

var (
	// Version is the version of the fake runtime.
	Version = "0.1.0"

	// RuntimeName is the name of the fake runtime.
	RuntimeName = "fakeRuntime"
)

func filterInLabels(labelSelector, lbls map[string]string) bool {
	return labels.SelectorFromSet(labelSelector).Matches(labels.Set(lbls))
}

const defaultIDLength = 63

func generateID(length int) string {
	data := make([]byte, (length/2)+1)
	for {
		_, _ = rand.Read(data)
		id := hex.EncodeToString(data)

		// Truncated versions of the id should not be numerical.
		if _, err := strconv.ParseInt(id[:12], 10, 64); err != nil {
			continue
		}

		return id[:length]
	}
}

type FakeInstance struct {
	iri.Instance
}

type FakeStatus struct {
	Capacity    iri.RuntimeResources
	Allocatable iri.RuntimeResources
}

type FakeRuntimeService struct {
	sync.RWMutex

	Instances   map[string]*FakeInstance
	Capacity    *iri.RuntimeResources
	Allocatable *iri.RuntimeResources
	GetExecURL  func(req *iri.ExecRequest) string
}

func NewFakeRuntimeService() *FakeRuntimeService {
	return &FakeRuntimeService{
		Instances: make(map[string]*FakeInstance),
	}
}

func (r *FakeRuntimeService) GetInstance(id string) (*FakeInstance, error) {
	r.RLock()
	defer r.RUnlock()

	for _, inst := range r.Instances {
		if actualID := inst.GetMetadata().GetId(); actualID == id {
			return inst, nil
		}
	}
	return nil, fmt.Errorf("no instance with id %q", id)
}

func (r *FakeRuntimeService) UpdateInstance(inst *FakeInstance) error {
	r.Lock()
	defer r.Unlock()

	id := inst.GetMetadata().GetId()

	for i, found := range r.Instances {
		if found.GetMetadata().GetId() == id {
			r.Instances[i] = inst
			return nil
		}
	}
	return fmt.Errorf("no instance with id %q found", id)
}

func (r *FakeRuntimeService) GetFirstInstanceByLabel(label, value string) (*FakeInstance, error) {
	r.RLock()
	defer r.RUnlock()

	for _, inst := range r.Instances {
		if actual, ok := inst.GetMetadata().GetLabels()[label]; ok && actual == value {
			return inst, nil
		}
	}
	return nil, fmt.Errorf("no instance with label %s=%s found", label, value)
}

func (r *FakeRuntimeService) SetInstances(instances []*FakeInstance) {
	r.Lock()
	defer r.Unlock()

	r.Instances = make(map[string]*FakeInstance)
	for _, instance := range instances {
		r.Instances[instance.Metadata.Id] = instance
	}
}

func (r *FakeRuntimeService) SetStatus(capacity, allocatable *iri.RuntimeResources) {
	r.Lock()
	defer r.Unlock()

	r.Capacity = capacity
	r.Allocatable = allocatable
}

func (r *FakeRuntimeService) SetGetExecURL(f func(req *iri.ExecRequest) string) {
	r.Lock()
	defer r.Unlock()

	r.GetExecURL = f
}

func (r *FakeRuntimeService) Version(ctx context.Context, req *iri.VersionRequest) (*iri.VersionResponse, error) {
	return &iri.VersionResponse{
		RuntimeName:    RuntimeName,
		RuntimeVersion: Version,
	}, nil
}

func (r *FakeRuntimeService) ListInstances(ctx context.Context, req *iri.ListInstancesRequest) (*iri.ListInstancesResponse, error) {
	r.Lock()
	defer r.Unlock()

	filter := req.Filter

	var res []*iri.Instance
	for _, m := range r.Instances {
		if filter != nil {
			if filter.Id != "" && filter.Id != m.Metadata.Id {
				continue
			}
			if filter.LabelSelector != nil && !filterInLabels(filter.LabelSelector, m.Metadata.Labels) {
				continue
			}
		}

		instance := proto.Clone(&m.Instance).(*iri.Instance) //nolint
		res = append(res, instance)
	}
	return &iri.ListInstancesResponse{Instances: res}, nil
}

func (r *FakeRuntimeService) CreateInstance(ctx context.Context, req *iri.CreateInstanceRequest) (*iri.CreateInstanceResponse, error) {
	r.Lock()
	defer r.Unlock()

	fakeInst := &FakeInstance{}
	proto.Merge(&fakeInst.Instance, req.Instance)

	fakeInst.Metadata.Id = generateID(defaultIDLength)
	fakeInst.Metadata.CreatedAt = time.Now().UnixNano()
	fakeInst.Status = &iri.InstanceStatus{
		State: iri.InstanceState_INSTANCE_PENDING,
	}

	r.Instances[fakeInst.Metadata.Id] = fakeInst

	return &iri.CreateInstanceResponse{
		Instance: &fakeInst.Instance,
	}, nil
}

func (r *FakeRuntimeService) DeleteInstance(ctx context.Context, req *iri.DeleteInstanceRequest) (*iri.DeleteInstanceResponse, error) {
	r.Lock()
	defer r.Unlock()

	instanceID := req.InstanceId
	if _, ok := r.Instances[instanceID]; !ok {
		return nil, status.Errorf(codes.NotFound, "instance %q not found", instanceID)
	}

	delete(r.Instances, instanceID)
	return &iri.DeleteInstanceResponse{}, nil
}

func (r *FakeRuntimeService) UpdateInstanceAnnotations(ctx context.Context, req *iri.UpdateInstanceAnnotationsRequest) (*iri.UpdateInstanceAnnotationsResponse, error) {
	r.Lock()
	defer r.Unlock()

	instanceID := req.InstanceId
	instance, ok := r.Instances[instanceID]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "instance %q not found", instanceID)
	}

	instance.Metadata.Annotations = req.Annotations
	return &iri.UpdateInstanceAnnotationsResponse{}, nil
}

func (r *FakeRuntimeService) UpdateInstancePower(ctx context.Context, req *iri.UpdateInstancePowerRequest) (*iri.UpdateInstancePowerResponse, error) {
	r.Lock()
	defer r.Unlock()

	instanceID := req.InstanceId
	instance, ok := r.Instances[instanceID]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "instance %q not found", instanceID)
	}

	instance.Spec.Power = req.Power
	return &iri.UpdateInstancePowerResponse{}, nil
}

func (r *FakeRuntimeService) AttachDisk(ctx context.Context, req *iri.AttachDiskRequest) (*iri.AttachDiskResponse, error) {
	r.Lock()
	defer r.Unlock()

	instanceID := req.InstanceId
	instance, ok := r.Instances[instanceID]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "instance %q not found", instanceID)
	}

	instance.Spec.Disks = append(instance.Spec.Disks, req.Disk)
	return &iri.AttachDiskResponse{}, nil
}

func (r *FakeRuntimeService) DetachDisk(ctx context.Context, req *iri.DetachDiskRequest) (*iri.DetachDiskResponse, error) {
	r.Lock()
	defer r.Unlock()

	instanceID := req.InstanceId
	instance, ok := r.Instances[instanceID]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "instance %q not found", instanceID)
	}

	var (
		filtered []*iri.Disk
		found    bool
	)
	for _, attachment := range instance.Spec.Disks {
		if attachment.Name == req.Name {
			found = true
			continue
		}

		filtered = append(filtered, attachment)
	}
	if !found {
		return nil, status.Errorf(codes.NotFound, "instance %q disk attachment %q not found", instanceID, req.Name)
	}

	instance.Spec.Disks = filtered
	return &iri.DetachDiskResponse{}, nil
}

func (r *FakeRuntimeService) AttachNetworkInterface(ctx context.Context, req *iri.AttachNetworkInterfaceRequest) (*iri.AttachNetworkInterfaceResponse, error) {
	r.Lock()
	defer r.Unlock()

	instanceID := req.InstanceId
	instance, ok := r.Instances[instanceID]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "instance %q not found", instanceID)
	}

	instance.Spec.NetworkInterfaces = append(instance.Spec.NetworkInterfaces, req.NetworkInterface)
	return &iri.AttachNetworkInterfaceResponse{}, nil
}

func (r *FakeRuntimeService) DetachNetworkInterface(ctx context.Context, req *iri.DetachNetworkInterfaceRequest) (*iri.DetachNetworkInterfaceResponse, error) {
	r.Lock()
	defer r.Unlock()

	instanceID := req.InstanceId
	instance, ok := r.Instances[instanceID]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "instance %q not found", instanceID)
	}

	var (
		filtered []*iri.NetworkInterface
		found    bool
	)
	for _, attachment := range instance.Spec.NetworkInterfaces {
		if attachment.Name == req.Name {
			found = true
			continue
		}

		filtered = append(filtered, attachment)
	}
	if !found {
		return nil, status.Errorf(codes.NotFound, "instance %q network interface attachment %q not found", instanceID, req.Name)
	}

	instance.Spec.NetworkInterfaces = filtered
	return &iri.DetachNetworkInterfaceResponse{}, nil
}

func (r *FakeRuntimeService) Status(ctx context.Context, req *iri.StatusRequest) (*iri.StatusResponse, error) {
	r.Lock()
	defer r.Unlock()

	return &iri.StatusResponse{
		Capacity:    r.Capacity,
		Allocatable: r.Allocatable,
	}, nil
}

func (r *FakeRuntimeService) Exec(ctx context.Context, req *iri.ExecRequest) (*iri.ExecResponse, error) {
	r.Lock()
	defer r.Unlock()

	var url string
	if r.GetExecURL != nil {
		url = r.GetExecURL(req)
	}
	return &iri.ExecResponse{Url: url}, nil
}
