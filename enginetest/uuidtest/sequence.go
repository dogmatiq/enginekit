package uuidtest

import (
	"encoding/binary"
	"sync/atomic"

	"github.com/dogmatiq/enginekit/protobuf/uuidpb"
)

// Sequence is a sequence of deterministic UUIDs.
type Sequence struct {
	ns *uuidpb.UUID
	n  atomic.Uint64
}

// NewSequence returns a new deterministic UUID sequence.
func NewSequence() *Sequence {
	return &Sequence{
		ns: uuidpb.Generate(),
	}
}

// Next returns the next UUID in the sequence.
func (s *Sequence) Next() *uuidpb.UUID {
	n := s.n.Add(1) - 1
	return s.Nth(n)
}

// Nth returns the n'th UUID in the sequence.
func (s *Sequence) Nth(n uint64) *uuidpb.UUID {
	var data [8]byte
	binary.BigEndian.PutUint64(data[:], n)
	return uuidpb.Derive(s.ns, data[:])
}

// IsNth returns true if x is the n'th UUID in the sequence.
func (s *Sequence) IsNth(x *uuidpb.UUID, n uint64) bool {
	return x.Equal(s.Nth(n))
}
