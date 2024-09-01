// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package debug

import (
	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

type typedDebugPredicate[object any] struct {
	log         logr.Logger
	predicate   predicate.TypedPredicate[object]
	objectValue func(object) any
}

func (d *typedDebugPredicate[object]) Create(evt event.TypedCreateEvent[object]) bool {
	log := d.log.WithValues("Event", "Create", "Object", d.objectValue(evt.Object))
	log.Info("Handling event")
	res := d.predicate.Create(evt)
	log.Info("Handled event", "Result", res)
	return res
}

func (d *typedDebugPredicate[object]) Delete(evt event.TypedDeleteEvent[object]) bool {
	log := d.log.WithValues("Event", "Delete", "Object", d.objectValue(evt.Object))
	log.Info("Handling event")
	res := d.predicate.Delete(evt)
	log.Info("Handled event", "Result", res)
	return res
}

func (d *typedDebugPredicate[object]) Update(evt event.TypedUpdateEvent[object]) bool {
	log := d.log.WithValues("Event", "Update", "ObjectOld", d.objectValue(evt.ObjectOld), "ObjectNew", d.objectValue(evt.ObjectNew))
	log.Info("Handling event")
	res := d.predicate.Update(evt)
	log.Info("Handled event", "Result", res)
	return res
}

func (d *typedDebugPredicate[object]) Generic(evt event.TypedGenericEvent[object]) bool {
	log := d.log.WithValues("Event", "Generic", "Object", d.objectValue(evt.Object))
	log.Info("Handling event")
	res := d.predicate.Generic(evt)
	log.Info("Handled event", "Result", res)
	return res
}

// TypedPredicate allows debugging a predicate.Predicate by wrapping it and logging each action it does.
//
// Caution: This has a heavy toll on runtime performance and should *not* be used in production code.
// Use only for debugging predicates and remove once done.
func TypedPredicate[object any](name string, prct predicate.TypedPredicate[object], opts ...TypedPredicateOption[object]) predicate.TypedPredicate[object] {
	o := (&TypedPredicateOptions[object]{}).ApplyOptions(opts)
	setPredicateOptionsDefaults(o)

	return &typedDebugPredicate[object]{
		log:         o.Log.WithName(name),
		predicate:   prct,
		objectValue: o.ObjectValue,
	}
}
