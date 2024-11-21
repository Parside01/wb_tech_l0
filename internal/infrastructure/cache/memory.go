package cache

import (
	"container/list"
	"sync"
)

type lruCacheItem struct {
	cacheable Cacheable
	element   *list.Element
}

type lruMemoryCache struct {
	mutex    sync.RWMutex
	memory   map[string]*lruCacheItem
	capacity int
	list     *list.List
}

func NewMemoryCache(capacity int) Cache {
	return &lruMemoryCache{
		memory:   make(map[string]*lruCacheItem, capacity),
		mutex:    sync.RWMutex{},
		list:     list.New(),
		capacity: capacity,
	}
}

func (c *lruMemoryCache) Set(key string, value Cacheable) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if e, ok := c.memory[key]; ok {
		c.list.MoveToFront(e.element)
		e.cacheable = value
		return
	}

	if c.list.Len() == c.capacity {
		back := c.list.Back()
		if back != nil {
			c.list.Remove(back)
			delete(c.memory, back.Value.(*lruCacheItem).cacheable.Key())
		}
	}

	newItem := &lruCacheItem{cacheable: value}
	newElem := c.list.PushFront(newItem)
	newItem.element = newElem

	c.memory[key] = newItem
}

func (c *lruMemoryCache) Get(key string) (Cacheable, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	item, ok := c.memory[key]
	if !ok {
		return nil, false
	}

	c.list.MoveToFront(item.element)
	return item.cacheable, true
}

func (c *lruMemoryCache) Delete(key string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if item, ok := c.memory[key]; ok {
		c.list.Remove(item.element)
		delete(c.memory, key)
	}
}
