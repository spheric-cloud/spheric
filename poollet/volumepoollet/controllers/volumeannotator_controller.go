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
	storagev1alpha1 "spheric.cloud/spheric/api/storage/v1alpha1"
	"spheric.cloud/spheric/poollet/srievent"
	volumepoolletv1alpha1 "spheric.cloud/spheric/poollet/volumepoollet/api/v1alpha1"
	srimeta "spheric.cloud/spheric/sri/apis/meta/v1alpha1"
	sri "spheric.cloud/spheric/sri/apis/volume/v1alpha1"
	sphericclient "spheric.cloud/spheric/utils/client"
)

type VolumeAnnotatorReconciler struct {
	client.Client

	VolumeEvents srievent.Source[*sri.Volume]
}

func (r *VolumeAnnotatorReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	volume := &storagev1alpha1.Volume{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: req.Namespace,
			Name:      req.Name,
		},
	}

	if err := sphericclient.PatchAddReconcileAnnotation(ctx, r.Client, volume); client.IgnoreNotFound(err) != nil {
		return ctrl.Result{}, fmt.Errorf("error patching volume: %w", err)
	}
	return ctrl.Result{}, nil
}

func volumeAnnotatorEventHandler[O srimeta.Object](log logr.Logger, c chan<- event.GenericEvent) srievent.HandlerFuncs[O] {
	handleEvent := func(obj srimeta.Object) {
		namespace, ok := obj.GetMetadata().Labels[volumepoolletv1alpha1.VolumeNamespaceLabel]
		if !ok {
			return
		}

		name, ok := obj.GetMetadata().Labels[volumepoolletv1alpha1.VolumeNameLabel]
		if !ok {
			return
		}

		volume := &storagev1alpha1.Volume{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: namespace,
				Name:      name,
			},
		}

		select {
		case c <- event.GenericEvent{Object: volume}:
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

func (r *VolumeAnnotatorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	c, err := controller.New("volumeannotator", mgr, controller.Options{
		Reconciler: r,
	})
	if err != nil {
		return err
	}

	src, err := r.sriVolumeEventSource(mgr)
	if err != nil {
		return err
	}

	if err := c.Watch(src, &handler.EnqueueRequestForObject{}); err != nil {
		return err
	}

	return nil
}

func (r *VolumeAnnotatorReconciler) sriVolumeEventSource(mgr ctrl.Manager) (source.Source, error) {
	ch := make(chan event.GenericEvent, 1024)

	if err := mgr.Add(manager.RunnableFunc(func(ctx context.Context) error {
		log := ctrl.LoggerFrom(ctx).WithName("volumeannotator").WithName("srieventhandlers")

		registrationFuncs := []func() (srievent.HandlerRegistration, error){
			func() (srievent.HandlerRegistration, error) {
				return r.VolumeEvents.AddHandler(volumeAnnotatorEventHandler[*sri.Volume](log, ch))
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
