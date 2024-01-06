// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package volume

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"strconv"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/apimachinery/pkg/labels"
	sri "spheric.cloud/spheric/sri/apis/volume/v1alpha1"
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

type FakeVolume struct {
	sri.Volume
}

type FakeVolumeClassStatus struct {
	sri.VolumeClassStatus
}

type FakeRuntimeService struct {
	sync.Mutex

	Volumes             map[string]*FakeVolume
	VolumeClassesStatus map[string]*FakeVolumeClassStatus
}

func NewFakeRuntimeService() *FakeRuntimeService {
	return &FakeRuntimeService{
		Volumes:             make(map[string]*FakeVolume),
		VolumeClassesStatus: make(map[string]*FakeVolumeClassStatus),
	}
}

func (r *FakeRuntimeService) SetVolumes(volumes []*FakeVolume) {
	r.Lock()
	defer r.Unlock()

	r.Volumes = make(map[string]*FakeVolume)
	for _, volume := range volumes {
		r.Volumes[volume.Metadata.Id] = volume
	}
}

func (r *FakeRuntimeService) SetVolumeClasses(volumeClassStatus []*FakeVolumeClassStatus) {
	r.Lock()
	defer r.Unlock()

	r.VolumeClassesStatus = make(map[string]*FakeVolumeClassStatus)
	for _, status := range volumeClassStatus {
		r.VolumeClassesStatus[status.VolumeClass.Name] = status
	}
}

func (r *FakeRuntimeService) ListVolumes(ctx context.Context, req *sri.ListVolumesRequest, opts ...grpc.CallOption) (*sri.ListVolumesResponse, error) {
	r.Lock()
	defer r.Unlock()

	filter := req.Filter

	var res []*sri.Volume
	for _, v := range r.Volumes {
		if filter != nil {
			if filter.Id != "" && filter.Id != v.Metadata.Id {
				continue
			}
			if filter.LabelSelector != nil && !filterInLabels(filter.LabelSelector, v.Metadata.Labels) {
				continue
			}
		}

		volume := v.Volume
		res = append(res, &volume)
	}
	return &sri.ListVolumesResponse{Volumes: res}, nil
}

func (r *FakeRuntimeService) CreateVolume(ctx context.Context, req *sri.CreateVolumeRequest, opts ...grpc.CallOption) (*sri.CreateVolumeResponse, error) {
	r.Lock()
	defer r.Unlock()

	volume := *req.Volume
	volume.Metadata.Id = generateID(defaultIDLength)
	volume.Metadata.CreatedAt = time.Now().UnixNano()
	volume.Status = &sri.VolumeStatus{}

	r.Volumes[volume.Metadata.Id] = &FakeVolume{
		Volume: volume,
	}

	return &sri.CreateVolumeResponse{
		Volume: &volume,
	}, nil
}

func (r *FakeRuntimeService) ExpandVolume(ctx context.Context, req *sri.ExpandVolumeRequest, opts ...grpc.CallOption) (*sri.ExpandVolumeResponse, error) {
	r.Lock()
	defer r.Unlock()

	volume, ok := r.Volumes[req.VolumeId]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "volume %q not found", req.VolumeId)
	}

	volume.Spec.Resources.StorageBytes = req.Resources.StorageBytes

	return &sri.ExpandVolumeResponse{}, nil
}

func (r *FakeRuntimeService) DeleteVolume(ctx context.Context, req *sri.DeleteVolumeRequest, opts ...grpc.CallOption) (*sri.DeleteVolumeResponse, error) {
	r.Lock()
	defer r.Unlock()

	volumeID := req.VolumeId
	if _, ok := r.Volumes[volumeID]; !ok {
		return nil, status.Errorf(codes.NotFound, "volume %q not found", volumeID)
	}

	delete(r.Volumes, volumeID)
	return &sri.DeleteVolumeResponse{}, nil
}

func (r *FakeRuntimeService) Status(ctx context.Context, req *sri.StatusRequest, opts ...grpc.CallOption) (*sri.StatusResponse, error) {
	r.Lock()
	defer r.Unlock()

	var res []*sri.VolumeClassStatus
	for _, m := range r.VolumeClassesStatus {
		volumeClassStatus := m.VolumeClassStatus
		res = append(res, &volumeClassStatus)
	}
	return &sri.StatusResponse{VolumeClassStatus: res}, nil
}
