package xsync

import (
	"sync/atomic"
)

// Observable is an atomic value that can be observed for changes.
type Observable[T any] struct {
	state atomic.Pointer[observableState[T]]
}

// Load returns the current value and a channel that is closed when that is
// replaced by a subsequent call to Store.
func (o *Observable[T]) Load() (T, <-chan struct{}) {
	if prev := o.state.Load(); prev != nil {
		return prev.value, prev.replaced
	}

	next := &observableState[T]{
		replaced: make(chan struct{}),
	}

	if o.state.CompareAndSwap(nil, next) {
		return next.value, next.replaced
	}

	prev := o.state.Load()
	return prev.value, prev.replaced
}

// Store updates the value.
func (o *Observable[T]) Store(v T) {
	next := &observableState[T]{
		value:    v,
		replaced: make(chan struct{}),
	}

retry:
	prev := o.state.Load()

	if !o.state.CompareAndSwap(prev, next) {
		goto retry
	}

	if prev != nil {
		close(prev.replaced)
	}
}

// ChangeObservable changes the value of o to v, unless v is already its current
// value.
//
// It returns true if the value was changed, in which case the channel returned
// by any prior call to [Observable.Load] is closed.
func ChangeObservable[T comparable](o *Observable[T], v T) bool {
	var next *observableState[T]

	for {
		prev := o.state.Load()

		var current T
		if prev != nil {
			current = prev.value
		}

		if current == v {
			return false
		}

		if next == nil {
			next = &observableState[T]{
				value:    v,
				replaced: make(chan struct{}),
			}
		}

		if o.state.CompareAndSwap(prev, next) {
			if prev != nil {
				close(prev.replaced)
			}

			return true
		}
	}
}

type observableState[T any] struct {
	value    T
	replaced chan struct{}
}
