package uuidpb

import (
	"slices"
)

// Set is an ordered set of UUIDs.
type Set []*UUID

// Add adds id to the set.
func (s *Set) Add(id *UUID) {
	if i, ok := s.search(id); !ok {
		*s = slices.Insert(*s, i, id)
	}
}

// Has returns true if id is in the set.
func (s Set) Has(id *UUID) bool {
	_, ok := s.search(id)
	return ok
}

// Delete removes id from the set.
func (s *Set) Delete(id *UUID) {
	if i, ok := s.search(id); ok {
		*s = slices.Delete(*s, i, i+1)
	}
}

func (s Set) search(id *UUID) (int, bool) {
	return slices.BinarySearchFunc(
		s,
		id,
		(*UUID).Compare,
	)
}
