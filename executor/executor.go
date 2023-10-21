package executor

import (
	"context"
	"sync"
	"sync/atomic"

	"golang.org/x/sync/semaphore"
)

type Executor[T any] struct {
	sem     *semaphore.Weighted
	wg      sync.WaitGroup
	running atomic.Int64
	pending atomic.Int64
}

func New[T any](maxParallel int64) *Executor[T] {
	var sem *semaphore.Weighted
	if maxParallel > 0 {
		sem = semaphore.NewWeighted(maxParallel)
	}
	return &Executor[T]{
		sem: sem,
	}
}

type output[T any] struct {
	value T
	err   error
}

func (e *Executor[T]) Launch(ctx context.Context, fn func(context.Context) (T, error)) func(context.Context) (T, error) {
	resultC := make(chan output[T], 1)

	e.wg.Add(1)
	go func() {
		defer close(resultC)

		// wait for parallelism limits, if set
		if e.sem != nil {
			e.pending.Add(1)
			if err := e.sem.Acquire(ctx, 1); err != nil {
				e.pending.Add(-1)
				return
			}
			e.pending.Add(-1)
			defer e.sem.Release(1)
		}

		// run
		e.running.Add(1)

		wrapped := func() output[T] {
			value, err := fn(ctx)
			return output[T]{value: value, err: err}
		}

		defer e.running.Add(-1)
		select {
		case resultC <- wrapped():
		case <-ctx.Done():
		}
	}()

	// return waitable function
	return func(ctx context.Context) (T, error) {
		select {
		case out := <-resultC:
			return out.value, out.err
		case <-ctx.Done():
			var zeroVal T
			return zeroVal, ctx.Err()
		}
	}
}

func (e *Executor[T]) Wait(ctx context.Context) error {
	c := make(chan struct{})
	go func() {
		e.wg.Wait()
		close(c)
	}()
	select {
	case <-c:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (e *Executor[T]) Running() int64 {
	return e.running.Load()
}

func (e *Executor[T]) Pending() int64 {
	return e.pending.Load()
}
