package sets

import (
	"iter"
	"maps"
)

// Set is an unordered set of unique T values.
type Set[T comparable] struct {
	members map[T]struct{}
}

// New returns a [Set] containing the given members.
func New[T comparable](members ...T) Set[T] {
	var s Set[T]
	s.Add(members...)
	return s
}

// Add adds the given members to the set.
func (s *Set[T]) Add(members ...T) {
	if s.members == nil {
		s.members = make(map[T]struct{}, len(members))
	}
	for _, m := range members {
		s.members[m] = struct{}{}
	}
}

// Remove removes the given members from the set.
func (s *Set[T]) Remove(members ...T) {
	for _, m := range members {
		delete(s.members, m)
	}
}

// Clear removes all members from the set.
func (s *Set[T]) Clear() {
	clear(s.members)
}

// Len returns the number of members in the set.
func (s Set[T]) Len() int {
	return len(s.members)
}

// Has returns true if all of the given values are members of the set.
func (s Set[T]) Has(members ...T) bool {
	for _, m := range members {
		if _, ok := s.members[m]; !ok {
			return false
		}
	}
	return true
}

// IsEqual returns true if s and x have the same members.
func (s Set[T]) IsEqual(x Set[T]) bool {
	if s.Len() != x.Len() {
		return false
	}

	for m := range s.members {
		if _, ok := x.members[m]; !ok {
			return false
		}
	}

	return true
}

// IsSuperset returns true if s has all of the members of x.
func (s Set[T]) IsSuperset(x Set[T]) bool {
	if s.Len() < x.Len() {
		return false
	}

	for m := range x.members {
		if _, ok := s.members[m]; !ok {
			return false
		}
	}

	return true
}

// IsSubset returns true if x has all of the members of s.
func (s Set[T]) IsSubset(x Set[T]) bool {
	return x.IsSuperset(s)
}

// IsStrictSuperset returns true if s has all of the members of x and at least
// one member that is not in x.
func (s Set[T]) IsStrictSuperset(x Set[T]) bool {
	return s.Len() > x.Len() && s.IsSuperset(x)
}

// IsStrictSubset returns true if x has all of the members of s and at least one
// member that is not in s.
func (s Set[T]) IsStrictSubset(x Set[T]) bool {
	return x.IsStrictSuperset(s)
}

// Clone returns a shallow copy of the set.
func (s Set[T]) Clone() Set[T] {
	return Set[T]{maps.Clone(s.members)}
}

// Union returns a set containing all members of s and x.
func (s Set[T]) Union(x Set[T]) Set[T] {
	if s.Len() == 0 {
		return x.Clone()
	}

	if x.Len() == 0 {
		return s.Clone()
	}

	smaller, larger := s, x
	if larger.Len() < smaller.Len() {
		smaller, larger = larger, smaller
	}

	union := larger.Clone()

	for m := range smaller.members {
		union.members[m] = struct{}{}
	}

	return union
}

// Select returns the subset of s containing members for which the given
// predicate function returns true.
func (s Set[T]) Select(pred func(T) bool) Set[T] {
	subset := Set[T]{
		members: map[T]struct{}{},
	}

	for m := range s.members {
		if pred(m) {
			subset.members[m] = struct{}{}
		}
	}

	return subset
}

// All returns an iterator that yields all members of the set in no particular
// order.
func (s Set[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		for m := range s.members {
			if !yield(m) {
				return
			}
		}
	}
}
