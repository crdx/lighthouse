package cache

import (
	"sync"
	"time"
)

type TemporalCache[T comparable] struct {
	cache map[T]time.Time
	mutex sync.Mutex
}

func NewTemporal[T comparable]() *TemporalCache[T] {
	return &TemporalCache[T]{
		cache: make(map[T]time.Time),
	}
}

// SeenWithinLast checks if a key been seen within the last duration.
func (self *TemporalCache[T]) SeenWithinLast(key T, duration time.Duration) bool {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	if lastSeen, found := self.cache[key]; found {
		if time.Since(lastSeen) < duration {
			return true
		}
	}

	self.cache[key] = time.Now()
	return false
}
