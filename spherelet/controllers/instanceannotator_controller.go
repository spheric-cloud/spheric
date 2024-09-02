// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package controllers

import (
	"context"
	"fmt"

	"sigs.k8s.io/controller-runtime/pkg/source"

	sphereletevent "spheric.cloud/spheric/spherelet/event/instanceevent"

	"github.com/go-logr/logr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	corev1alpha1 "spheric.cloud/spheric/api/core/v1alpha1"
	iri "spheric.cloud/spheric/iri-api/apis/runtime/v1alpha1"
	sphereletv1alpha1 "spheric.cloud/spheric/spherelet/api/v1alpha1"
	sphericclient "spheric.cloud/spheric/utils/client"
)

type InstanceAnnotatorReconciler struct {
	client.Client

	InstanceEvents sphereletevent.Source
}

func (r *InstanceAnnotatorReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	instance := &corev1alpha1.Instance{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: req.Namespace,
			Name:      req.Name,
		},
	}

	if err := sphericclient.PatchAddReconcileAnnotation(ctx, r.Client, instance); client.IgnoreNotFound(err) != nil {
		return ctrl.Result{}, fmt.Errorf("error patching instance: %w", err)
	}
	return ctrl.Result{}, nil
}

func instanceAnnotatorEventHandler(log logr.Logger, c chan<- event.GenericEvent) sphereletevent.HandlerFuncs {
	handleEvent := func(obj *iri.Instance) {
		namespace, ok := obj.GetMetadata().Labels[sphereletv1alpha1.InstanceNamespaceLabel]
		if !ok {
			return
		}

		name, ok := obj.GetMetadata().Labels[sphereletv1alpha1.InstanceNameLabel]
		if !ok {
			return
		}

		instance := &corev1alpha1.Instance{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: namespace,
				Name:      name,
			},
		}

		select {
		case c <- event.GenericEvent{Object: instance}:
		default:
			log.V(5).Info("Channel full, discarding event")
		}
	}

	return sphereletevent.HandlerFuncs{
		CreateFunc: func(event sphereletevent.CreateEvent) {
			handleEvent(event.Object)
		},
		UpdateFunc: func(event sphereletevent.UpdateEvent) {
			handleEvent(event.ObjectNew)
		},
		DeleteFunc: func(event sphereletevent.DeleteEvent) {
			handleEvent(event.Object)
		},
		GenericFunc: func(event sphereletevent.GenericEvent) {
			handleEvent(event.Object)
		},
	}
}

func (r *InstanceAnnotatorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	instanceEventChannel, err := r.iriInstanceEventChannel(mgr)
	if err != nil {
		return err
	}

	return ctrl.NewControllerManagedBy(mgr).
		Named("instanceannotator").
		WatchesRawSource(source.Channel(instanceEventChannel, &handler.EnqueueRequestForObject{})).
		Complete(r)
}

func (r *InstanceAnnotatorReconciler) iriInstanceEventChannel(mgr ctrl.Manager) (<-chan event.GenericEvent, error) {
	ch := make(chan event.GenericEvent, 1024)

	if err := mgr.Add(manager.RunnableFunc(func(ctx context.Context) error {
		log := ctrl.LoggerFrom(ctx).WithName("instanceannotator").WithName("irieventhandlers")

		registrationFuncs := []func() (sphereletevent.HandlerRegistration, error){
			func() (sphereletevent.HandlerRegistration, error) {
				return r.InstanceEvents.AddHandler(instanceAnnotatorEventHandler(log, ch))
			},
		}

		var handles []sphereletevent.HandlerRegistration
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
