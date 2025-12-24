package bus

import (
	"log/slog"
	"runtime/debug"
	"sync"
	"sync/atomic"
)

type Event struct {
	Topic   Topic
	Payload any
}

type Bus struct {
	mu     sync.RWMutex
	topics map[Topic]map[uint64]chan any
	emitCh chan Event
	stop   chan struct{}
	wg     sync.WaitGroup
	nextID atomic.Uint64
	subBuf int
	closed atomic.Bool
}

const (
	defaultEmitBuffer = 1024
	defaultSubBuffer  = 128
)

func New() *Bus {
	b := &Bus{
		topics: make(map[Topic]map[uint64]chan any),
		emitCh: make(chan Event, defaultEmitBuffer),
		stop:   make(chan struct{}),
		subBuf: defaultSubBuffer,
	}
	b.wg.Add(1)
	go b.dispatch()
	return b
}

// Listen subscribes to a topic and handles payloads asynchronously.
// It returns an unsubscribe function to stop receiving events.
func (b *Bus) Listen(topic Topic, handler func(any)) (unsubscribe func()) {
	if handler == nil {
		return func() {}
	}

	ch := make(chan any, b.subBuf)
	id := b.nextID.Add(1)

	b.mu.Lock()
	if b.topics[topic] == nil {
		b.topics[topic] = make(map[uint64]chan any)
	}
	b.topics[topic][id] = ch
	b.mu.Unlock()

	b.wg.Add(1)
	go func() {
		defer b.wg.Done()
		for payload := range ch {
			func() {
				defer func() {
					if r := recover(); r != nil {
						slog.Error("bus handler panicked", "topic", topic, "panic", r, "stack", string(debug.Stack()))
					}
				}()
				handler(payload)
			}()
		}
	}()

	return func() {
		b.mu.Lock()
		if subs, ok := b.topics[topic]; ok {
			if subCh, ok := subs[id]; ok {
				delete(subs, id)
				close(subCh)
			}
			if len(subs) == 0 {
				delete(b.topics, topic)
			}
		}
		b.mu.Unlock()
	}
}

// Emit publishes a payload to a topic asynchronously.
func (b *Bus) Emit(topic Topic, payload any) {
	if b.closed.Load() {
		return
	}
	ev := Event{Topic: topic, Payload: payload}

	// Reliable delivery: apply backpressure instead of dropping events.
	select {
	case b.emitCh <- ev:
	case <-b.stop:
	}
}

// Close gracefully stops the bus and all subscriber goroutines.
func (b *Bus) Close() {
	if !b.closed.CompareAndSwap(false, true) {
		return
	}
	close(b.stop)

	// Close all subscriber channels
	b.mu.Lock()
	for _, subs := range b.topics {
		for id, ch := range subs {
			delete(subs, id)
			close(ch)
		}
	}
	b.topics = make(map[Topic]map[uint64]chan any)
	b.mu.Unlock()

	b.wg.Wait()
}

func (b *Bus) dispatch() {
	defer b.wg.Done()
	for {
		select {
		case <-b.stop:
			return
		case ev := <-b.emitCh:
			// Snapshot subscriber channels to avoid holding lock during sends
			b.mu.RLock()
			var targets []chan any
			if subs, ok := b.topics[ev.Topic]; ok {
				targets = make([]chan any, 0, len(subs))
				for _, ch := range subs {
					targets = append(targets, ch)
				}
			}
			b.mu.RUnlock()

			for _, ch := range targets {
				func(ch chan any) {
					defer func() { _ = recover() }() // ignore sends to channels closed by unsubscribe races
					select {
					case ch <- ev.Payload:
					case <-b.stop:
					}
				}(ch)
			}
		}
	}
}
