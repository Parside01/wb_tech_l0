package cache

import (
	"sync"
	"wb_tech_l0/internal/entity"
)

type memoryCache struct {
	mutex  sync.RWMutex
	memory map[string]*entity.Order
}

func NewMemoryCache() Cache {
	return &memoryCache{
		memory: make(map[string]*entity.Order),
		mutex:  sync.RWMutex{},
	}
}

func (c *memoryCache) Set(key string, order *entity.Order) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.memory[key] = order
}

func (c *memoryCache) Get(key string) (*entity.Order, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	order, ok := c.memory[key]
	return order, ok
}

func (c *memoryCache) Delete(key string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	delete(c.memory, key)
}
