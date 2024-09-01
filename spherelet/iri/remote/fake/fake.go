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
	sri "spheric.cloud/spheric/iri-api/apis/runtime/v1alpha1"
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
	sri.Instance
}

type FakeStatus struct {
	Capacity    sri.RuntimeResources
	Allocatable sri.RuntimeResources
}

type FakeRuntimeService struct {
	sync.RWMutex

	Instances   map[string]*FakeInstance
	Capacity    *sri.RuntimeResources
	Allocatable *sri.RuntimeResources
	GetExecURL  func(req *sri.ExecRequest) string
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

func (r *FakeRuntimeService) SetStatus(capacity, allocatable *sri.RuntimeResources) {
	r.Lock()
	defer r.Unlock()

	r.Capacity = capacity
	r.Allocatable = allocatable
}

func (r *FakeRuntimeService) SetGetExecURL(f func(req *sri.ExecRequest) string) {
	r.Lock()
	defer r.Unlock()

	r.GetExecURL = f
}

func (r *FakeRuntimeService) Version(ctx context.Context, req *sri.VersionRequest) (*sri.VersionResponse, error) {
	return &sri.VersionResponse{
		RuntimeName:    RuntimeName,
		RuntimeVersion: Version,
	}, nil
}

func (r *FakeRuntimeService) ListInstances(ctx context.Context, req *sri.ListInstancesRequest) (*sri.ListInstancesResponse, error) {
	r.Lock()
	defer r.Unlock()

	filter := req.Filter

	var res []*sri.Instance
	for _, m := range r.Instances {
		if filter != nil {
			if filter.Id != "" && filter.Id != m.Metadata.Id {
				continue
			}
			if filter.LabelSelector != nil && !filterInLabels(filter.LabelSelector, m.Metadata.Labels) {
				continue
			}
		}

		instance := proto.Clone(&m.Instance).(*sri.Instance) //nolint
		res = append(res, instance)
	}
	return &sri.ListInstancesResponse{Instances: res}, nil
}

func (r *FakeRuntimeService) CreateInstance(ctx context.Context, req *sri.CreateInstanceRequest) (*sri.CreateInstanceResponse, error) {
	r.Lock()
	defer r.Unlock()

	fakeInst := &FakeInstance{}
	proto.Merge(&fakeInst.Instance, req.Instance)

	fakeInst.Instance.Metadata.Id = generateID(defaultIDLength)
	fakeInst.Instance.Metadata.CreatedAt = time.Now().UnixNano()
	fakeInst.Instance.Status = &sri.InstanceStatus{
		State: sri.InstanceState_INSTANCE_PENDING,
	}

	r.Instances[fakeInst.Instance.Metadata.Id] = fakeInst

	return &sri.CreateInstanceResponse{
		Instance: &fakeInst.Instance,
	}, nil
}

func (r *FakeRuntimeService) DeleteInstance(ctx context.Context, req *sri.DeleteInstanceRequest) (*sri.DeleteInstanceResponse, error) {
	r.Lock()
	defer r.Unlock()

	instanceID := req.InstanceId
	if _, ok := r.Instances[instanceID]; !ok {
		return nil, status.Errorf(codes.NotFound, "instance %q not found", instanceID)
	}

	delete(r.Instances, instanceID)
	return &sri.DeleteInstanceResponse{}, nil
}

func (r *FakeRuntimeService) UpdateInstanceAnnotations(ctx context.Context, req *sri.UpdateInstanceAnnotationsRequest) (*sri.UpdateInstanceAnnotationsResponse, error) {
	r.Lock()
	defer r.Unlock()

	instanceID := req.InstanceId
	instance, ok := r.Instances[instanceID]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "instance %q not found", instanceID)
	}

	instance.Metadata.Annotations = req.Annotations
	return &sri.UpdateInstanceAnnotationsResponse{}, nil
}

func (r *FakeRuntimeService) UpdateInstancePower(ctx context.Context, req *sri.UpdateInstancePowerRequest) (*sri.UpdateInstancePowerResponse, error) {
	r.Lock()
	defer r.Unlock()

	instanceID := req.InstanceId
	instance, ok := r.Instances[instanceID]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "instance %q not found", instanceID)
	}

	instance.Spec.Power = req.Power
	return &sri.UpdateInstancePowerResponse{}, nil
}

func (r *FakeRuntimeService) AttachDisk(ctx context.Context, req *sri.AttachDiskRequest) (*sri.AttachDiskResponse, error) {
	r.Lock()
	defer r.Unlock()

	instanceID := req.InstanceId
	instance, ok := r.Instances[instanceID]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "instance %q not found", instanceID)
	}

	instance.Spec.Disks = append(instance.Spec.Disks, req.Disk)
	return &sri.AttachDiskResponse{}, nil
}

func (r *FakeRuntimeService) DetachDisk(ctx context.Context, req *sri.DetachDiskRequest) (*sri.DetachDiskResponse, error) {
	r.Lock()
	defer r.Unlock()

	instanceID := req.InstanceId
	instance, ok := r.Instances[instanceID]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "instance %q not found", instanceID)
	}

	var (
		filtered []*sri.Disk
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
	return &sri.DetachDiskResponse{}, nil
}

func (r *FakeRuntimeService) AttachNetworkInterface(ctx context.Context, req *sri.AttachNetworkInterfaceRequest) (*sri.AttachNetworkInterfaceResponse, error) {
	r.Lock()
	defer r.Unlock()

	instanceID := req.InstanceId
	instance, ok := r.Instances[instanceID]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "instance %q not found", instanceID)
	}

	instance.Spec.NetworkInterfaces = append(instance.Spec.NetworkInterfaces, req.NetworkInterface)
	return &sri.AttachNetworkInterfaceResponse{}, nil
}

func (r *FakeRuntimeService) DetachNetworkInterface(ctx context.Context, req *sri.DetachNetworkInterfaceRequest) (*sri.DetachNetworkInterfaceResponse, error) {
	r.Lock()
	defer r.Unlock()

	instanceID := req.InstanceId
	instance, ok := r.Instances[instanceID]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "instance %q not found", instanceID)
	}

	var (
		filtered []*sri.NetworkInterface
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
	return &sri.DetachNetworkInterfaceResponse{}, nil
}

func (r *FakeRuntimeService) Status(ctx context.Context, req *sri.StatusRequest) (*sri.StatusResponse, error) {
	r.Lock()
	defer r.Unlock()

	return &sri.StatusResponse{
		Capacity:    r.Capacity,
		Allocatable: r.Allocatable,
	}, nil
}

func (r *FakeRuntimeService) Exec(ctx context.Context, req *sri.ExecRequest) (*sri.ExecResponse, error) {
	r.Lock()
	defer r.Unlock()

	var url string
	if r.GetExecURL != nil {
		url = r.GetExecURL(req)
	}
	return &sri.ExecResponse{Url: url}, nil
}
