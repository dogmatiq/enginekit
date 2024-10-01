package sets

import (
	"slices"

	"github.com/dogmatiq/enginekit/collection/constraints"
)

// OrderedByMember is an ordered set of unique T values with the order defined
// by the T.Compare method.
type OrderedByMember[T constraints.Ordered[T]] struct {
	ordered[T]
}

// NewOrderedByMember returns an [OrderedByMember] containing the given members.
func NewOrderedByMember[T constraints.Ordered[T]](members ...T) OrderedByMember[T] {
	var s OrderedByMember[T]
	s.Add(members...)
	return s
}

// Add adds the given members to the set.
func (s *OrderedByMember[T]) Add(members ...T) {
	s.ordered.Add(T.Compare, members...)
}

// Remove removes the given members from the set.
func (s *OrderedByMember[T]) Remove(members ...T) {
	s.ordered.Remove(T.Compare, members...)
}

// Has returns true if all of the given values are members of the set.
func (s OrderedByMember[T]) Has(members ...T) bool {
	return s.ordered.Has(T.Compare, members...)
}

// IsEqual returns true if s and x have the same members.
func (s OrderedByMember[T]) IsEqual(x OrderedByMember[T]) bool {
	return s.ordered.IsEqual(T.Compare, x.ordered)
}

// IsSuperset returns true if s has all of the members of x.
func (s OrderedByMember[T]) IsSuperset(x OrderedByMember[T]) bool {
	return s.ordered.IsSuperset(T.Compare, x.ordered)
}

// IsSubset returns true if x has all of the members of s.
func (s OrderedByMember[T]) IsSubset(x OrderedByMember[T]) bool {
	return x.IsSuperset(s)
}

// IsStrictSuperset returns true if s has all of the members of x and at least
// one member that is not in x.
func (s OrderedByMember[T]) IsStrictSuperset(x OrderedByMember[T]) bool {
	return s.Len() > x.Len() && s.IsSuperset(x)
}

// IsStrictSubset returns true if x has all of the members of s and at least one
// member that is not in s.
func (s OrderedByMember[T]) IsStrictSubset(x OrderedByMember[T]) bool {
	return x.IsStrictSuperset(s)
}

// Clone returns a shallow copy of the set.
func (s OrderedByMember[T]) Clone() OrderedByMember[T] {
	return OrderedByMember[T]{slices.Clone(s.ordered)}
}

// Union returns a set containing all members of s and x.
func (s OrderedByMember[T]) Union(x OrderedByMember[T]) OrderedByMember[T] {
	return OrderedByMember[T]{s.ordered.Union(T.Compare, x.ordered)}
}

// Select returns the subset of s containing members for which the given
// predicate function returns true.
func (s OrderedByMember[T]) Select(pred func(T) bool) OrderedByMember[T] {
	return OrderedByMember[T]{s.ordered.Select(pred)}
}
