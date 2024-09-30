package collection

import (
	"iter"
	"maps"
	"slices"
)

// Set is a collection of unique elements.
type Set[E any] interface {
	Add(...E)
	Remove(...E)
	Clear()
	Has(...E) bool
	Len() int
	Elements() iter.Seq[E]
}

type setptr[E, T any] interface {
	*T
	Set[E]
}

// IsEquivalentSet returns true if a and b contain the same elements.
func IsEquivalentSet[E any](a, b Set[E]) bool {
	if a.Len() != b.Len() {
		return false
	}

	for e := range a.Elements() {
		if !b.Has(e) {
			return false
		}
	}

	return true
}

// Union returns a new set containing all elements a and b.
func Union[
	E any,
	A setptr[E, T],
	T any,
](a A, b Set[E]) A {
	var union A = new(T)

	for e := range a.Elements() {
		union.Add(e)
	}

	for e := range b.Elements() {
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

	for e := range set.Elements() {
		if pred(e) {
			subset.Add(e)
		}
	}

	return subset
}

// Clone returns a new set containing all elements of the given set.
func Clone[E any, S setptr[E, T], T any](set S) S {
	var clone S = new(T)

	for e := range set.Elements() {
		clone.Add(e)
	}

	return clone
}

// UnorderedSet is an unordered set of unique T values.
type UnorderedSet[E comparable] struct {
	elements map[E]struct{}
}

// NewUnorderedSet returns a new [UnorderedSet] containing the given elements.
func NewUnorderedSet[E comparable](elements ...E) *UnorderedSet[E] {
	var s UnorderedSet[E]
	s.Add(elements...)
	return &s
}

// Add adds the given elements to the set.
func (s *UnorderedSet[E]) Add(elements ...E) {
	if s.elements == nil {
		s.elements = make(map[E]struct{})
	}

	for _, e := range elements {
		s.elements[e] = struct{}{}
	}
}

// Remove removes the given elements from the set.
func (s *UnorderedSet[E]) Remove(elements ...E) {
	for _, e := range elements {
		delete(s.elements, e)
	}
}

// Clear removes all elements from the set.
func (s *UnorderedSet[E]) Clear() {
	clear(s.elements)
}

// Has returns true if the set contains all of the given elements.
func (s UnorderedSet[E]) Has(elements ...E) bool {
	for _, e := range elements {
		if _, ok := s.elements[e]; !ok {
			return false
		}
	}
	return true
}

// Len returns the number of elements in the set.
func (s UnorderedSet[E]) Len() int {
	return len(s.elements)
}

// Elements returns an iterator that yields all elements in the set.
func (s UnorderedSet[E]) Elements() iter.Seq[E] {
	return maps.Keys(s.elements)
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

// Elements returns an iterator that yields all elements in the set, in order.
func (s OrderedSet[E]) Elements() iter.Seq[E] {
	return slices.Values(s.elements)
}

func (s OrderedSet[E]) search(e E) (int, bool) {
	return slices.BinarySearchFunc(s.elements, e, E.Compare)
}
