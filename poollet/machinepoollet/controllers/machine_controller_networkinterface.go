// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package controllers

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/gogo/protobuf/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	"sigs.k8s.io/controller-runtime/pkg/client"
	commonv1alpha1 "spheric.cloud/spheric/api/common/v1alpha1"
	computev1alpha1 "spheric.cloud/spheric/api/compute/v1alpha1"
	ipamv1alpha1 "spheric.cloud/spheric/api/ipam/v1alpha1"
	networkingv1alpha1 "spheric.cloud/spheric/api/networking/v1alpha1"
	"spheric.cloud/spheric/poollet/machinepoollet/api/v1alpha1"
	"spheric.cloud/spheric/poollet/machinepoollet/controllers/events"
	sri "spheric.cloud/spheric/sri/apis/machine/v1alpha1"
	"spheric.cloud/spheric/utils/claimmanager"
	utilslices "spheric.cloud/spheric/utils/slices"
)

type networkInterfaceClaimStrategy struct {
	client.Client
}

func (s *networkInterfaceClaimStrategy) ClaimState(claimer client.Object, obj client.Object) claimmanager.ClaimState {
	nic := obj.(*networkingv1alpha1.NetworkInterface)
	if machineRef := nic.Spec.MachineRef; machineRef != nil {
		if machineRef.UID == claimer.GetUID() {
			return claimmanager.ClaimStateClaimed
		}
		return claimmanager.ClaimStateTaken
	}
	return claimmanager.ClaimStateFree
}

func (s *networkInterfaceClaimStrategy) Adopt(ctx context.Context, claimer client.Object, obj client.Object) error {
	nic := obj.(*networkingv1alpha1.NetworkInterface)
	base := nic.DeepCopy()
	nic.Spec.MachineRef = commonv1alpha1.NewLocalObjUIDRef(claimer)
	nic.Spec.ProviderID = ""
	return s.Patch(ctx, nic, client.StrategicMergeFrom(base))
}

func (s *networkInterfaceClaimStrategy) Release(ctx context.Context, claimer client.Object, obj client.Object) error {
	nic := obj.(*networkingv1alpha1.NetworkInterface)
	base := nic.DeepCopy()
	nic.Spec.ProviderID = ""
	nic.Spec.MachineRef = nil
	return s.Patch(ctx, nic, client.StrategicMergeFrom(base))
}

func (r *MachineReconciler) networkInterfaceNameToMachineNetworkInterfaceName(machine *computev1alpha1.Machine) map[string]string {
	sel := make(map[string]string)
	for _, machineNic := range machine.Spec.NetworkInterfaces {
		nicName := computev1alpha1.MachineNetworkInterfaceName(machine.Name, machineNic)
		sel[nicName] = machineNic.Name
	}
	return sel
}

func (r *MachineReconciler) machineNetworkInterfaceSelector(machine *computev1alpha1.Machine) claimmanager.Selector {
	names := sets.New(computev1alpha1.MachineNetworkInterfaceNames(machine)...)
	return claimmanager.SelectorFunc(func(obj client.Object) bool {
		nic := obj.(*networkingv1alpha1.NetworkInterface)
		return names.Has(nic.Name)
	})
}

func (r *MachineReconciler) getNetworkInterfacesForMachine(ctx context.Context, machine *computev1alpha1.Machine) ([]networkingv1alpha1.NetworkInterface, error) {
	nicList := &networkingv1alpha1.NetworkInterfaceList{}
	if err := r.List(ctx, nicList,
		client.InNamespace(machine.Namespace),
	); err != nil {
		return nil, fmt.Errorf("error listing network interfaces: %w", err)
	}

	var (
		sel      = r.machineNetworkInterfaceSelector(machine)
		claimMgr = claimmanager.New(machine, sel, &networkInterfaceClaimStrategy{r.Client})
		nics     []networkingv1alpha1.NetworkInterface
		errs     []error
	)
	for _, nic := range nicList.Items {
		ok, err := claimMgr.Claim(ctx, &nic)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		if !ok {
			continue
		}

		nics = append(nics, nic)
	}
	return nics, errors.Join(errs...)
}

