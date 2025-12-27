package throttle

import (
	"math"
	"time"

	"github.com/abdelrahman146/kyora/internal/platform/cache"
)

type cooldownState struct {
	Start int64 `json:"start"` // unix seconds
}

// Cooldown implements a per-key cooldown window.
//
// If the key does not exist, it is created with a TTL equal to cooldown and the action is allowed.
// If the key exists, the action is denied and the remaining time until the cooldown ends is returned.
//
// This uses an atomic cache Add to make concurrent calls safe.
// When cache isn't configured or cache operations fail, it defaults to allowing the action
// to avoid causing outages.
func Cooldown(c *cache.Cache, key string, cooldown time.Duration) (allowed bool, retryAfter time.Duration) {
	if c == nil {
		return true, cooldown
	}
	if cooldown <= 0 {
		return true, 0
	}

	now := time.Now().UTC()
	st := cooldownState{Start: now.Unix()}

	b, err := c.Marshal(st)
	if err != nil {
		return true, cooldown
	}

	ttlSeconds := int32(math.Ceil(cooldown.Seconds()))
	if ttlSeconds <= 0 {
		ttlSeconds = 1
	}

	if err := c.AddX(key, b, ttlSeconds); err == nil {
		return true, cooldown
	} else if !cache.IsNotStored(err) {
		// best-effort: allow if cache is flaky
		return true, cooldown
	}

	// Not stored: key already exists
	data, err := c.Get(key)
	if err != nil || len(data) == 0 {
		return false, cooldown
	}

	var existing cooldownState
	if err := c.Unmarshal(data, &existing); err != nil || existing.Start == 0 {
		return false, cooldown
	}

	elapsed := now.Sub(time.Unix(existing.Start, 0))
	remaining := cooldown - elapsed
	if remaining < 0 {
		remaining = 0
	}
	return false, remaining
}
