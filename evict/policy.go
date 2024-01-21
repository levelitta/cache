package evict

import "container/list"

type Policy[K comparable] interface {
	IncScore(item *list.Element)
	DecScore(item *list.Element)
	Push(key K) *list.Element
	Del(item *list.Element) K
	Evict() (K, bool)
	Len() int
}

type LruPolicy[K comparable] struct {
	list *list.List
}

func NewLruPolicy[K comparable]() *LruPolicy[K] {
	return &LruPolicy[K]{
		list: list.New(),
	}
}

// IncScore TODO :check send item via value
func (p *LruPolicy[K]) IncScore(item *list.Element) {
	p.list.MoveToFront(item)
}

func (p *LruPolicy[K]) DecScore(item *list.Element) {
	p.list.MoveToBack(item)
}

func (p *LruPolicy[K]) Push(key K) *list.Element {
	return p.list.PushFront(key)
}

func (p *LruPolicy[K]) Del(item *list.Element) K {
	p.list.Remove(item)

	return item.Value.(K)
}

func (p *LruPolicy[K]) Evict() (K, bool) {
	item := p.list.Back()
	if item == nil {
		return zeroValue[K](), false
	}

	return p.Del(item), true
}

func (p *LruPolicy[K]) Len() int {
	return p.list.Len()
}

func zeroValue[T any]() T {
	var zero T
	return zero
}
