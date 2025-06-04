// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0

package controllers

import (
	"context"
	"errors"
	"net/http"
	"os"
	"path/filepath"

	"spheric.cloud/spheric/actuo/client"
	"spheric.cloud/spheric/actuo/reconcile"
	cloudhypervisor "spheric.cloud/spheric/cloud-hypervisor"
	chypclient "spheric.cloud/spheric/cloud-hypervisor/client"
	oapiclient "spheric.cloud/spheric/cloud-hypervisor/oapi-client"
	"spheric.cloud/spheric/vee/api"
)

type InstanceReconciler struct {
	client.Client[string, *api.Instance]
	CloudHypervisor cloudhypervisor.CloudHypervisor
	SocketDir       string
}

func (r *InstanceReconciler) Reconcile(ctx context.Context, id string) (reconcile.Result, error) {
	instance, err := r.Get(ctx, id)
	if err != nil {
		return reconcile.Result{}, err
	}

	apiSocket := filepath.Join(r.SocketDir, instance.ID)
	cHyp, err := chypclient.Connect(apiSocket)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return reconcile.Result{}, err
		}

		newCHyp, err := r.CloudHypervisor.Start(ctx, apiSocket)
		if err != nil {
			return reconcile.Result{}, err
		}

		cHyp = newCHyp
	}

	vmInfo, err := cHyp.GetVMInfo(ctx)
	if err != nil {
		if !chypclient.IsStatusError(err, http.StatusNotFound) {
			return reconcile.Result{}, err
		}
	}

	_ = vmInfo
	return reconcile.Result{}, nil
}

func (r *InstanceReconciler) create(ctx context.Context, cHyp chypclient.Client, instance *api.Instance) error {
	cpusConfig := &oapiclient.CpusConfig{
		BootVcpus: int(instance.Spec.CPUCount),
	}

	memoryConfig := &oapiclient.MemoryConfig{
		Size: instance.Spec.MemoryBytes,
	}

	if err := cHyp.CreateVM(ctx, oapiclient.CreateVMJSONRequestBody{
		Cpus:    cpusConfig,
		Memory:  memoryConfig,
		Payload: oapiclient.PayloadConfig{},
	}); err != nil {
		return err
	}

	vmInfo, err := cHyp.GetVMInfo(ctx)
	if err != nil {
		return err
	}

	_ = vmInfo
	return nil
}
