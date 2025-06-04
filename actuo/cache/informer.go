// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0

package cache

import (
	"context"
	"fmt"
	"sync"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/util/sets"
	"spheric.cloud/spheric/utils/chans"
	"spheric.cloud/spheric/utils/container/squeue"
)

func NewInformer[Key comparable, Object any](
	keyFunc KeyFunc[Key, Object],
	lw ListerWatcher[Object],
	handler ResourceEventHandler[Object],
) (Store[Key, Object], Synced, Controller) {
	store := NewCache[Key, Object](keyFunc)
	reconciler := NewDeltaFIFOReconciler(keyFunc, store, handler)
	ctrl := NewController(keyFunc, lw, store, reconciler, ControllerOptions{})
	return store, reconciler, ctrl
}

type notification interface {
	notification()
}

type resyncedNotification struct{}

func (resyncedNotification) notification() {}

type populatedNotification struct{}

func (populatedNotification) notification() {}

type addNotification[Object any] struct {
	newObj        Object
	isInitialList bool
}

func (addNotification[Object]) notification() {}

type updateNotification[Object any] struct {
	oldObj Object
	newObj Object
}

func (updateNotification[Object]) notification() {}

type deleteNotification[Object any] struct {
	oldObj Object
}

func (deleteNotification[Object]) notification() {}

type listener[Object any] struct {
	notifications chan notification
	handler       ResourceEventHandler[Object]
	inResync      <-chan struct{}
	reqResync     chan<- *listener[Object]
	synced        chan struct{}
}

func (l *listener[Object]) add(notification notification) {
	l.notifications <- notification
}

func (l *listener[Object]) run() {
	notifications := make(chan notification)
	stopPump := make(chan struct{})
	defer close(stopPump)

	go func() {
		defer close(notifications)
		chans.BufferedPumpExitOnClose(notifications, l.notifications, squeue.New[notification](1024))
	}()

	var (
		nReqResync = chans.ToggleOf[*listener[Object]](l.reqResync).Disable()
		nInResync  = chans.ToggleOfFunc[struct{}](func() <-chan struct{} { return l.inResync })
	)

	for {
		select {
		case _, ok := <-nInResync.C():
			// Since we could read a resync request,
			// Temporarily disable reading further resync requests until
			// we can dispatch a resync request to the central resyncer.
			nInResync = nil

			if !ok {
				// If the incoming resync request channel is closed,
				// set it to nil to eliminate this case statement forever.
				l.inResync = nil
				continue
			}

			nReqResync.Enable()
		case nReqResync.C() <- l:
			// We could dispatch a resync request - disable until we get a resynced event.
			nReqResync.Disable()
		case n, ok := <-notifications:
			if !ok {
				return
			}

			switch n := n.(type) {
			case populatedNotification:
				close(l.synced)
			case resyncedNotification:
				// Listen for incoming resync requests again
				nInResync.Enable()
			case addNotification[Object]:
				l.handler.OnAdd(n.newObj, n.isInitialList)
			case updateNotification[Object]:
				l.handler.OnUpdate(n.oldObj, n.newObj)
			case deleteNotification[Object]:
				l.handler.OnDelete(n.oldObj)
			}
		}
	}
}

func (l *listener[Object]) Synced() <-chan struct{} {
	return l.synced
}

type sharedInformer[Key comparable, Object any] struct {
	startedMu sync.Mutex
	started   bool

	wg sync.WaitGroup

	store      Store[Key, Object]
	keyFunc    KeyFunc[Key, Object]
	controller Controller

	// resyncReqs are resync requests by the listeners
	resyncReqs chan *listener[Object]
	// resync is the channel to the controller to trigger resyncs
	resync chan struct{}
	// resyncCompleted is a channel to signal if a resync completed
	resyncCompleted chan struct{}

	blockDeltas sync.Mutex

	listenersMu sync.RWMutex
	listeners   map[*listener[Object]]bool

	synced chan struct{}

	logger logr.Logger
}

