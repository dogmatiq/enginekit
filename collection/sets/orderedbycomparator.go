package sets

import (
	"slices"

	"github.com/dogmatiq/enginekit/collection/constraints"
)

// OrderedByComparator is an ordered set of unique T values where the ordering
// is defined by a separate comparitor.
type OrderedByComparator[T any, C constraints.Comparator[T]] struct {
	ordered[T]

	Comparator C
}

// NewOrderedByComparator returns an [OrderedByComparator] containing the given
// members.
func NewOrderedByComparator[T any, C constraints.Comparator[T]](
	cmp C,
	members ...T,
) OrderedByComparator[T, C] {
	set := OrderedByComparator[T, C]{
		Comparator: cmp,
	}
	set.Add(members...)
	return set
}

// Add adds the given members to the set.
func (s *OrderedByComparator[T, C]) Add(members ...T) {
	s.ordered.Add(s.Comparator.Compare, members...)
}

// Remove removes the given members from the set.
func (s *OrderedByComparator[T, C]) Remove(members ...T) {
	s.ordered.Remove(s.Comparator.Compare, members...)
}

// Has returns true if all of the given values are members of the set.
func (s OrderedByComparator[T, C]) Has(members ...T) bool {
	return s.ordered.Has(s.Comparator.Compare, members...)
}

// IsEqual returns true if s and x have the same members.
func (s OrderedByComparator[T, C]) IsEqual(x OrderedByComparator[T, C]) bool {
	return s.ordered.IsEqual(s.Comparator.Compare, x.ordered)
}

// IsSuperset returns true if s has all of the members of x.
func (s OrderedByComparator[T, C]) IsSuperset(x OrderedByComparator[T, C]) bool {
	return s.ordered.IsSuperset(s.Comparator.Compare, x.ordered)
}

// IsSubset returns true if x has all of the members of s.
func (s OrderedByComparator[T, C]) IsSubset(x OrderedByComparator[T, C]) bool {
	return x.IsSuperset(s)
}

// IsStrictSuperset returns true if s has all of the members of x and at least
// one member that is not in x.
func (s OrderedByComparator[T, C]) IsStrictSuperset(x OrderedByComparator[T, C]) bool {
	return s.Len() > x.Len() && s.IsSuperset(x)
}

// IsStrictSubset returns true if x has all of the members of s and at least one
// member that is not in s.
func (s OrderedByComparator[T, C]) IsStrictSubset(x OrderedByComparator[T, C]) bool {
	return x.IsStrictSuperset(s)
}

// Clone returns a shallow copy of the set.
func (s OrderedByComparator[T, C]) Clone() OrderedByComparator[T, C] {
	return OrderedByComparator[T, C]{
		slices.Clone(s.ordered),
		s.Comparator,
	}
}

// Union returns a set containing all members of s and x.
func (s OrderedByComparator[T, C]) Union(x OrderedByComparator[T, C]) OrderedByComparator[T, C] {
	return OrderedByComparator[T, C]{
		s.ordered.Union(s.Comparator.Compare, x.ordered),
		s.Comparator,
	}
}

// Select returns the subset of s containing members for which the given
// predicate function returns true.
func (s OrderedByComparator[T, C]) Select(pred func(T) bool) OrderedByComparator[T, C] {
	return OrderedByComparator[T, C]{
		s.ordered.Select(pred),
		s.Comparator,
	}
}
