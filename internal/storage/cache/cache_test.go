package cache_test

import (
	"testing"

	"github.com/de1phin/music-transfer/internal/storage/cache"
)

func TestCache(t *testing.T) {
	cache := cache.NewCacheStorage[int]()
	cache.Put(3, 10)
	if cache.Get(3) != 10 {
		t.Fatal("cache.Get(3): Expected 10, got", cache.Get(3))
	}
	if !cache.Exist(3) {
		t.Fatal("cache.Exist(3): Expected true, got false")
	}
	if cache.Exist(10) {
		t.Fatal("cache.Exist(10): Expected false, got true")
	}
}
