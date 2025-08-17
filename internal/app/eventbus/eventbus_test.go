package eventbus

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"os"
	"strings"
	"sync"
	"testing"
	"time"
)

type event1 struct {
	V      int
	Called int
}
type event2 struct {
	V      int
	Called int
}

type syncBuffer struct {
	buffer bytes.Buffer
	sync.Mutex
}

var w io.Writer = &syncBuffer{}

func newSyncBuffer() *syncBuffer {
	return &syncBuffer{
		buffer: *bytes.NewBuffer(nil),
	}
}

func (s *syncBuffer) Write(p []byte) (n int, err error) {
	s.Lock()
	defer s.Unlock()
	return s.buffer.Write(p)
}

func (s *syncBuffer) String() string {
	s.Lock()
	defer s.Unlock()
	return s.buffer.String()
}

func TestSubscribe(t *testing.T) {

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))

	t.Run("subscribe and unsubscribe", func(t *testing.T) {
		bus := New(logger)
		subs := []Subscription{}
		for range 5 {
			sub := SubscribeWith(bus, func(ctx context.Context, ev event1) {
				ev.Called += 1
			})

			subs = append(subs, sub)
		}

		total := len(subs)

		if len(bus.subscribers[subs[0].EventType()]) != total {
			t.Fatalf("expected subscribers to be %d after subscribe", total)
		}

		for _, sub := range subs {
			UnsubscribeWith(bus, sub)
			total -= 1

			current := len(bus.subscribers[sub.EventType()])
			if current != total {
				t.Fatalf("expected subscribers to be %d after unsubscribe but was %d", total, current)
			}
		}
	})
	t.Run("publish 1 event", func(t *testing.T) {
		bus := New(logger)

		called := make(chan struct{})
		sub := SubscribeWith(bus, func(ctx context.Context, ev *event1) {
			ev.Called += 1
			close(called)
		})

		bus.Start(t.Context())

		event := &event1{V: 10}
		PublishWith(bus, t.Context(), event)
		<-called
		if event.Called != 1 {
			t.Fatalf("expected event called to be 1 instead it was %d", event.Called)
		}

		UnsubscribeWith(bus, sub)
		defer bus.Stop(1 * time.Second)
	})
	t.Run("publish event with no subscribers", func(t *testing.T) {
		bus := New(logger)

		bus.Start(t.Context())

		event := &event1{V: 10}
		PublishWith(bus, t.Context(), event)
		<-time.After(50 * time.Millisecond)
		defer bus.Stop(1 * time.Second)
	})
	t.Run("publish 2 different events", func(t *testing.T) {
		bus := New(logger)

		wg := sync.WaitGroup{}
		subs := []Subscription{}
		subs = append(subs, SubscribeWith(bus, func(ctx context.Context, ev *event1) {
			ev.Called += 1
			wg.Done()
		}))
		subs = append(subs, SubscribeWith(bus, func(ctx context.Context, ev *event2) {
			ev.Called += 1
			wg.Done()
		}))

		bus.Start(t.Context())

		wg.Add(2 * 10) // 2 publish * 10 events
		for i := range 10 {
			PublishWith(bus, t.Context(), &event1{
				V: i,
			})
			PublishWith(bus, t.Context(), &event2{
				V: i + 10,
			})
		}

		wg.Wait()

		for _, s := range subs {
			UnsubscribeWith(bus, s)
		}
		defer bus.Stop(1 * time.Second)
	})
	t.Run("panic in publish results in error", func(t *testing.T) {
		buf := newSyncBuffer()

		logger := slog.New(slog.NewTextHandler(buf, nil))
		bus := New(logger)

		wg := sync.WaitGroup{}
		sub := SubscribeWith(bus, func(ctx context.Context, ev *event1) {
			defer wg.Done()
			panic("I am an error")
		})

		bus.Start(t.Context())

		event := &event1{V: 10}
		wg.Add(1)
		PublishWith(bus, t.Context(), event)
		wg.Wait()
		// have to have this small sleep here otherwise things complete too quickly and hte defer doesnt' fire
		time.Sleep(50 * time.Millisecond)

		logContent := buf.String()
		if !strings.Contains(logContent, "I am an error") {
			t.Fatalf("expected string 'I am an error' in log content but got: \n%s", logContent)
		}

		UnsubscribeWith(bus, sub)
		defer bus.Stop(1 * time.Second)
	})

}
