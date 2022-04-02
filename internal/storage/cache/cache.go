package cache

import (
	"log"
	"sync"
)

type cacheStorage[T any] struct {
	// TODO: RWMutex
	sync.Mutex
	storage map[int64]T
}

func NewCacheStorage[T any]() *cacheStorage[T] {
	var cache cacheStorage[T]
	cache.storage = make(map[int64]T)
	return &cache
}

func (cs *cacheStorage[T]) Exist(id int64) bool {
	cs.Lock()
	defer cs.Unlock()
	_, ok := cs.storage[id]
	return ok
}

func (cs *cacheStorage[T]) Get(id int64) T {
	cs.Lock()
	defer cs.Unlock()
	return cs.storage[id]
}

func (cs *cacheStorage[T]) Put(id int64, data T) {
	cs.Lock()
	defer cs.Unlock()
	cs.storage[id] = data
	log.Println(data)
}
