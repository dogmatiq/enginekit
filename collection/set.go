package collection

import (
	"iter"
	"maps"
	"slices"
)

type set[E any] interface {
	Has(...E) bool
	Len() int
	All() iter.Seq[E]
}

type setptr[E, T any] interface {
	*T
	set[E]
	Add(...E)
}

// IsEquivalentSet returns true if a and b contain the same elements.
func IsEquivalentSet[E any, A, B set[E]](a A, b B) bool {
	if a.Len() != b.Len() {
		return false
	}

	for e := range a.All() {
		if !b.Has(e) {
			return false
		}
	}

	return true
}

// Union returns a new set containing all elements from a and b.
func Union[
	E any,
	A setptr[E, T],
	B set[E],
	T any,
](a A, b B) A {
	var union A = new(T)

	for e := range a.All() {
		union.Add(e)
	}

	for e := range b.All() {
		union.Add(e)
	}

	return union
}

// Subset returns a new set containing the elements of the given set for which
// pred returns true.
func Subset[E any, S setptr[E, T], T any](
	set S,
	pred func(E) bool,
) S {
	var subset S = new(T)

	for e := range set.All() {
		if pred(e) {
			subset.Add(e)
		}
	}

	return subset
}

// Set is an unordered set of unique T values.
type Set[E comparable] struct {
	elements map[E]struct{}
}

// NewSet returns a new [Set] containing the given elements.
func NewSet[E comparable](elements ...E) *Set[E] {
	var s Set[E]
	s.Add(elements...)
	return &s
}

// Add adds the given elements to the set.
func (s *Set[E]) Add(elements ...E) {
	if s.elements == nil {
		s.elements = make(map[E]struct{})
	}

	for _, e := range elements {
		s.elements[e] = struct{}{}
	}
}

// Remove removes the given elements from the set.
func (s *Set[E]) Remove(elements ...E) {
	for _, e := range elements {
		delete(s.elements, e)
	}
}

// Clear removes all elements from the set.
func (s *Set[E]) Clear() {
	clear(s.elements)
}

// Has returns true if the set contains all of the given elements.
func (s Set[E]) Has(elements ...E) bool {
	for _, e := range elements {
		if _, ok := s.elements[e]; !ok {
			return false
		}
	}
	return true
}

// Len returns the number of elements in the set.
func (s Set[E]) Len() int {
	return len(s.elements)
}

// All returns an iterator that yields all elements in the set.
func (s Set[E]) All() iter.Seq[E] {
	return maps.Keys(s.elements)
}

// Clone returns a shallow copy of the set.
func (s Set[E]) Clone() *Set[E] {
	return &Set[E]{maps.Clone(s.elements)}
}

// OrderedSet is an ordered set of unique T values.
type OrderedSet[E Ordered[E]] struct {
	elements []E
}

// NewOrderedSet returns a new [OrderedSet] containing the given elements.
func NewOrderedSet[E Ordered[E]](elements ...E) *OrderedSet[E] {
	var s OrderedSet[E]
	s.Add(elements...)
	return &s
}

// Add adds the given elements to the set.
func (s *OrderedSet[E]) Add(elements ...E) {
	for _, e := range elements {
		if i, ok := s.search(e); !ok {
			s.elements = slices.Insert(s.elements, i, e)
		}
	}
}

// Remove removes the given elements from the set.
func (s *OrderedSet[E]) Remove(elements ...E) {
	for _, e := range elements {
		if i, ok := s.search(e); ok {
			s.elements = slices.Delete(s.elements, i, i+1)
		}
	}
}

// Clear removes all elements from the set.
func (s *OrderedSet[E]) Clear() {
	clear(s.elements)
	s.elements = s.elements[:0]
}

// Has returns true if the set contains all of the given elements.
func (s OrderedSet[E]) Has(elements ...E) bool {
	for _, e := range elements {
		if _, ok := s.search(e); !ok {
			return false
		}
	}
	return true
}

// Len returns the number of elements in the set.
func (s OrderedSet[E]) Len() int {
	return len(s.elements)
}

// All returns an iterator that yields all elements in the set, in order.
func (s OrderedSet[E]) All() iter.Seq[E] {
	return slices.Values(s.elements)
}

// Clone returns a shallow copy of the set.
func (s OrderedSet[E]) Clone() *OrderedSet[E] {
	return &OrderedSet[E]{slices.Clone(s.elements)}
}

func (s OrderedSet[E]) search(e E) (int, bool) {
	return slices.BinarySearchFunc(s.elements, e, E.Compare)
}
