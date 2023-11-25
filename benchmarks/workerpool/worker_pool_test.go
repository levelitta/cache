package workerpool

import (
	"context"
	"github.com/stretchr/testify/require"
	"sync"
	"sync/atomic"
	"testing"
)

func TestWorkerPool_AddTask(t *testing.T) {

	t.Run("success - task execute", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		wp := NewWorkerPool(ctx, 3)
		wg := sync.WaitGroup{}

		execCount := atomic.Int32{}

		for i := 0; i < 10; i++ {
			wg.Add(1)
			wp.AddTask(func() {
				defer wg.Done()
				execCount.Add(1)
			})
		}

		wg.Wait()

		require.Equal(t, int32(10), execCount.Load())

		cancel()
	})
}
