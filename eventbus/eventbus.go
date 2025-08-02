package eventbus

import (
	"context"
	"fmt"
	"math/rand"
	"reflect"
	"slices"
	"sync"
	"time"
)

type HandlerFunc[T any] func(ctx context.Context, event T)
type handlerWrapper func(ctx context.Context, event any)

type eventPayload struct {
	ctx   context.Context
	event any
}

type Bus struct {
	subscribers      map[reflect.Type][]*subscription
	dispatchersCount int
	events           chan *eventPayload

	running bool
	ctx     context.Context
	cancel  context.CancelFunc

	wg sync.WaitGroup
	mu sync.RWMutex
}

type Subscription interface {
	ID() int
	EventType() reflect.Type
	Close()
}

type subscription struct {
	id        int
	bus       *Bus
	eventType reflect.Type
	handler   handlerWrapper
}

func (s *subscription) Close() {}

func (s *subscription) ID() int                 { return s.id }
func (s *subscription) EventType() reflect.Type { return s.eventType }

func New() *Bus {
	return &Bus{
		subscribers:      make(map[reflect.Type][]*subscription),
		dispatchersCount: 5,
	}
}

func (bus *Bus) Start(ctx context.Context) {
	bus.mu.Lock()
	defer bus.mu.Unlock()

	bus.ctx, bus.cancel = context.WithCancel(ctx)
	bus.events = make(chan *eventPayload, 100)

	for range bus.dispatchersCount {
		bus.wg.Add(1)
		go func() {
			defer bus.wg.Done()
			for {
				select {
				case <-bus.ctx.Done():
					return
				case payload, ok := <-bus.events:
					{
						if !ok {
							return
						}
						key := reflect.TypeOf(payload.event)
						bus.mu.RLock()
						subs := bus.subscribers[key]
						subs = slices.Clone(subs)
						bus.mu.RUnlock()

						for _, s := range subs {
							s.handler(payload.ctx, payload.event)
						}
					}
				}
			}
		}()
	}

	bus.running = true

}

func (bus *Bus) Stop(timeout time.Duration) error {
	done := make(chan struct{})

	go func() {
		bus.cancel()
		bus.wg.Wait()
		close(done)
	}()

	var err error
	select {
	case <-time.After(timeout):
		err = fmt.Errorf("bus shutdown timed out after %v", timeout)

	case <-done:
		// dispatchers shutdown
	}

	bus.mu.Lock()
	close(bus.events)
	bus.running = false
	bus.mu.Unlock()

	return err
}

func SubscribeWith[T any](bus *Bus, handler HandlerFunc[T]) Subscription {
	var typ T
	key := reflect.TypeOf(typ)
	bus.mu.RLock()
	subs := bus.subscribers[key]
	subs = slices.Clone(subs)
	bus.mu.RUnlock()

	id := rand.Int()
	subscription := subscription{
		id:        id,
		bus:       bus,
		eventType: key,
		handler: func(ctx context.Context, event any) {
			ev, ok := event.(T)
			if !ok {
				return
			}

			handler(ctx, ev)
		},
	}

	bus.mu.Lock()
	bus.subscribers[key] = subs
	bus.mu.Unlock()

	return &subscription
}

func PublishWith[T any](bus *Bus, ctx context.Context, event T) {
	bus.mu.RLock()
	defer bus.mu.RUnlock()
	if !bus.running {
		return
	}
	bus.events <- &eventPayload{ctx: ctx, event: event}
}

func UnsubscribeWith(bus *Bus, sub Subscription) {
	bus.mu.Lock()
	defer bus.mu.Unlock()
	subs, ok := bus.subscribers[sub.EventType()]
	if !ok {
		return
	}
	subs = slices.DeleteFunc(subs, func(i *subscription) bool { return i.id == sub.ID() })
	bus.subscribers[sub.EventType()] = subs
}