func (s *sharedInformer[Key, Object]) resyncer(ctx context.Context) {
	var (
		resyncReqs = chans.ToggleOf[*listener[Object]](s.resyncReqs)

		// Initialize this as disabled since we don't want to send
		// resyncs until we got any request
		resync = chans.ToggleOf[struct{}](s.resync).Disable()
	)

	for {
		select {
		case <-ctx.Done():
			return
		case req := <-resyncReqs.C():
			func() {
				s.listenersMu.Lock()
				defer s.listenersMu.Unlock()

				if _, ok := s.listeners[req]; !ok {
					// Got a resync request from a listener that was removed.
					return
				}

				s.blockDeltas.Lock()
				defer s.blockDeltas.Unlock()

				// TODO: Do we want to set a timeout here in case of a busy resync sender?
				reqs := sets.New[*listener[Object]]()
			GatherFurtherReqs:
				for {
					select {
					case req := <-resyncReqs.C():
						reqs.Insert(req)
					default:
						break GatherFurtherReqs
					}
				}

				for req := range reqs {
					// TODO: Set to syncing
					s.listeners[req] = true
				}

				// Do not accept resync requests until we could dispatch a resync.
				resyncReqs.Disable()
				resync.Enable()
			}()
		case resync.C() <- struct{}{}:
			// Do not accept sending resyncs until we have resynced
			resync.Disable()
		case <-s.resyncCompleted:
			// Accept resync requests again
			resyncReqs.Enable()
		}
	}
}

func (s *sharedInformer[Key, Object]) Reconcile(ctx context.Context, req ReconcileRequest[Object]) error {
	s.blockDeltas.Lock()
	defer s.blockDeltas.Unlock()

	switch req.Type {
	case Populated:
		close(s.synced)
		s.distribute(populatedNotification{}, false)
		return nil
	case Resynced:
		select {
		case <-ctx.Done():
			return ctx.Err()
		case s.resyncCompleted <- struct{}{}:
			s.distribute(resyncedNotification{}, true)
			return nil
		}
	case NewDeltas:
		return s.processDeltas(req.Deltas)
	default:
		return fmt.Errorf("unknown event type: %q", req.Type)
	}
}

func (s *sharedInformer[Key, Object]) processDeltas(deltas Deltas[Object]) error {
	for _, d := range deltas {
		obj := d.Object
		key, err := s.keyFunc(obj)
		if err != nil {
			return err
		}

		isSync := d.Type == Resync
		switch d.Type {
		case Resync, Populate, Added, Updated:
			if old, err := s.store.Get(key); err == nil {
				if err := s.store.Set(obj); err != nil {
					return err
				}
				s.distribute(updateNotification[Object]{
					oldObj: old,
					newObj: obj,
				}, isSync)
			} else {
				if err := s.store.Set(obj); err != nil {
					return err
				}
				s.distribute(addNotification[Object]{
					newObj:        obj,
					isInitialList: d.Type == Populate,
				}, isSync)
			}
		case Deleted:
			if err := s.store.Delete(key); err != nil {
				return err
			}
			s.distribute(deleteNotification[Object]{
				oldObj: obj,
			}, isSync)
		}
	}
	return nil
}

func (s *sharedInformer[Key, Object]) addListener(listener *listener[Object]) {
	s.listenersMu.Lock()
	defer s.listenersMu.Unlock()

	s.listeners[listener] = true
}

func (s *sharedInformer[Key, Object]) AddEventHandler(handler ResourceEventHandler[Object], opts EventHandlerOptions) (ResourceEventHandlerRegistration, error) {
	s.startedMu.Lock()
	defer s.startedMu.Unlock()

	l := &listener[Object]{
		notifications: make(chan notification),
		handler:       handler,
		synced:        make(chan struct{}),
		inResync:      opts.Resync,
	}

	if !s.started {
		s.addListener(l)
		return l, nil
	}

	s.blockDeltas.Lock()
	defer s.blockDeltas.Unlock()

	s.addListener(l)

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		l.run()
	}()

	// Add synthetic add notifications and (iff already synced) a populated notification
	for obj := range s.store.All() {
		l.add(addNotification[Object]{newObj: obj, isInitialList: true})
	}
	select {
	case <-s.synced:
		l.add(populatedNotification{})
	default:
	}

	return l, nil
}

