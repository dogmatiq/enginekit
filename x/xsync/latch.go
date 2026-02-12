package xsync

import (
	"sync"
	"sync/atomic"
)

// Latch is a synchronization primitive that allows multiple waiters to wait for
// a single event. Latches cannot be reset.
type Latch struct {
	created, closed atomic.Bool
	m               sync.Mutex
	ch              chan struct{}
}

// Set sets the latch, allowing all waiters to proceed.
func (l *Latch) Set() {
	if l.closed.Load() {
		return
	}

	l.m.Lock()
	defer l.m.Unlock()

	if l.ch == nil {
		l.ch = make(chan struct{})
		l.created.Store(true)
	}

	close(l.ch)
	l.closed.Store(true)
}

// IsSet reports whether the latch has been set.
func (l *Latch) IsSet() bool {
	return l.closed.Load()
}

// Wait blocks until the latch is set.
func (l *Latch) Wait() {
	<-l.Chan()
}

// Chan returns a channel that is closed when the latch is set.
func (l *Latch) Chan() <-chan struct{} {
	if l.created.Load() {
		return l.ch
	}

	l.m.Lock()
	defer l.m.Unlock()

	if l.ch == nil {
		l.ch = make(chan struct{})
		l.created.Store(true)
	}

	return l.ch
}
