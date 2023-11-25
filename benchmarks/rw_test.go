package benchmarks

import (
	"fmt"
	"github.com/levelitta/cache"
	"testing"
)

var capacity uint64 = 10_000_000

func BenchmarkCache(b *testing.B) {
	b.Run("empty cache: reads=0%, write=100%", func(b *testing.B) {
		b.StopTimer()
		c := cache.NewCache[string, int](capacity)
		b.StartTimer()
		for i := 0; i < b.N; i++ {
			c.Set(fmt.Sprintf("key:%v", i), i)
		}
	})

	b.Run("empty cache: reads=50%, write=50%", func(b *testing.B) {
		b.StopTimer()
		c := cache.NewCache[string, int](capacity)
		b.StartTimer()
		for i := 0; i < b.N; i++ {
			c.Set(fmt.Sprintf("key:%v", i), i)
			c.Get(fmt.Sprintf("key:%v", i))
		}
	})

	var filledCache = cache.NewCache[string, uint64](capacity)

	for i := uint64(0); i < capacity; i++ {
		filledCache.Set(fmt.Sprintf("key:%v", i), i)
	}

	b.Run("filled cache: reads=0%, write=100%", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			filledCache.Set(fmt.Sprintf("key:%v", i), uint64(i))
		}
	})

	b.Run("filled cache: reads=50%, write=50%", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			filledCache.Set(fmt.Sprintf("key:%v", i), uint64(i))
			filledCache.Get(fmt.Sprintf("key:%v", i))
		}
	})
}

func BenchmarkSet(b *testing.B) {
	b.Run("levelitta cache", func(b *testing.B) {
		b.StopTimer()
		c := cache.NewCache[string, int](capacity)
		b.StartTimer()
		for i := 0; i < b.N; i++ {
			c.Set(fmt.Sprintf("key:%v", i), i)
		}
	})
}

func BenchmarkRead(b *testing.B) {
	b.Run("levelitta cache", func(b *testing.B) {
		b.StopTimer()
		c := cache.NewCache[string, int](capacity)
		for i := 0; i < b.N; i++ {
			c.Set(fmt.Sprintf("key:%v", i), i)
		}
		b.StartTimer()
		for i := 0; i < b.N; i++ {
			c.Get(fmt.Sprintf("key:%v", i))
		}
	})
}

func BenchmarkSetWithEvict(b *testing.B) {

}
