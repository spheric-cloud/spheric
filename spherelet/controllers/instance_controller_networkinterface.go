// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package controllers

import (
	"context"
	"errors"
	"fmt"

	corev1alpha1 "spheric.cloud/spheric/api/core/v1alpha1"

	"github.com/go-logr/logr"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	"sigs.k8s.io/controller-runtime/pkg/client"
	iri "spheric.cloud/spheric/iri-api/apis/runtime/v1alpha1"
	"spheric.cloud/spheric/spherelet/controllers/events"
	utilslices "spheric.cloud/spheric/utils/slices"
)

func (r *InstanceReconciler) prepareIRINetworkInterfacesForInstance(
	ctx context.Context,
	instance *corev1alpha1.Instance,
) ([]*iri.NetworkInterface, bool, error) {
	iriNics, err := r.getIRINetworkInterfacesForInstance(ctx, instance)
	if err != nil {
		return nil, false, err
	}

	if len(iriNics) != len(instance.Spec.NetworkInterfaces) {
		expectedNicNames := utilslices.ToSetFunc(instance.Spec.NetworkInterfaces, func(v corev1alpha1.NetworkInterface) string { return v.Name })
		actualNicNames := utilslices.ToSetFunc(iriNics, (*iri.NetworkInterface).GetName)
		missingNicNames := sets.List(expectedNicNames.Difference(actualNicNames))
		r.Eventf(instance, corev1.EventTypeNormal, events.NetworkInterfaceNotReady, "Instance network interfaces are not ready: %v", missingNicNames)
		return nil, false, nil
	}

	return iriNics, true, nil
}

func (r *InstanceReconciler) getIRINetworkInterfacesForInstance(
	ctx context.Context,
	instance *corev1alpha1.Instance,
) ([]*iri.NetworkInterface, error) {
	var (
		iriNics []*iri.NetworkInterface
		errs    []error
	)
	for _, nic := range instance.Spec.NetworkInterfaces {
		iriNic, ok, err := r.prepareIRINetworkInterface(ctx, instance, &nic)
		if err != nil {
			errs = append(errs, fmt.Errorf("[network interface %s] error preparing: %w", nic.Name, err))
			continue
		}
		if !ok {
			continue
		}

		iriNics = append(iriNics, iriNic)
	}

	return iriNics, errors.Join(errs...)
}

func (r *InstanceReconciler) prepareIRINetworkInterface(
	ctx context.Context,
	instance *corev1alpha1.Instance,
	nic *corev1alpha1.NetworkInterface,
) (*iri.NetworkInterface, bool, error) {
	if nic.SubnetRef.Name == "" {
		r.Eventf(instance, events.NetworkInterfaceNotReady, "Network interface %s is not yet assigned to a subnet", nic.Name)
		return nil, false, nil
	}

	subnet := &corev1alpha1.Subnet{}
	subnetKey := client.ObjectKey{Namespace: instance.Namespace, Name: nic.SubnetRef.Name}
	if err := r.Get(ctx, subnetKey, subnet); err != nil {
		if !apierrors.IsNotFound(err) {
			return nil, false, fmt.Errorf("error getting subnet %s: %w", subnetKey.Name, err)
		}

		r.Eventf(instance, corev1.EventTypeNormal, events.NetworkInterfaceNotReady, "Network interface %s subnet %s not found", nic.Name, subnetKey.Name)
		return nil, false, nil
	}

	network := &corev1alpha1.Network{}
	networkKey := client.ObjectKey{Namespace: instance.Namespace, Name: nic.SubnetRef.NetworkName}
	if err := r.Get(ctx, networkKey, network); err != nil {
		if !apierrors.IsNotFound(err) {
			return nil, false, fmt.Errorf("error getting network %s: %w", networkKey.Name, err)
		}

		r.Eventf(instance, corev1.EventTypeNormal, events.NetworkInterfaceNotReady, "Network interface %s network %s not found", nic.Name, networkKey.Name)
		return nil, false, nil
	}

	return &iri.NetworkInterface{
		Name: nic.Name,
		SubnetMetadata: &iri.NetworkInterfaceSubnetMetadata{
			NetworkName: network.Name,
			NetworkUid:  string(network.UID),
			SubnetName:  subnet.Name,
			SubnetUid:   string(subnet.UID),
		},
		Ips:         nic.IPs,
		SubnetCidrs: subnet.Spec.CIDRs,
	}, true, nil
}

func (r *InstanceReconciler) getExistingIRINetworkInterfacesForInstance(
	ctx context.Context,
	log logr.Logger,
	iriInstance *iri.Instance,
	desiredIRINics []*iri.NetworkInterface,
) ([]*iri.NetworkInterface, error) {
	var (
		iriNics              []*iri.NetworkInterface
		desiredIRINicsByName = utilslices.ToMapByKey(desiredIRINics, (*iri.NetworkInterface).GetName)
		errs                 []error
	)

	for _, iriNic := range iriInstance.Spec.NetworkInterfaces {
		log := log.WithValues("NetworkInterface", iriNic.Name)

		desiredIRINic, desiredNicPresent := desiredIRINicsByName[iriNic.Name]
		if desiredNicPresent && proto.Equal(desiredIRINic, iriNic) {
			log.V(1).Info("Existing IRI network interface is up-to-date")
			iriNics = append(iriNics, iriNic)
			continue
		}

		log.V(1).Info("Detaching outdated IRI network interface")
		_, err := r.InstanceRuntime.DetachNetworkInterface(ctx, &iri.DetachNetworkInterfaceRequest{
			InstanceId: iriInstance.Metadata.Id,
			Name:       iriNic.Name,
		})
		if err != nil {
			if status.Code(err) != codes.NotFound {
				errs = append(errs, fmt.Errorf("[network interface %s] %w", iriNic.Name, err))
				continue
			}
		}
	}
	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}
	return iriNics, nil
}