func (r *MachineReconciler) prepareSRINetworkInterfacesForMachine(
	ctx context.Context,
	machine *computev1alpha1.Machine,
	nics []networkingv1alpha1.NetworkInterface,
) ([]*sri.NetworkInterface, map[string]v1alpha1.ObjectUIDRef, bool, error) {
	sriNics, mapping, err := r.getSRINetworkInterfacesForMachine(ctx, machine, nics)
	if err != nil {
		return nil, nil, false, err
	}

	if len(sriNics) != len(machine.Spec.Volumes) {
		expectedNicNames := utilslices.ToSetFunc(machine.Spec.NetworkInterfaces, func(v computev1alpha1.NetworkInterface) string { return v.Name })
		actualNicNames := utilslices.ToSetFunc(sriNics, (*sri.NetworkInterface).GetName)
		missingNicNames := sets.List(expectedNicNames.Difference(actualNicNames))
		r.Eventf(machine, corev1.EventTypeNormal, events.NetworkInterfaceNotReady, "Machine network interfaces are not ready: %v", missingNicNames)
		return nil, nil, false, nil
	}

	return sriNics, mapping, true, err
}

func (r *MachineReconciler) getSRINetworkInterfacesForMachine(
	ctx context.Context,
	machine *computev1alpha1.Machine,
	nics []networkingv1alpha1.NetworkInterface,
) ([]*sri.NetworkInterface, map[string]v1alpha1.ObjectUIDRef, error) {
	var (
		nicNameToMachineNicName = r.networkInterfaceNameToMachineNetworkInterfaceName(machine)

		sriNics                []*sri.NetworkInterface
		machineNicNameToUIDRef = make(map[string]v1alpha1.ObjectUIDRef)
		errs                   []error
	)
	for _, nic := range nics {
		machineNicName := nicNameToMachineNicName[nic.Name]
		sriNic, ok, err := r.prepareSRINetworkInterface(ctx, machine, &nic, machineNicName)
		if err != nil {
			errs = append(errs, fmt.Errorf("[network interface %s] error preparing: %w", machineNicName, err))
			continue
		}
		if !ok {
			continue
		}

		sriNics = append(sriNics, sriNic)
		machineNicNameToUIDRef[machineNicName] = v1alpha1.ObjUID(&nic)
	}
	if err := errors.Join(errs...); err != nil {
		return nil, nil, err
	}
	return sriNics, machineNicNameToUIDRef, nil
}

func (r *MachineReconciler) getNetworkInterfaceIP(
	ctx context.Context,
	machine *computev1alpha1.Machine,
	nic *networkingv1alpha1.NetworkInterface,
	idx int,
	nicIP networkingv1alpha1.IPSource,
) (commonv1alpha1.IP, bool, error) {
	switch {
	case nicIP.Value != nil:
		return *nicIP.Value, true, nil
	case nicIP.Ephemeral != nil:
		prefix := &ipamv1alpha1.Prefix{}
		prefixName := networkingv1alpha1.NetworkInterfaceIPIPAMPrefixName(nic.Name, idx)
		prefixKey := client.ObjectKey{Namespace: nic.Namespace, Name: prefixName}
		if err := r.Get(ctx, prefixKey, prefix); err != nil {
			if !apierrors.IsNotFound(err) {
				return commonv1alpha1.IP{}, false, fmt.Errorf("error getting prefix %s: %w", prefixName, err)
			}

			r.Eventf(machine, corev1.EventTypeNormal, events.NetworkInterfaceNotReady, "Network interface prefix %s not found", prefixName)
			return commonv1alpha1.IP{}, false, nil
		}

		if !metav1.IsControlledBy(prefix, nic) {
			r.Eventf(machine, corev1.EventTypeNormal, events.NetworkInterfaceNotReady, "Network interface prefix %s not controlled by network interface", prefixName, nic.Name)
			return commonv1alpha1.IP{}, false, nil
		}

		if prefix.Status.Phase != ipamv1alpha1.PrefixPhaseAllocated {
			r.Eventf(machine, corev1.EventTypeNormal, events.NetworkInterfaceNotReady, "Network interface prefix %s is not yet allocated", prefixName)
			return commonv1alpha1.IP{}, false, nil
		}

		return prefix.Spec.Prefix.IP(), true, nil
	default:
		return commonv1alpha1.IP{}, false, fmt.Errorf("unrecognized network interface ip %#v", nicIP)
	}
}

func (r *MachineReconciler) getNetworkInterfaceIPs(
	ctx context.Context,
	machine *computev1alpha1.Machine,
	nic *networkingv1alpha1.NetworkInterface,
) ([]commonv1alpha1.IP, bool, error) {
	var ips []commonv1alpha1.IP
	for i, nicIP := range nic.Spec.IPs {
		ip, ok, err := r.getNetworkInterfaceIP(ctx, machine, nic, i, nicIP)
		if err != nil || !ok {
			return nil, false, err
		}

		ips = append(ips, ip)
	}
	return ips, true, nil
}

