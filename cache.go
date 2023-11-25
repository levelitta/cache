package cache

import (
	"container/list"
	"sync"
)

func zeroValue[T any]() T {
	var zero T
	return zero
}

type Cache[K comparable, V any] struct {
	items      map[K]*Item[V]
	mu         sync.Mutex // TODO:заменить на RWMutes и проверить производительность
	evictQueue *list.List
	capacity   uint64
}

type Item[V any] struct {
	value      V
	evictQueue *list.Element
}

func NewCache[K comparable, V any](capacity uint64) *Cache[K, V] {
	return &Cache[K, V]{
		items:      make(map[K]*Item[V], capacity),
		evictQueue: list.New(),
	}
}

func (c *Cache[K, V]) Set(key K, value V) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if item, has := c.items[key]; has {
		item.value = value
		c.evictQueue.MoveToFront(item.evictQueue)
	} else {
		item = &Item[V]{value: value}
		item.evictQueue = c.evictQueue.PushFront(item)
		c.items[key] = item
	}
}

func (c *Cache[K, V]) Get(key K) (V, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, ok := c.items[key]
	if !ok {
		return zeroValue[V](), false
	}
	return item.value, ok
}

func (c *Cache[K, V]) Has(key K) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, has := c.items[key]
	return has
}

func (c *Cache[K, V]) Del(key K) {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, ok := c.items[key]
	if !ok {
		return
	}

	c.evictQueue.Remove(item.evictQueue) // TODO: попробовать при удалении элемента класть его в sync.Pool
	delete(c.items, key)
}

func (c *Cache[K, V]) evict() {

}
