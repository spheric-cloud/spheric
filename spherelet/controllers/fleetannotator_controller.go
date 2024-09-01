// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package controllers

import (
	"context"
	"fmt"

	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/source"
	"spheric.cloud/spheric/spherelet/event/runtimeevent"

	"github.com/go-logr/logr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	corev1alpha1 "spheric.cloud/spheric/api/core/v1alpha1"
	sphericclient "spheric.cloud/spheric/utils/client"
)

type FleetAnnotatorReconciler struct {
	client.Client

	FleetName     string
	RuntimeEvents runtimeevent.Source
}

func (r *FleetAnnotatorReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	machinePool := &corev1alpha1.Fleet{
		ObjectMeta: metav1.ObjectMeta{
			Name: req.Name,
		},
	}

	if err := sphericclient.PatchAddReconcileAnnotation(ctx, r.Client, machinePool); client.IgnoreNotFound(err) != nil {
		return ctrl.Result{}, fmt.Errorf("error patching fleet: %w", err)
	}
	return ctrl.Result{}, nil
}

func (r *FleetAnnotatorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	runtimeEventsChannel, err := r.runtimeEventChannel(mgr)
	if err != nil {
		return err
	}

	return ctrl.NewControllerManagedBy(mgr).
		Named("fleetannotator").
		WatchesRawSource(source.Channel(runtimeEventsChannel, &handler.EnqueueRequestForObject{})).
		Complete(r)
}

func (r *FleetAnnotatorReconciler) runtimeEventChannel(mgr ctrl.Manager) (<-chan event.GenericEvent, error) {
	ch := make(chan event.GenericEvent, 1024)

	if err := mgr.Add(manager.RunnableFunc(func(ctx context.Context) error {
		log := ctrl.LoggerFrom(ctx).WithName("fleet").WithName("irieventhandlers")

		registrationFuncs := []func() (runtimeevent.HandlerRegistration, error){
			func() (runtimeevent.HandlerRegistration, error) {
				return r.RuntimeEvents.AddHandler(r.fleetAnnotatorEventHandler(log, ch))
			},
		}

		var handles []runtimeevent.HandlerRegistration
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

	return ch, nil
}

func (r *FleetAnnotatorReconciler) fleetAnnotatorEventHandler(log logr.Logger, c chan<- event.GenericEvent) runtimeevent.Handler {
	handleEvent := func() {
		select {
		case c <- event.GenericEvent{Object: &corev1alpha1.Fleet{ObjectMeta: metav1.ObjectMeta{
			Name: r.FleetName,
		}}}:
			log.V(1).Info("Added item to queue")
		default:
			log.V(5).Info("Channel full, discarding event")
		}
	}

	return runtimeevent.HandlerFuncs{
		UpdateResourcesFunc: func(*runtimeevent.UpdateResourcesEvent) {
			handleEvent()
		},
	}
}
