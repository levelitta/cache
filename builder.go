package cache

type Builder[K comparable, V any] struct {
	capacity uint64
}

func NewBuilder[K comparable, V any]() *Builder[K, V] {
	return &Builder[K, V]{}
}

func (b *Builder[K, V]) Capacity(cap uint64) {
	b.capacity = cap
}

func (b *Builder[K, V]) Build() *Cache[K, V] {
	return &Cache[K, V]{
		capacity: b.capacity,
	}
}
