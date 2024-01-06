// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package machine

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"strconv"
	"sync"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/apimachinery/pkg/labels"
	sri "spheric.cloud/spheric/sri/apis/machine/v1alpha1"
)

var (
	// FakeVersion is the version of the fake runtime.
	FakeVersion = "0.1.0"

	// FakeRuntimeName is the name of the fake runtime.
	FakeRuntimeName = "fakeRuntime"
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

type FakeMachine struct {
	sri.Machine
}

type FakeVolume struct {
	sri.Volume
}

type FakeNetworkInterface struct {
	sri.NetworkInterface
}

type FakeMachineClassStatus struct {
	sri.MachineClassStatus
}

type FakeRuntimeService struct {
	sync.Mutex

	Machines           map[string]*FakeMachine
	MachineClassStatus map[string]*FakeMachineClassStatus
	GetExecURL         func(req *sri.ExecRequest) string
}

func NewFakeRuntimeService() *FakeRuntimeService {
	return &FakeRuntimeService{
		Machines:           make(map[string]*FakeMachine),
		MachineClassStatus: make(map[string]*FakeMachineClassStatus),
	}
}

func (r *FakeRuntimeService) SetMachines(machines []*FakeMachine) {
	r.Lock()
	defer r.Unlock()

	r.Machines = make(map[string]*FakeMachine)
	for _, machine := range machines {
		r.Machines[machine.Metadata.Id] = machine
	}
}

func (r *FakeRuntimeService) SetMachineClasses(machineClassStatus []*FakeMachineClassStatus) {
	r.Lock()
	defer r.Unlock()

	r.MachineClassStatus = make(map[string]*FakeMachineClassStatus)
	for _, status := range machineClassStatus {
		r.MachineClassStatus[status.MachineClass.Name] = status
	}
}

func (r *FakeRuntimeService) SetGetExecURL(f func(req *sri.ExecRequest) string) {
	r.Lock()
	defer r.Unlock()

	r.GetExecURL = f
}

func (r *FakeRuntimeService) Version(ctx context.Context, req *sri.VersionRequest) (*sri.VersionResponse, error) {
	return &sri.VersionResponse{
		RuntimeName:    FakeRuntimeName,
		RuntimeVersion: FakeVersion,
	}, nil
}

func (r *FakeRuntimeService) ListMachines(ctx context.Context, req *sri.ListMachinesRequest) (*sri.ListMachinesResponse, error) {
	r.Lock()
	defer r.Unlock()

	filter := req.Filter

	var res []*sri.Machine
	for _, m := range r.Machines {
		if filter != nil {
			if filter.Id != "" && filter.Id != m.Metadata.Id {
				continue
			}
			if filter.LabelSelector != nil && !filterInLabels(filter.LabelSelector, m.Metadata.Labels) {
				continue
			}
		}

		machine := m.Machine
		res = append(res, &machine)
	}
	return &sri.ListMachinesResponse{Machines: res}, nil
}

func (r *FakeRuntimeService) CreateMachine(ctx context.Context, req *sri.CreateMachineRequest) (*sri.CreateMachineResponse, error) {
	r.Lock()
	defer r.Unlock()

	machine := *req.Machine
	machine.Metadata.Id = generateID(defaultIDLength)
	machine.Metadata.CreatedAt = time.Now().UnixNano()
	machine.Status = &sri.MachineStatus{
		State: sri.MachineState_MACHINE_PENDING,
	}

	r.Machines[machine.Metadata.Id] = &FakeMachine{
		Machine: machine,
	}

	return &sri.CreateMachineResponse{
		Machine: &machine,
	}, nil
}

func (r *FakeRuntimeService) DeleteMachine(ctx context.Context, req *sri.DeleteMachineRequest) (*sri.DeleteMachineResponse, error) {
	r.Lock()
	defer r.Unlock()

	machineID := req.MachineId
	if _, ok := r.Machines[machineID]; !ok {
		return nil, status.Errorf(codes.NotFound, "machine %q not found", machineID)
	}

	delete(r.Machines, machineID)
	return &sri.DeleteMachineResponse{}, nil
}

func (r *FakeRuntimeService) UpdateMachineAnnotations(ctx context.Context, req *sri.UpdateMachineAnnotationsRequest) (*sri.UpdateMachineAnnotationsResponse, error) {
	r.Lock()
	defer r.Unlock()

	machineID := req.MachineId
	machine, ok := r.Machines[machineID]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "machine %q not found", machineID)
	}

	machine.Metadata.Annotations = req.Annotations
	return &sri.UpdateMachineAnnotationsResponse{}, nil
}

