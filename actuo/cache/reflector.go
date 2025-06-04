// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0

package cache

import (
	"context"
	"fmt"
	"slices"

	"github.com/go-logr/logr"
	"spheric.cloud/spheric/actuo/watch"
)

type Reflector[Object any] struct {
	lw     ListerWatcher[Object]
	store  ReflectorSink[Object]
	logger logr.Logger
}

type ReflectorOptions struct {
	Logger logr.Logger
}

func NewReflector[Object any](lw ListerWatcher[Object], store ReflectorSink[Object], opts ReflectorOptions) *Reflector[Object] {
	logger := opts.Logger
	if logger.GetSink() == nil {
		logger = logr.Discard()
	}

	return &Reflector[Object]{
		lw:     lw,
		store:  store,
		logger: logger,
	}
}

func (r *Reflector[Object]) Run(ctx context.Context) error {
	objs, err := r.lw.List(ctx)
	if err != nil {
		return err
	}

	if err := r.store.Populate(ctx, slices.Values(objs)); err != nil {
		return err
	}

	w, err := r.lw.Watch(ctx)
	if err != nil {
		return err
	}
	defer w.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case evt, ok := <-w.Events():
			if !ok {
				return fmt.Errorf("watch closed")
			}

			switch evt.Type {
			case watch.EventTypeCreated:
				if err := r.store.Add(ctx, evt.Object); err != nil {
					r.logger.Error(err, "Error adding object")
				}
			case watch.EventTypeUpdated:
				if err := r.store.Update(ctx, evt.Object); err != nil {
					r.logger.Error(err, "Error updating object")
				}
			case watch.EventTypeDeleted:
				if err := r.store.Delete(ctx, evt.Object); err != nil {
					r.logger.Error(err, "Error deleting object")
				}
			}
		}
	}
}
