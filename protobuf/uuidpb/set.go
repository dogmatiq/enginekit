package uuidpb

import (
	"iter"
	"maps"
)

// Set is a collection of [UUID] values.
type Set struct {
	m map[key]struct{}
}

// Has reports whether the set contains the given [UUID].
func (s *Set) Has(v *UUID) bool {
	if s.Len() == 0 {
		return false
	}

	_, ok := s.m[asKey(v)]
	return ok
}

// Add adds the given [UUID] to the set.
func (s *Set) Add(v *UUID) {
	if s.m == nil {
		s.m = map[key]struct{}{}
	}

	s.m[asKey(v)] = struct{}{}
}

// Delete removes the given [UUID] from the set.
func (s *Set) Delete(v *UUID) {
	if s.Len() != 0 {
		delete(s.m, asKey(v))
	}
}

// All yields all members of the set.
func (s *Set) All() iter.Seq[*UUID] {
	return func(yield func(*UUID) bool) {
		if s != nil {
			for v := range s.m {
				if !yield(v.asUUID()) {
					return
				}
			}
		}
	}
}

// Len returns the number of members in the set.
func (s *Set) Len() int {
	if s == nil {
		return 0
	}

	return len(s.m)
}

// Clear removes all members from the set.
func (s *Set) Clear() {
	if s != nil {
		clear(s.m)
	}
}

// Clone returns a shallow copy of the set.
func (s *Set) Clone() *Set {
	if s == nil {
		return nil
	}

	return &Set{
		m: maps.Clone(s.m),
	}
}
