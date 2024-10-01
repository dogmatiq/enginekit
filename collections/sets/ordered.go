package sets

import (
	"cmp"
	"iter"

	"github.com/dogmatiq/enginekit/collections/internal/nocopy"
)

// Ordered is an ordered set of unique T values.
type Ordered[T cmp.Ordered] struct {
	_       nocopy.NoCopy
	members []T
}

// NewOrdered returns an [Ordered] containing the given members.
func NewOrdered[T cmp.Ordered](members ...T) *Ordered[T] {
	return newOrdered[T, *Ordered[T]](members)
}

// Add adds the given members to the set.
func (s *Ordered[T]) Add(members ...T) {
	orderedAdd(s, members)
}

// Remove removes the given members from the set.
func (s *Ordered[T]) Remove(members ...T) {
	orderedRemove(s, members)
}

// Clear removes all members from the set.
func (s *Ordered[T]) Clear() {
	orderedClear[T](s)
}

// Len returns the number of members in the set.
func (s *Ordered[T]) Len() int {
	return orderedLen[T](s)
}

// Has returns true if all of the given values are members of the set.
func (s *Ordered[T]) Has(members ...T) bool {
	return orderedHas[T](s, members)
}

// IsEqual returns true if s and x have the same members.
func (s *Ordered[T]) IsEqual(x *Ordered[T]) bool {
	return orderedIsEqual[T](s, x)
}

// IsSuperset returns true if s has all of the members of x.
func (s *Ordered[T]) IsSuperset(x *Ordered[T]) bool {
	return orderedIsSuperset[T](s, x)
}

// IsSubset returns true if x has all of the members of s.
func (s *Ordered[T]) IsSubset(x *Ordered[T]) bool {
	return x.IsSuperset(s)
}

// IsStrictSuperset returns true if s has all of the members of x and at least
// one member that is not in x.
func (s *Ordered[T]) IsStrictSuperset(x *Ordered[T]) bool {
	return s.Len() > x.Len() && s.IsSuperset(x)
}

// IsStrictSubset returns true if x has all of the members of s and at least one
// member that is not in s.
func (s *Ordered[T]) IsStrictSubset(x *Ordered[T]) bool {
	return x.IsStrictSuperset(s)
}

// Clone returns a shallow copy of the set.
func (s *Ordered[T]) Clone() *Ordered[T] {
	return orderedClone[T](s)
}

// Union returns a set containing all members of s and x.
func (s *Ordered[T]) Union(x *Ordered[T]) *Ordered[T] {
	return orderedUnion[T](s, x)
}

// Select returns the subset of s containing members for which the given
// predicate function returns true.
func (s *Ordered[T]) Select(pred func(T) bool) *Ordered[T] {
	return orderedSelect[T](s, pred)
}

// All returns an iterator that yeilds all members of the set in order.
func (s *Ordered[T]) All() iter.Seq[T] {
	return orderedAll[T](s)
}

// Reverse returns an iterator that yields all members of the set in reverse
// order.
func (s *Ordered[T]) Reverse() iter.Seq[T] {
	return orderedReverse[T](s)
}

func (s *Ordered[T]) new(members []T) *Ordered[T] {
	return &Ordered[T]{
		members: members,
	}
}

func (s *Ordered[T]) ptr() *[]T {
	return &s.members
}

func (s *Ordered[T]) cmp(x, y T) int {
	return cmp.Compare(x, y)
}