func (r *InstanceReconciler) getNewAttachIRINetworkInterfaces(
	ctx context.Context,
	log logr.Logger,
	iriInstance *iri.Instance,
	desiredIRINics, existingIRINics []*iri.NetworkInterface,
) ([]*iri.NetworkInterface, error) {
	var (
		desiredNewIRINics = FindNewIRINetworkInterfaces(desiredIRINics, existingIRINics)
		iriNics           []*iri.NetworkInterface
		errs              []error
	)
	for _, newIRINic := range desiredNewIRINics {
		log := log.WithValues("NetworkInterface", newIRINic.Name)
		log.V(1).Info("Attaching new network interface")
		if _, err := r.InstanceRuntime.AttachNetworkInterface(ctx, &iri.AttachNetworkInterfaceRequest{
			InstanceId:       iriInstance.Metadata.Id,
			NetworkInterface: newIRINic,
		}); err != nil {
			errs = append(errs, fmt.Errorf("[network interface %s] %w", newIRINic.Name, err))
			continue
		}

		iriNics = append(iriNics, newIRINic)
	}
	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}
	return iriNics, nil
}

func (r *InstanceReconciler) updateIRINetworkInterfaces(
	ctx context.Context,
	log logr.Logger,
	instance *corev1alpha1.Instance,
	iriInstance *iri.Instance,
) error {
	desiredIRINics, err := r.getIRINetworkInterfacesForInstance(ctx, instance)
	if err != nil {
		return fmt.Errorf("error preparing iri network interfaces: %w", err)
	}

	existingIRINics, err := r.getExistingIRINetworkInterfacesForInstance(ctx, log, iriInstance, desiredIRINics)
	if err != nil {
		return fmt.Errorf("error getting existing iri network interfaces for instance: %w", err)
	}

	_, err = r.getNewAttachIRINetworkInterfaces(ctx, log, iriInstance, desiredIRINics, existingIRINics)
	if err != nil {
		return fmt.Errorf("error getting new iri network interfaces for instance: %w", err)
	}

	return nil
}

var iriNetworkInterfaceStateToNetworkInterfaceState = map[iri.NetworkInterfaceState]corev1alpha1.NetworkInterfaceState{
	iri.NetworkInterfaceState_NETWORK_INTERFACE_PENDING:  corev1alpha1.NetworkInterfaceStatePending,
	iri.NetworkInterfaceState_NETWORK_INTERFACE_ATTACHED: corev1alpha1.NetworkInterfaceStateAttached,
}

func (r *InstanceReconciler) convertIRINetworkInterfaceState(state iri.NetworkInterfaceState) (corev1alpha1.NetworkInterfaceState, error) {
	if res, ok := iriNetworkInterfaceStateToNetworkInterfaceState[state]; ok {
		return res, nil
	}
	return "", fmt.Errorf("unknown network interface attachment state %v", state)
}

func (r *InstanceReconciler) convertIRINetworkInterfaceStatus(status *iri.NetworkInterfaceStatus) (corev1alpha1.NetworkInterfaceStatus, error) {
	state, err := r.convertIRINetworkInterfaceState(status.State)
	if err != nil {
		return corev1alpha1.NetworkInterfaceStatus{}, err
	}

	return corev1alpha1.NetworkInterfaceStatus{
		Name:  status.Name,
		State: state,
	}, nil
}

func (r *InstanceReconciler) addNetworkInterfaceStatusValues(now metav1.Time, existing, newValues *corev1alpha1.NetworkInterfaceStatus) {
	if existing.State != newValues.State {
		existing.LastStateTransitionTime = &now
	}
	existing.Name = newValues.Name
	existing.State = newValues.State
}

func (r *InstanceReconciler) getNetworkInterfaceStatusesForInstance(
	instance *corev1alpha1.Instance,
	iriInstance *iri.Instance,
	now metav1.Time,
) ([]corev1alpha1.NetworkInterfaceStatus, error) {
	var (
		iriNicStatusByName        = utilslices.ToMapByKey(iriInstance.Status.NetworkInterfaces, (*iri.NetworkInterfaceStatus).GetName)
		existingNicStatusesByName = utilslices.ToMapByKey(instance.Status.NetworkInterfaces, func(status corev1alpha1.NetworkInterfaceStatus) string { return status.Name })
		nicStatuses               []corev1alpha1.NetworkInterfaceStatus
		errs                      []error
	)

	for _, instanceNic := range instance.Spec.NetworkInterfaces {
		var (
			iriNicStatus, ok = iriNicStatusByName[instanceNic.Name]
			nicStatusValues  corev1alpha1.NetworkInterfaceStatus
		)
		if ok {
			var err error
			nicStatusValues, err = r.convertIRINetworkInterfaceStatus(iriNicStatus)
			if err != nil {
				return nil, fmt.Errorf("[network interface %s] %w", instanceNic.Name, err)
			}
		} else {
			nicStatusValues = corev1alpha1.NetworkInterfaceStatus{
				Name:  instanceNic.Name,
				State: corev1alpha1.NetworkInterfaceStatePending,
			}
		}

		nicStatus := existingNicStatusesByName[instanceNic.Name]
		r.addNetworkInterfaceStatusValues(now, &nicStatus, &nicStatusValues)
		nicStatuses = append(nicStatuses, nicStatus)
	}
	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}
	return nicStatuses, nil
}
