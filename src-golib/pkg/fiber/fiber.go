package fiber

import (
	"sync/atomic"
)

const (
	Started = 1 << iota
	Running
)

type fiber[T any] struct {
	status uint32
	fn     func(SuspendFunc[T])
	in     chan T
	out    chan T
}

type Fiber[T any] interface {
	Start() T
	Resume(v T) T

	IsStarted() bool
	IsRunning() bool
	IsSuspended() bool
	IsTerminated() bool
}

type SuspendFunc[T any] func(T) T

func New[T any](fn func(SuspendFunc[T])) Fiber[T] {
	return &fiber[T]{
		fn:  fn,
		in:  make(chan T),
		out: make(chan T),
	}
}
func (f *fiber[T]) Start() (v T) {
	if atomic.SwapUint32(&f.status, Started|Running) != 0 {
		var zero T
		return zero
	}
	go func() {
		f.fn(f.suspend)
		close(f.in)
		close(f.out)
		f.status = 0 // Terminated
	}()
	return <-f.out
}

func (f *fiber[T]) Resume(v T) T {
	if atomic.SwapUint32(&f.status, f.status|Running) != Started {
		var zero T
		return zero
	}
	f.in <- v
	return <-f.out
}

func (f *fiber[T]) suspend(v T) T {
	f.status &^= Running
	f.out <- v
	return <-f.in
}

func (f *fiber[T]) IsStarted() bool {
	return f.status&Started > 0
}

func (f *fiber[T]) IsRunning() bool {
	return f.status&Running > 0
}

func (f *fiber[T]) IsSuspended() bool {
	return !f.IsRunning()
}

func (f *fiber[T]) IsTerminated() bool {
	return !f.IsStarted()
}
