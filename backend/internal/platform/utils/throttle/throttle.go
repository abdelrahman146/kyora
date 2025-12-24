package throttle

import (
	"time"

	"github.com/abdelrahman146/kyora/internal/platform/cache"
)

type state struct {
	Count int   `json:"count"`
	Last  int64 `json:"last"` // unix seconds
}

// Allow implements a best-effort token bucket using cache with JSON state.
//
// It returns true when the action is allowed and false when it should be rate-limited.
// When cache isn't configured or cache operations fail, it defaults to allowing the action
// to avoid causing outages.
func Allow(c *cache.Cache, key string, window time.Duration, max int, minInterval time.Duration) bool {
	if c == nil {
		return true
	}
	if max <= 0 {
		return true
	}

	now := time.Now()
	var st state
	if data, err := c.Get(key); err == nil && len(data) > 0 {
		_ = c.Unmarshal(data, &st)
	}

	if minInterval > 0 && st.Last != 0 && now.Sub(time.Unix(st.Last, 0)) < minInterval {
		return false
	}

	st.Count++
	st.Last = now.Unix()

	ttlSeconds := int32(window.Seconds())
	if ttlSeconds <= 0 {
		ttlSeconds = 1
	}

	persist := func() {
		if b, err := c.Marshal(st); err == nil {
			_ = c.SetX(key, b, ttlSeconds)
		}
	}

	if st.Count > max {
		persist()
		return false
	}

	persist()
	return true
}
