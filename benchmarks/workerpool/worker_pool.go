package workerpool

import (
	"context"
	"sync"
)

type Task func()

type WorkerPool struct {
	tasks chan Task
	wg    sync.WaitGroup
}

type Option func(pool *WorkerPool)

func NewWorkerPool(ctx context.Context, workersCount uint32) *WorkerPool {
	wp := &WorkerPool{
		tasks: make(chan Task, workersCount),
	}

	for i := uint32(0); i < workersCount; i++ {
		go func() {
			wp.wg.Add(1)
			defer wp.wg.Done()

			w := NewWorker()
			w.Run(ctx, wp.tasks)
		}()
	}

	return wp
}

func (wp *WorkerPool) Wait() {
	wp.wg.Wait()
}

func (wp *WorkerPool) AddTask(t Task) {
	wp.tasks <- t
}

type Worker struct{}

func NewWorker() *Worker {
	return &Worker{}
}

func (w *Worker) Run(ctx context.Context, tasks <-chan Task) {
	for {
		select {
		case t, ok := <-tasks:
			if !ok {
				return
			}
			t()
		case <-ctx.Done():
			return
		}
	}
}

func (w *Worker) Cancel(ctx context.Context, tasks <-chan Task) {
	for {
		select {
		case t, ok := <-tasks:
			if !ok {
				return
			}
			t()
		case <-ctx.Done():
			return
		}
	}
}
