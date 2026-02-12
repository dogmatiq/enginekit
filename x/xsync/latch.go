package xsync

import (
	"sync"
	"sync/atomic"
	"weak"
)

// Latch is a synchronization primitive that allows multiple waiters to wait for
// a single event. Latches cannot be reset.
type Latch struct {
	hasChan, isSet atomic.Bool

	m          sync.Mutex
	ch         chan struct{}
	downstream []weak.Pointer[Latch]
}

// Set sets the latch, allowing all waiters to proceed.
//
// It returns true if the latch was set by this call, or false if it was already
// set.
func (l *Latch) Set() bool {
	if l.isSet.Load() {
		return false
	}

	l.m.Lock()
	defer l.m.Unlock()

	if l.isSet.Load() {
		return false
	}

	if l.ch == nil {
		l.ch = make(chan struct{})
		l.hasChan.Store(true)
	}

	close(l.ch)
	l.isSet.Store(true)

	for _, p := range l.downstream {
		if x := p.Value(); x != nil {
			x.Set()
		}
	}

	return true
}

// Link creates a unidirectional relationship that sets l when the upstream
// latch is closed.
func (l *Latch) Link(upstream *Latch) {
	if upstream.isSet.Load() {
		l.Set()
		return
	}

	upstream.m.Lock()
	defer upstream.m.Unlock()

	if upstream.isSet.Load() {
		l.Set()
		return
	}

	p := weak.Make(l)

	for i, x := range upstream.downstream {
		if x.Value() == nil {
			upstream.downstream[i] = p
			return
		}
	}

	upstream.downstream = append(upstream.downstream, p)
}

// IsSet reports whether the latch has been set.
func (l *Latch) IsSet() bool {
	return l.isSet.Load()
}

// Wait blocks until the latch is set.
func (l *Latch) Wait() {
	<-l.Chan()
}

// Chan returns a channel that is closed when the latch is set.
func (l *Latch) Chan() <-chan struct{} {
	if l.hasChan.Load() {
		return l.ch
	}

	l.m.Lock()
	defer l.m.Unlock()

	if l.ch == nil {
		l.ch = make(chan struct{})
		l.hasChan.Store(true)
	}

	return l.ch
}
