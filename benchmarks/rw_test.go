package benchmarks

import (
	"context"
	"fmt"
	"github.com/levelitta/cache"
	"github.com/levelitta/cache/benchmarks/workerpool"
	"sync"
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

	b.Run("delete item and set new", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			filledCache.Del(fmt.Sprintf("key:%v", i))
			filledCache.Set(fmt.Sprintf("key:%v", i), uint64(i))
		}
	})
}

func BenchmarkCache_ParallelOperations(b *testing.B) {
	var filledCache = cache.NewCache[string, uint64](capacity)

	for i := uint64(0); i < capacity; i++ {
		filledCache.Set(fmt.Sprintf("key:%v", i), i)
	}

	b.Run("reads=50%, write=50%", func(b *testing.B) {
		ctx, cancel := context.WithCancel(context.Background())
		wp := workerpool.NewWorkerPool(ctx, 32)

		operationsCount := b.N / 2

		wg := sync.WaitGroup{}
		wg.Add(1)
		go func() {
			defer wg.Done()

			for i := 0; i < operationsCount; i++ {
				wp.AddTask(func() {
					filledCache.Set(fmt.Sprintf("key:%v", i), uint64(i))
				})
			}
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()

			for i := 0; i < operationsCount; i++ {
				wp.AddTask(func() {
					filledCache.Get(fmt.Sprintf("key:%v", i))
				})
			}
		}()

		wg.Wait()
		cancel()
		wp.Wait()
	})

	b.Run("reads=75%, write=25%", func(b *testing.B) {
		ctx, cancel := context.WithCancel(context.Background())
		wp := workerpool.NewWorkerPool(ctx, 32)
		wg := sync.WaitGroup{}

		operationsCount := b.N / 4

		setsFn := func() {
			defer wg.Done()

			for i := 0; i < operationsCount; i++ {
				wp.AddTask(func() {
					filledCache.Set(fmt.Sprintf("key:%v", i), uint64(i))
				})
			}
		}

		getsFn := func() {
			defer wg.Done()

			for i := 0; i < operationsCount; i++ {
				wp.AddTask(func() {
					filledCache.Get(fmt.Sprintf("key:%v", i))
				})
			}
		}

		wg.Add(1)
		go setsFn()

		wg.Add(3)
		for i := 0; i < 3; i++ {
			go getsFn()
		}

		wg.Wait()
		cancel()
		wp.Wait()
	})

	b.Run("reads=95%, write=5%", func(b *testing.B) {
		ctx, cancel := context.WithCancel(context.Background())
		wp := workerpool.NewWorkerPool(ctx, 32)
		wg := sync.WaitGroup{}

		operationsCount := b.N / 20

		setsFn := func() {
			defer wg.Done()

			for i := 0; i < operationsCount; i++ {
				wp.AddTask(func() {
					filledCache.Set(fmt.Sprintf("key:%v", i), uint64(i))
				})
			}
		}

		getsFn := func(offset int) {
			defer wg.Done()

			for i := offset; i < operationsCount+offset; i++ {
				wp.AddTask(func() {
					filledCache.Get(fmt.Sprintf("key:%v", i))
				})
			}
		}

		wg.Add(20)
		for i := 0; i < 19; i++ {
			go getsFn(i * 100_000)
		}
		go setsFn()

		wg.Wait()
		cancel()
		wp.Wait()
	})

	b.Run("reads=5%, write=95%", func(b *testing.B) {
		ctx, cancel := context.WithCancel(context.Background())
		wp := workerpool.NewWorkerPool(ctx, 32)
		wg := sync.WaitGroup{}

		operationsCount := b.N / 20

		setsFn := func(offset int) {
			defer wg.Done()

			for i := offset; i < operationsCount+offset; i++ {
				wp.AddTask(func() {
					filledCache.Get(fmt.Sprintf("key:%v", i))
				})
			}

		}

		getsFn := func() {
			defer wg.Done()

			for i := 0; i < operationsCount; i++ {
				wp.AddTask(func() {
					filledCache.Set(fmt.Sprintf("key:%v", i), uint64(i))
				})
			}
		}

		wg.Add(20)
		for i := 0; i < 19; i++ {
			go setsFn(i * 100_000)
		}
		go getsFn()

		wg.Wait()
		cancel()
		wp.Wait()
	})
}

func BenchmarkSetWithEvict(b *testing.B) {

}
