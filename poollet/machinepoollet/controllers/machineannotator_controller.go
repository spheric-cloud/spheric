// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package controllers

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/source"
	computev1alpha1 "spheric.cloud/spheric/api/compute/v1alpha1"
	machinepoolletv1alpha1 "spheric.cloud/spheric/poollet/machinepoollet/api/v1alpha1"
	"spheric.cloud/spheric/poollet/srievent"
	sri "spheric.cloud/spheric/sri/apis/machine/v1alpha1"
	srimeta "spheric.cloud/spheric/sri/apis/meta/v1alpha1"
	sphericclient "spheric.cloud/spheric/utils/client"
)

type MachineAnnotatorReconciler struct {
	client.Client

	MachineEvents srievent.Source[*sri.Machine]
}

func (r *MachineAnnotatorReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	machine := &computev1alpha1.Machine{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: req.Namespace,
			Name:      req.Name,
		},
	}

	if err := sphericclient.PatchAddReconcileAnnotation(ctx, r.Client, machine); client.IgnoreNotFound(err) != nil {
		return ctrl.Result{}, fmt.Errorf("error patching machine: %w", err)
	}
	return ctrl.Result{}, nil
}

func machineAnnotatorEventHandler[O srimeta.Object](log logr.Logger, c chan<- event.GenericEvent) srievent.HandlerFuncs[O] {
	handleEvent := func(obj srimeta.Object) {
		namespace, ok := obj.GetMetadata().Labels[machinepoolletv1alpha1.MachineNamespaceLabel]
		if !ok {
			return
		}

		name, ok := obj.GetMetadata().Labels[machinepoolletv1alpha1.MachineNameLabel]
		if !ok {
			return
		}

		machine := &computev1alpha1.Machine{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: namespace,
				Name:      name,
			},
		}

		select {
		case c <- event.GenericEvent{Object: machine}:
		default:
			log.V(5).Info("Channel full, discarding event")
		}
	}

	return srievent.HandlerFuncs[O]{
		CreateFunc: func(event srievent.CreateEvent[O]) {
			handleEvent(event.Object)
		},
		UpdateFunc: func(event srievent.UpdateEvent[O]) {
			handleEvent(event.ObjectNew)
		},
		DeleteFunc: func(event srievent.DeleteEvent[O]) {
			handleEvent(event.Object)
		},
		GenericFunc: func(event srievent.GenericEvent[O]) {
			handleEvent(event.Object)
		},
	}
}

func (r *MachineAnnotatorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	c, err := controller.New("machineannotator", mgr, controller.Options{
		Reconciler: r,
	})
	if err != nil {
		return err
	}

	src, err := r.sriMachineEventSource(mgr)
	if err != nil {
		return err
	}

	if err := c.Watch(src, &handler.EnqueueRequestForObject{}); err != nil {
		return err
	}

	return nil
}

func (r *MachineAnnotatorReconciler) sriMachineEventSource(mgr ctrl.Manager) (source.Source, error) {
	ch := make(chan event.GenericEvent, 1024)

	if err := mgr.Add(manager.RunnableFunc(func(ctx context.Context) error {
		log := ctrl.LoggerFrom(ctx).WithName("machineannotator").WithName("srieventhandlers")

		registrationFuncs := []func() (srievent.HandlerRegistration, error){
			func() (srievent.HandlerRegistration, error) {
				return r.MachineEvents.AddHandler(machineAnnotatorEventHandler[*sri.Machine](log, ch))
			},
		}

		var handles []srievent.HandlerRegistration
		defer func() {
			log.V(1).Info("Removing handles")
			for _, handle := range handles {
				if err := handle.Remove(); err != nil {
					log.Error(err, "Error removing handle")
				}
			}
		}()

		for _, registrationFunc := range registrationFuncs {
			handle, err := registrationFunc()
			if err != nil {
				return err
			}

			handles = append(handles, handle)
		}

		<-ctx.Done()
		return nil
	})); err != nil {
		return nil, err
	}

	return &source.Channel{Source: ch}, nil
}