func (r *MachineReconciler) prepareSRINetworkInterface(
	ctx context.Context,
	machine *computev1alpha1.Machine,
	nic *networkingv1alpha1.NetworkInterface,
	machineNicName string,
) (*sri.NetworkInterface, bool, error) {
	network := &networkingv1alpha1.Network{}
	networkKey := client.ObjectKey{Namespace: nic.Namespace, Name: nic.Spec.NetworkRef.Name}
	if err := r.Get(ctx, networkKey, network); err != nil {
		if !apierrors.IsNotFound(err) {
			return nil, false, fmt.Errorf("error getting network %s: %w", networkKey.Name, err)
		}
		r.Eventf(machine, corev1.EventTypeNormal, events.NetworkInterfaceNotReady, "Network interface %s network %s not found", nic.Name, networkKey.Name)
		return nil, false, nil
	}

	ips, ok, err := r.getNetworkInterfaceIPs(ctx, machine, nic)
	if err != nil || !ok {
		return nil, false, err
	}

	return &sri.NetworkInterface{
		Name:       machineNicName,
		NetworkId:  network.Spec.ProviderID,
		Ips:        utilslices.Map(ips, commonv1alpha1.IP.String),
		Attributes: nic.Spec.Attributes,
	}, true, nil
}

func (r *MachineReconciler) getExistingSRINetworkInterfacesForMachine(
	ctx context.Context,
	log logr.Logger,
	sriMachine *sri.Machine,
	desiredSRINics []*sri.NetworkInterface,
) ([]*sri.NetworkInterface, error) {
	var (
		sriNics              []*sri.NetworkInterface
		desiredSRINicsByName = utilslices.ToMapByKey(desiredSRINics, (*sri.NetworkInterface).GetName)
		errs                 []error
	)

	for _, sriNic := range sriMachine.Spec.NetworkInterfaces {
		log := log.WithValues("NetworkInterface", sriNic.Name)

		desiredSRINic, desiredNicPresent := desiredSRINicsByName[sriNic.Name]
		if desiredNicPresent && proto.Equal(desiredSRINic, sriNic) {
			log.V(1).Info("Existing SRI network interface is up-to-date")
			sriNics = append(sriNics, sriNic)
			continue
		}

		log.V(1).Info("Detaching outdated SRI network interface")
		_, err := r.MachineRuntime.DetachNetworkInterface(ctx, &sri.DetachNetworkInterfaceRequest{
			MachineId: sriMachine.Metadata.Id,
			Name:      sriNic.Name,
		})
		if err != nil {
			if status.Code(err) != codes.NotFound {
				errs = append(errs, fmt.Errorf("[network interface %s] %w", sriNic.Name, err))
				continue
			}
		}
	}
	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}
	return sriNics, nil
}

func (r *MachineReconciler) getNewAttachSRINetworkInterfaces(
	ctx context.Context,
	log logr.Logger,
	sriMachine *sri.Machine,
	desiredSRINics, existingSRINics []*sri.NetworkInterface,
) ([]*sri.NetworkInterface, error) {
	var (
		desiredNewSRINics = FindNewSRINetworkInterfaces(desiredSRINics, existingSRINics)
		sriNics           []*sri.NetworkInterface
		errs              []error
	)
	for _, newSRINic := range desiredNewSRINics {
		log := log.WithValues("NetworkInterface", newSRINic.Name)
		log.V(1).Info("Attaching new network interface")
		if _, err := r.MachineRuntime.AttachNetworkInterface(ctx, &sri.AttachNetworkInterfaceRequest{
			MachineId:        sriMachine.Metadata.Id,
			NetworkInterface: newSRINic,
		}); err != nil {
			errs = append(errs, fmt.Errorf("[network interface %s] %w", newSRINic.Name, err))
			continue
		}

		sriNics = append(sriNics, newSRINic)
	}
	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}
	return sriNics, nil
}

func (r *MachineReconciler) updateSRINetworkInterfaces(
	ctx context.Context,
	log logr.Logger,
	machine *computev1alpha1.Machine,
	sriMachine *sri.Machine,
	nics []networkingv1alpha1.NetworkInterface,
) ([]*sri.NetworkInterface, error) {
	desiredSRINics, _, err := r.getSRINetworkInterfacesForMachine(ctx, machine, nics)
	if err != nil {
		return nil, fmt.Errorf("error preparing sri network interfaces: %w", err)
	}

	existingSRINics, err := r.getExistingSRINetworkInterfacesForMachine(ctx, log, sriMachine, desiredSRINics)
	if err != nil {
		return nil, fmt.Errorf("error getting existing sri network interfaces for machine: %w", err)
	}

	_, err = r.getNewAttachSRINetworkInterfaces(ctx, log, sriMachine, desiredSRINics, existingSRINics)
	if err != nil {
		return nil, fmt.Errorf("error getting new sri network interfaces for machine: %w", err)
	}

	return desiredSRINics, nil
}

