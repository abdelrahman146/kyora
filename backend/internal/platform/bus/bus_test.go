package bus_test

import (
	"sync"
	"testing"
	"time"

	"github.com/abdelrahman146/kyora/internal/platform/bus"
	"github.com/stretchr/testify/require"
)

func TestBus_DeliversAllEvents(t *testing.T) {
	t.Parallel()

	b := bus.New()
	defer b.Close()

	const n = 250

	var (
		mu       sync.Mutex
		received = make([]int, 0, n)
		wg       sync.WaitGroup
	)

	wg.Add(n)
	unsub := b.Listen("test.topic", func(v any) {
		mu.Lock()
		received = append(received, v.(int))
		mu.Unlock()
		wg.Done()
	})
	defer unsub()

	for i := 0; i < n; i++ {
		b.Emit("test.topic", i)
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		// ok
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for events")
	}

	mu.Lock()
	defer mu.Unlock()
	require.Len(t, received, n)
}
