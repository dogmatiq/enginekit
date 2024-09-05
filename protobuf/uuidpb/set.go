package uuidpb

import (
	"iter"
	"slices"
)

// OrderedSet is an ordered set of UUIDs.
type OrderedSet struct {
	elems []plain
}

// Len returns the number of elements in s.
func (s *OrderedSet) Len() int {
	return len(s.elems)
}

// All returns an iterator that yields each element in the set in order.
func (s *OrderedSet) All() iter.Seq[*UUID] {
	return func(yield func(*UUID) bool) {
		for _, id := range s.elems {
			if !yield(id.proto()) {
				return
			}
		}
	}
}

// Add adds id to s.
func (s *OrderedSet) Add(id *UUID) {
	if i, ok := s.search(id); !ok {
		s.elems = slices.Insert(s.elems, i, id.plain())
	}
}

// Has returns true if id is an element of s.
func (s OrderedSet) Has(id *UUID) bool {
	_, ok := s.search(id)
	return ok
}

// Delete removes id from s.
func (s *OrderedSet) Delete(id *UUID) {
	if i, ok := s.search(id); ok {
		s.elems = slices.Delete(s.elems, i, i+1)
	}
}

func (s *OrderedSet) search(id *UUID) (int, bool) {
	return slices.BinarySearchFunc(
		s.elems,
		id,
		func(a plain, b *UUID) int {
			return compare(
				a.upper, a.lower,
				b.GetUpper(), b.GetLower(),
			)
		},
	)
}
