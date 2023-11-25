package cache

import (
	"container/list"
	"sync"
	"time"
)

func zeroValue[T any]() T {
	var zero T
	return zero
}

type Cache[K comparable, V any] struct {
	items       map[K]*Item[V]
	mu          sync.Mutex
	evictPolicy *list.List
	capacity    uint64
	len         uint64
	now         func() int64
}

type Item[V any] struct {
	value      V
	expired    int64
	evictQueue *list.Element
}

func NewCache[K comparable, V any](capacity uint64) *Cache[K, V] {
	return &Cache[K, V]{
		items:       make(map[K]*Item[V], capacity),
		evictPolicy: list.New(),
		now: func() int64 {
			return time.Now().Unix()
		},
	}
}

func (c *Cache[K, V]) Set(key K, value V) {
	c.SetWithTTL(key, value, 0)
}

func (c *Cache[K, V]) SetWithTTL(key K, value V, ttl time.Duration) {
	var expired int64
	if ttl > 0 {
		expired = c.now() + int64(ttl.Seconds())
	}

	c.mu.Lock()

	if item, has := c.items[key]; has {
		item.value = value
		item.expired = expired
		c.evictPolicy.MoveToFront(item.evictQueue)
	} else {
		item = &Item[V]{
			value:      value,
			expired:    expired,
			evictQueue: c.evictPolicy.PushFront(item),
		}
		c.items[key] = item
	}

	c.mu.Unlock()
}

func (c *Cache[K, V]) Get(key K) (V, bool) {
	c.mu.Lock()
	item, has := c.items[key]
	c.mu.Unlock()

	if !has ||
		item.expired > 0 && c.now() >= item.expired {
		return zeroValue[V](), false
	}

	return item.value, has
}

func (c *Cache[K, V]) Has(key K) bool {
	c.mu.Lock()
	item, has := c.items[key]
	c.mu.Unlock()

	return has && (item.expired == 0 || c.now() < item.expired)
}

func (c *Cache[K, V]) Del(key K) {
	c.mu.Lock()
	item, has := c.items[key]
	c.mu.Unlock()

	if !has {
		return
	}

	c.evictPolicy.Remove(item.evictQueue)
	delete(c.items, key)
}

func (c *Cache[K, V]) evict() {

}
