package cache

import (
	"sync"
)

type cacheStorage[Key comparable, T any] struct {
	// TODO: RWMutex
	sync.Mutex
	storage map[Key]T
}

func NewCacheStorage[Key comparable, T any]() *cacheStorage[Key, T] {
	var cache cacheStorage[Key, T]
	cache.storage = make(map[Key]T)
	return &cache
}

func (cs *cacheStorage[Key, T]) Exist(id Key) bool {
	cs.Lock()
	defer cs.Unlock()
	_, ok := cs.storage[id]
	return ok
}

func (cs *cacheStorage[Key, T]) Get(id Key) T {
	cs.Lock()
	defer cs.Unlock()
	return cs.storage[id]
}

func (cs *cacheStorage[Key, T]) Put(id Key, data T) {
	cs.Lock()
	defer cs.Unlock()
	cs.storage[id] = data
}
