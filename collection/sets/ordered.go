package sets

import (
	"cmp"
	"slices"
)

// Ordered is an ordered set of unique T values.
type Ordered[T cmp.Ordered] struct {
	ordered[T]
}

// NewOrdered returns an [Ordered] containing the given members.
func NewOrdered[T cmp.Ordered](members ...T) Ordered[T] {
	var s Ordered[T]
	s.Add(members...)
	return s
}

// Add adds the given members to the set.
func (s *Ordered[T]) Add(members ...T) {
	s.ordered.Add(cmp.Compare, members...)
}

// Remove removes the given members from the set.
func (s *Ordered[T]) Remove(members ...T) {
	s.ordered.Remove(cmp.Compare, members...)
}

// Has returns true if all of the given values are members of the set.
func (s Ordered[T]) Has(members ...T) bool {
	return s.ordered.Has(cmp.Compare, members...)
}

// IsEqual returns true if s and x have the same members.
func (s Ordered[T]) IsEqual(x Ordered[T]) bool {
	return s.ordered.IsEqual(cmp.Compare, x.ordered)
}

// IsSuperset returns true if s has all of the members of x.
func (s Ordered[T]) IsSuperset(x Ordered[T]) bool {
	return s.ordered.IsSuperset(cmp.Compare, x.ordered)
}

// IsSubset returns true if x has all of the members of s.
func (s Ordered[T]) IsSubset(x Ordered[T]) bool {
	return x.IsSuperset(s)
}

// IsStrictSuperset returns true if s has all of the members of x and at least
// one member that is not in x.
func (s Ordered[T]) IsStrictSuperset(x Ordered[T]) bool {
	return s.Len() > x.Len() && s.IsSuperset(x)
}

// IsStrictSubset returns true if x has all of the members of s and at least one
// member that is not in s.
func (s Ordered[T]) IsStrictSubset(x Ordered[T]) bool {
	return x.IsStrictSuperset(s)
}

// Clone returns a shallow copy of the set.
func (s Ordered[T]) Clone() Ordered[T] {
	return Ordered[T]{slices.Clone(s.ordered)}
}

// Union returns a set containing all members of s and x.
func (s Ordered[T]) Union(x Ordered[T]) Ordered[T] {
	return Ordered[T]{s.ordered.Union(cmp.Compare, x.ordered)}
}

// Select returns the subset of s containing members for which the given
// predicate function returns true.
func (s Ordered[T]) Select(pred func(T) bool) Ordered[T] {
	return Ordered[T]{s.ordered.Select(pred)}
}