func (r *FakeRuntimeService) UpdateMachinePower(ctx context.Context, req *sri.UpdateMachinePowerRequest) (*sri.UpdateMachinePowerResponse, error) {
	r.Lock()
	defer r.Unlock()

	machineID := req.MachineId
	machine, ok := r.Machines[machineID]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "machine %q not found", machineID)
	}

	machine.Spec.Power = req.Power
	return &sri.UpdateMachinePowerResponse{}, nil
}

func (r *FakeRuntimeService) AttachVolume(ctx context.Context, req *sri.AttachVolumeRequest) (*sri.AttachVolumeResponse, error) {
	r.Lock()
	defer r.Unlock()

	machineID := req.MachineId
	machine, ok := r.Machines[machineID]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "machine %q not found", machineID)
	}

	machine.Spec.Volumes = append(machine.Spec.Volumes, req.Volume)
	return &sri.AttachVolumeResponse{}, nil
}

func (r *FakeRuntimeService) DetachVolume(ctx context.Context, req *sri.DetachVolumeRequest) (*sri.DetachVolumeResponse, error) {
	r.Lock()
	defer r.Unlock()

	machineID := req.MachineId
	machine, ok := r.Machines[machineID]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "machine %q not found", machineID)
	}

	var (
		filtered []*sri.Volume
		found    bool
	)
	for _, attachment := range machine.Spec.Volumes {
		if attachment.Name == req.Name {
			found = true
			continue
		}

		filtered = append(filtered, attachment)
	}
	if !found {
		return nil, status.Errorf(codes.NotFound, "machine %q volume attachment %q not found", machineID, req.Name)
	}

	machine.Spec.Volumes = filtered
	return &sri.DetachVolumeResponse{}, nil
}

func (r *FakeRuntimeService) AttachNetworkInterface(ctx context.Context, req *sri.AttachNetworkInterfaceRequest) (*sri.AttachNetworkInterfaceResponse, error) {
	r.Lock()
	defer r.Unlock()

	machineID := req.MachineId
	machine, ok := r.Machines[machineID]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "machine %q not found", machineID)
	}

	machine.Spec.NetworkInterfaces = append(machine.Spec.NetworkInterfaces, req.NetworkInterface)
	return &sri.AttachNetworkInterfaceResponse{}, nil
}

func (r *FakeRuntimeService) DetachNetworkInterface(ctx context.Context, req *sri.DetachNetworkInterfaceRequest) (*sri.DetachNetworkInterfaceResponse, error) {
	r.Lock()
	defer r.Unlock()

	machineID := req.MachineId
	machine, ok := r.Machines[machineID]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "machine %q not found", machineID)
	}

	var (
		filtered []*sri.NetworkInterface
		found    bool
	)
	for _, attachment := range machine.Spec.NetworkInterfaces {
		if attachment.Name == req.Name {
			found = true
			continue
		}

		filtered = append(filtered, attachment)
	}
	if !found {
		return nil, status.Errorf(codes.NotFound, "machine %q network interface attachment %q not found", machineID, req.Name)
	}

	machine.Spec.NetworkInterfaces = filtered
	return &sri.DetachNetworkInterfaceResponse{}, nil
}

func (r *FakeRuntimeService) Status(ctx context.Context, req *sri.StatusRequest) (*sri.StatusResponse, error) {
	r.Lock()
	defer r.Unlock()

	var res []*sri.MachineClassStatus
	for _, m := range r.MachineClassStatus {
		machineClassStatus := m.MachineClassStatus
		res = append(res, &machineClassStatus)
	}
	return &sri.StatusResponse{MachineClassStatus: res}, nil
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
