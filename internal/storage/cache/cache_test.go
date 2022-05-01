package cache_test

import (
	"testing"

	"github.com/de1phin/music-transfer/internal/storage/cache"
)

func TestCache(t *testing.T) {
	cache := cache.NewCacheStorage[int64, int]()
	cache.Set(3, 10)
	if val, _ := cache.Get(3); val != 10 {
		t.Fatal("cache.Get(3): Expected 10, got", val)
	}
	if val, _ := cache.Exist(3); !val {
		t.Fatal("cache.Exist(3): Expected true, got false")
	}
	if val, _ := cache.Exist(10); val {
		t.Fatal("cache.Exist(10): Expected false, got true")
	}
}
