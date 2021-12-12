package limit_util

import (
	"sync/atomic"
	"time"
)

type CountLimiter struct {
	counter      int64
	limit        int64
	intervalNano int64
	unixNano     int64
}

func NewCountLimiter(interval time.Duration, limit int64) *CountLimiter {
	return &CountLimiter{
		counter:      0,
		limit:        limit,
		intervalNano: int64(interval),
		unixNano:     time.Now().UnixNano(),
	}
}

func (c *CountLimiter) Allow() bool {
	now := time.Now().UnixNano()
	if now-c.unixNano > c.intervalNano {
		atomic.StoreInt64(&c.counter, 0)
		atomic.StoreInt64(&c.unixNano, now)
		return true
	}
	atomic.AddInt64(&c.counter, 1)
	return c.counter <= c.limit
}
