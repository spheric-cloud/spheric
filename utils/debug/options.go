// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package debug

import (
	"github.com/go-logr/logr"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// TypedHandlerOptions are options for construction a debug handler.
type TypedHandlerOptions[object any, request comparable] struct {
	// Log is the logger to use. If unspecified, the debug package logger will be used.
	Log logr.Logger

	// ObjectValue controls how an object will be represented as in the log values.
	ObjectValue func(object) any
}

func (o *TypedHandlerOptions[object, request]) ApplyOptions(opts []TypedHandlerOption[object, request]) *TypedHandlerOptions[object, request] {
	for _, opt := range opts {
		opt.ApplyToHandler(o)
	}
	return o
}

func setTypedHandlerOptionsDefaults[object any, request comparable](o *TypedHandlerOptions[object, request]) {
	if o.Log.GetSink() == nil {
		o.Log = handlerLog
	}
	if o.ObjectValue == nil {
		o.ObjectValue = DefaultObjectValue[object]()
	}
}

// TypedPredicateOptions are options for construction a debug predicate.
type TypedPredicateOptions[object any] struct {
	// Log is the logger to use. If unspecified, the debug package logger will be used.
	Log logr.Logger

	// ObjectValue controls how an object will be represented as in the log values.
	ObjectValue func(object) any
}

func (o *TypedPredicateOptions[object]) ApplyToPredicate(o2 *TypedPredicateOptions[object]) {
	if o.Log.GetSink() != nil {
		o2.Log = o.Log
	}
	if o.ObjectValue != nil {
		o2.ObjectValue = o.ObjectValue
	}
}

func (o *TypedPredicateOptions[object]) ApplyOptions(opts []TypedPredicateOption[object]) *TypedPredicateOptions[object] {
	for _, opt := range opts {
		opt.ApplyToPredicate(o)
	}
	return o
}

func setPredicateOptionsDefaults[object any](o *TypedPredicateOptions[object]) {
	if o.Log.GetSink() == nil {
		o.Log = predicateLog
	}
	if o.ObjectValue == nil {
		o.ObjectValue = DefaultObjectValue[object]()
	}
}

type TypedPredicateOption[object any] interface {
	ApplyToPredicate(o *TypedPredicateOptions[object])
}

// DefaultObjectValue provides object logging values by using klog.KObj.
func DefaultObjectValue[object any]() func(object) any {
	var obj object
	if _, ok := any(obj).(client.Object); ok {
		return func(obj object) any {
			return klog.KObj(any(obj).(client.Object))
		}
	}
	return func(obj object) any {
		return obj
	}
}

type TypedHandlerOption[object any, request comparable] interface {
	ApplyToHandler(o *TypedHandlerOptions[object, request])
}

// WithLog specifies the logger to use.
type WithLog[object any, request comparable] struct {
	Log logr.Logger
}

func (w WithLog[object, request]) ApplyToHandler(o *TypedHandlerOptions[object, request]) {
	o.Log = w.Log
}

func (w WithLog[object, request]) ApplyToPredicate(o *TypedPredicateOptions[object]) {
	o.Log = w.Log
}

// WithObjectValue specifies the function to log an client.Object's value with.
type WithObjectValue[object any, request comparable] func(obj object) any

func (w WithObjectValue[object, request]) ApplyToHandler(o *TypedHandlerOptions[object, request]) {
	o.ObjectValue = w
}

func (w WithObjectValue[object, request]) ApplyToPredicate(o *TypedPredicateOptions[object]) {
	o.ObjectValue = w
}