func (s *sharedInformer[Key, Object]) RemoveEventHandler(handle ResourceEventHandlerRegistration) error {
	l, ok := handle.(*listener[Object])
	if !ok {
		return fmt.Errorf("unrecognized registration type %T", handle)
	}

	s.blockDeltas.Lock()
	defer s.blockDeltas.Unlock()

	s.listenersMu.Lock()
	defer s.listenersMu.Unlock()

	if _, ok := s.listeners[l]; !ok {
		return nil
	}
	delete(s.listeners, l)

	close(l.notifications)

	return nil
}

func (s *sharedInformer[Key, Object]) Store() Store[Key, Object] {
	return s.store
}

func (s *sharedInformer[Key, Object]) Synced() <-chan struct{} {
	return s.synced
}

func (s *sharedInformer[Key, Object]) distribute(n notification, sync bool) {
	s.listenersMu.RLock()
	defer s.listenersMu.RUnlock()

	for l, isSyning := range s.listeners {
		switch {
		case !sync:
			// non-sync messages are delivered to every listener
			l.add(n)
		case isSyning:
			// sync messages are delivered to every syncing listener
			l.add(n)
		default:
			// skipping a sync message for a non-syncing listener
		}
	}
}

func (s *sharedInformer[Key, Object]) Run(ctx context.Context) error {
	err := func() error {
		s.startedMu.Lock()
		defer s.startedMu.Unlock()

		if s.started {
			return fmt.Errorf("shared informer already started")
		}

		s.started = true

		s.wg.Add(1)
		go func() {
			defer s.wg.Done()

			if err := s.controller.Run(ctx); err != nil {
				s.logger.Error(err, "Error running controller")
			}
		}()

		s.listenersMu.RLock()
		defer s.listenersMu.RUnlock()

		for l := range s.listeners {
			s.wg.Add(1)
			go func() {
				defer s.wg.Done()

				l.run()
			}()
		}

		return nil
	}()
	if err != nil {
		return err
	}

	<-ctx.Done()
	s.logger.Info("Context canceled, stopping listeners")
	func() {
		s.listenersMu.Lock()
		defer s.listenersMu.Unlock()

		for l := range s.listeners {
			delete(s.listeners, l)
			close(l.notifications)
		}
	}()

	s.logger.Info("Waiting for goroutines to finish")
	s.wg.Wait()

	s.logger.Info("Queue shutdown complete")
	return nil
}

type EventHandlerOptions struct {
	Resync <-chan struct{}
}

type SharedInformer[Key comparable, Object any] interface {
	AddEventHandler(handler ResourceEventHandler[Object], opts EventHandlerOptions) (ResourceEventHandlerRegistration, error)
	RemoveEventHandler(handle ResourceEventHandlerRegistration) error
	Store() Store[Key, Object]
	Run(ctx context.Context) error
	Synced
}

type ResourceEventHandlerRegistration interface {
	Synced
}

type SharedInformerOptions struct {
	Logger logr.Logger
}

func NewSharedInformer[Key comparable, Object any](
	keyFunc KeyFunc[Key, Object],
	lw ListerWatcher[Object],
	opts SharedInformerOptions,
) SharedInformer[Key, Object] {
	logger := opts.Logger
	if logger.GetSink() == nil {
		opts.Logger = logr.Discard()
	}

	si := &sharedInformer[Key, Object]{
		store:           NewCache[Key, Object](keyFunc),
		keyFunc:         keyFunc,
		resync:          make(chan struct{}),
		resyncCompleted: make(chan struct{}),
		listeners:       make(map[*listener[Object]]bool),
		synced:          make(chan struct{}),
		logger:          logger,
	}
	si.controller = NewController(keyFunc, lw, si.store, si, ControllerOptions{
		Resync: si.resync,
		Logger: logger.WithName("controller"),
	})
	return si
}
