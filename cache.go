package cache

import (
	"container/list"
	"github.com/levelitta/cache/evict"
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
	evictPolicy evict.Policy[K]
	capacity    uint64
	now         func() int64 // mockable field
}

type Item[V any] struct {
	value             V
	expired           int64
	evictQueueElement *list.Element
}

func NewCache[K comparable, V any](capacity uint64) *Cache[K, V] {
	return &Cache[K, V]{
		items:       make(map[K]*Item[V], capacity),
		capacity:    capacity,
		evictPolicy: evict.NewLruPolicy[K](),
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

	// TODO проверить вариант, когда лок вешается только на получение из мапы и отдельно на пуш нового элемента
	c.mu.Lock()
	if item, has := c.items[key]; has {
		item.value = value
		item.expired = expired
		c.evictPolicy.IncScore(item.evictQueueElement)
	} else {
		item = &Item[V]{
			value:             value,
			expired:           expired,
			evictQueueElement: c.evictPolicy.Push(key),
		}
		c.items[key] = item
	}
	if c.capacity != 0 && uint64(len(c.items)) > c.capacity {
		c.evict()
	}

	c.mu.Unlock()
}

func (c *Cache[K, V]) Get(key K) (V, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, has := c.items[key]
	if !has {
		return zeroValue[V](), false
	}
	if item.expired > 0 && c.now() >= item.expired {
		c.delItem(item)
		return zeroValue[V](), false
	}

	c.evictPolicy.IncScore(item.evictQueueElement)

	return item.value, has
}

func (c *Cache[K, V]) Has(key K) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, has := c.items[key]
	if !has {
		return false
	}
	if item.expired > 0 && c.now() >= item.expired {
		c.delItem(item)
		return false
	}

	c.evictPolicy.IncScore(item.evictQueueElement)

	return true
}

func (c *Cache[K, V]) Del(key K) {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, has := c.items[key]
	if !has {
		return
	}

	c.delItem(item)
}

func (c *Cache[K, V]) delItem(item *Item[V]) {
	key := c.evictPolicy.Del(item.evictQueueElement)
	delete(c.items, key)
}

func (c *Cache[K, V]) evict() {
	key, ok := c.evictPolicy.Evict()
	if ok {
		delete(c.items, key)
	}
}
