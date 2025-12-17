package stubs

import (
	"strconv"
	"sync/atomic"

	"github.com/dogmatiq/enginekit/protobuf/uuidpb"
)

// UUIDSequence is a generator of deterministic UUIDs for use in tests.
type UUIDSequence struct {
	namespace atomic.Pointer[uuidpb.UUID]
	counter   atomic.Uint64
}

// Next returns the next UUID in the sequence.
func (g *UUIDSequence) Next() *uuidpb.UUID {
	n := g.counter.Add(1) - 1
	return g.At(n)
}

// At returns the UUID at the given position in the sequence.
func (g *UUIDSequence) At(n uint64) *uuidpb.UUID {
	ns := g.namespace.Load()

	if ns == nil {
		ns = uuidpb.Generate()

		if !g.namespace.CompareAndSwap(nil, ns) {
			ns = g.namespace.Load()
		}
	}

	return uuidpb.Derive(
		ns,
		strconv.FormatUint(n, 10),
	)
}
