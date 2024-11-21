package cache

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"testing"
)

type TestCacheData struct {
	ID    string
	Value string
}

func (t *TestCacheData) Key() string {
	return t.ID
}

func TestLruMemoryCache(t *testing.T) {
	TestSetGetFromLruMemoryCache(t)
	TestDeleteFromLruMemoryCache(t)
	TestExceedingCapacityLruMemoryCache(t)
}

func TestSetGetFromLruMemoryCache(t *testing.T) {
	cacheCap := 20
	cache := NewMemoryCache(cacheCap)

	items := make([]*TestCacheData, cacheCap)
	for i := range items {
		items[i] = &TestCacheData{
			ID:    gofakeit.UUID(),
			Value: gofakeit.Sentence(5),
		}
		cache.Set(items[i].Key(), items[i])
	}

	for _, item := range items {
		if cacheItem, ok := cache.Get(item.Key()); ok {
			assert.Equal(t, item, cacheItem)
		} else {
			t.Errorf("Expected to find item with key %s, but got none", item.Key())
		}
	}
}

func TestExceedingCapacityLruMemoryCache(t *testing.T) {
	cacheCap := 20
	cache := NewMemoryCache(cacheCap)
	items := make([]*TestCacheData, cacheCap+1)
	for i := range items {
		items[i] = &TestCacheData{
			ID:    gofakeit.UUID(),
			Value: gofakeit.Sentence(10),
		}
		cache.Set(items[i].Key(), items[i])
	}

	_, found := cache.Get(items[0].Key())
	assert.False(t, found)
}

func TestDeleteFromLruMemoryCache(t *testing.T) {
	cacheCap := 20
	cache := NewMemoryCache(cacheCap)
	items := make([]*TestCacheData, cacheCap)
	for i := range items {
		items[i] = &TestCacheData{
			ID:    gofakeit.UUID(),
			Value: gofakeit.Sentence(10),
		}
		cache.Set(items[i].Key(), items[i])
	}
	for i := range items {
		cache.Delete(items[i].Key())
		_, found := cache.Get(items[i].Key())
		assert.False(t, found)
	}
}
