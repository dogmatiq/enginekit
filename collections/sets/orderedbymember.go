package sets

import (
	"iter"

	"github.com/dogmatiq/enginekit/collections/constraints"
)

// OrderedByMember is an ordered set of unique T values with the order defined
// by the T.Compare method.
type OrderedByMember[T constraints.Ordered[T]] struct {
	members []T
}

// NewOrderedByMember returns an [OrderedByMember] containing the given members.
func NewOrderedByMember[T constraints.Ordered[T]](members ...T) *OrderedByMember[T] {
	return newOrdered[T, *OrderedByMember[T]](members)
}

// Add adds the given members to the set.
func (s *OrderedByMember[T]) Add(members ...T) {
	orderedAdd(s, members)
}

// Remove removes the given members from the set.
func (s *OrderedByMember[T]) Remove(members ...T) {
	orderedRemove(s, members)
}

// Clear removes all members from the set.
func (s *OrderedByMember[T]) Clear() {
	orderedClear[T](s)
}

// Len returns the number of members in the set.
func (s *OrderedByMember[T]) Len() int {
	return orderedLen[T](s)
}

// Has returns true if all of the given values are members of the set.
func (s *OrderedByMember[T]) Has(members ...T) bool {
	return orderedHas[T](s, members)
}

// IsEqual returns true if s and x have the same members.
func (s *OrderedByMember[T]) IsEqual(x *OrderedByMember[T]) bool {
	return orderedIsEqual[T](s, x)
}

// IsSuperset returns true if s has all of the members of x.
func (s *OrderedByMember[T]) IsSuperset(x *OrderedByMember[T]) bool {
	return orderedIsSuperset[T](s, x)
}

// IsSubset returns true if x has all of the members of s.
func (s *OrderedByMember[T]) IsSubset(x *OrderedByMember[T]) bool {
	return x.IsSuperset(s)
}

// IsStrictSuperset returns true if s has all of the members of x and at least
// one member that is not in x.
func (s *OrderedByMember[T]) IsStrictSuperset(x *OrderedByMember[T]) bool {
	return s.Len() > x.Len() && s.IsSuperset(x)
}

// IsStrictSubset returns true if x has all of the members of s and at least one
// member that is not in s.
func (s *OrderedByMember[T]) IsStrictSubset(x *OrderedByMember[T]) bool {
	return x.IsStrictSuperset(s)
}

// Clone returns a shallow copy of the set.
func (s *OrderedByMember[T]) Clone() *OrderedByMember[T] {
	return orderedClone[T](s)
}

// Union returns a set containing all members of s and x.
func (s *OrderedByMember[T]) Union(x *OrderedByMember[T]) *OrderedByMember[T] {
	return orderedUnion[T](s, x)
}

// Select returns the subset of s containing members for which the given
// predicate function returns true.
func (s *OrderedByMember[T]) Select(pred func(T) bool) *OrderedByMember[T] {
	return orderedSelect[T](s, pred)
}

// All returns an iterator that yeilds all members of the set in order.
func (s *OrderedByMember[T]) All() iter.Seq[T] {
	return orderedAll[T](s)
}

// Reverse returns an iterator that yields all members of the set in reverse
// order.
func (s *OrderedByMember[T]) Reverse() iter.Seq[T] {
	return orderedReverse[T](s)
}

func (s *OrderedByMember[T]) new(members []T) *OrderedByMember[T] {
	return &OrderedByMember[T]{
		members: members,
	}
}

func (s *OrderedByMember[T]) ptr() *[]T {
	return &s.members
}

func (s *OrderedByMember[T]) cmp(x, y T) int {
	return x.Compare(y)
}