func (r *MachineReconciler) computeNetworkInterfaceMapping(
	machine *computev1alpha1.Machine,
	nics []networkingv1alpha1.NetworkInterface,
	sriNics []*sri.NetworkInterface,
) map[string]v1alpha1.ObjectUIDRef {
	nicByName := utilslices.ToMapByKey(nics,
		func(nic networkingv1alpha1.NetworkInterface) string { return nic.Name },
	)

	machineNicNameToNicName := make(map[string]string, len(machine.Spec.NetworkInterfaces))
	for _, machineNic := range machine.Spec.NetworkInterfaces {
		nicName := computev1alpha1.MachineNetworkInterfaceName(machine.Name, machineNic)
		machineNicNameToNicName[machineNic.Name] = nicName
	}

	nicMapping := make(map[string]v1alpha1.ObjectUIDRef, len(sriNics))
	for _, sriNic := range sriNics {
		nicName := machineNicNameToNicName[sriNic.Name]
		nic := nicByName[nicName]

		nicMapping[sriNic.Name] = v1alpha1.ObjUID(&nic)
	}
	return nicMapping
}

var sriNetworkInterfaceStateToNetworkInterfaceState = map[sri.NetworkInterfaceState]computev1alpha1.NetworkInterfaceState{
	sri.NetworkInterfaceState_NETWORK_INTERFACE_PENDING:  computev1alpha1.NetworkInterfaceStatePending,
	sri.NetworkInterfaceState_NETWORK_INTERFACE_ATTACHED: computev1alpha1.NetworkInterfaceStateAttached,
}

func (r *MachineReconciler) convertSRINetworkInterfaceState(state sri.NetworkInterfaceState) (computev1alpha1.NetworkInterfaceState, error) {
	if res, ok := sriNetworkInterfaceStateToNetworkInterfaceState[state]; ok {
		return res, nil
	}
	return "", fmt.Errorf("unknown network interface attachment state %v", state)
}

func (r *MachineReconciler) convertSRINetworkInterfaceStatus(status *sri.NetworkInterfaceStatus) (computev1alpha1.NetworkInterfaceStatus, error) {
	state, err := r.convertSRINetworkInterfaceState(status.State)
	if err != nil {
		return computev1alpha1.NetworkInterfaceStatus{}, err
	}

	return computev1alpha1.NetworkInterfaceStatus{
		Name:   status.Name,
		Handle: status.Handle,
		State:  state,
	}, nil
}

func (r *MachineReconciler) addNetworkInterfaceStatusValues(now metav1.Time, existing, newValues *computev1alpha1.NetworkInterfaceStatus) {
	if existing.State != newValues.State {
		existing.LastStateTransitionTime = &now
	}
	existing.Name = newValues.Name
	existing.State = newValues.State
	existing.Handle = newValues.Handle
}

func (r *MachineReconciler) getNetworkInterfaceStatusesForMachine(
	machine *computev1alpha1.Machine,
	sriMachine *sri.Machine,
	now metav1.Time,
) ([]computev1alpha1.NetworkInterfaceStatus, error) {
	var (
		sriNicStatusByName        = utilslices.ToMapByKey(sriMachine.Status.NetworkInterfaces, (*sri.NetworkInterfaceStatus).GetName)
		existingNicStatusesByName = utilslices.ToMapByKey(machine.Status.NetworkInterfaces, func(status computev1alpha1.NetworkInterfaceStatus) string { return status.Name })
		nicStatuses               []computev1alpha1.NetworkInterfaceStatus
		errs                      []error
	)

	for _, machineNic := range machine.Spec.NetworkInterfaces {
		var (
			sriNicStatus, ok = sriNicStatusByName[machineNic.Name]
			nicStatusValues  computev1alpha1.NetworkInterfaceStatus
		)
		if ok {
			var err error
			nicStatusValues, err = r.convertSRINetworkInterfaceStatus(sriNicStatus)
			if err != nil {
				return nil, fmt.Errorf("[network interface %s] %w", machineNic.Name, err)
			}
		} else {
			nicStatusValues = computev1alpha1.NetworkInterfaceStatus{
				Name:  machineNic.Name,
				State: computev1alpha1.NetworkInterfaceStatePending,
			}
		}

		nicStatus := existingNicStatusesByName[machineNic.Name]
		r.addNetworkInterfaceStatusValues(now, &nicStatus, &nicStatusValues)
		nicStatuses = append(nicStatuses, nicStatus)
	}
	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}
	return nicStatuses, nil
}
