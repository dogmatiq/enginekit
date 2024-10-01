package sets

import (
	"iter"
	"slices"
)

// ordered is a helper for implementing ordered sets based on pre-sorted slices.
type ordered[T any] []T

func (s *ordered[T]) Add(cmp func(T, T) int, members ...T) {
	for _, m := range members {
		if i, ok := slices.BinarySearchFunc(*s, m, cmp); !ok {
			*s = slices.Insert(*s, i, m)
		}
	}
}

func (s *ordered[T]) Remove(cmp func(T, T) int, members ...T) {
	for _, m := range members {
		if i, ok := slices.BinarySearchFunc(*s, m, cmp); ok {
			*s = slices.Delete(*s, i, i+1)
		}
	}
}

func (s *ordered[T]) Clear() {
	clear(*s)
	*s = (*s)[:0]
}

func (s ordered[T]) Len() int {
	return len(s)
}

func (s ordered[T]) Has(cmp func(T, T) int, members ...T) bool {
	for _, m := range members {
		if _, ok := slices.BinarySearchFunc(s, m, cmp); !ok {
			return false
		}
	}
	return true
}

func (s ordered[T]) IsEqual(cmp func(T, T) int, x ordered[T]) bool {
	if len(s) != len(x) {
		return false
	}

	for i, m := range s {
		if cmp(x[i], m) != 0 {
			return false
		}
	}

	return true
}

func (s ordered[T]) IsSuperset(cmp func(T, T) int, x ordered[T]) bool {
	if len(s) == len(x) {
		return s.IsEqual(cmp, x)
	}

	if len(s) < len(x) {
		return false
	}

	supIndex, subIndex := 0, 0

	for {
		if subIndex >= len(x) {
			return true
		}

		if supIndex >= len(s) {
			return false
		}

		supMember := s[supIndex]
		subMember := x[subIndex]

		c := cmp(subMember, supMember)

		if c < 0 {
			return false
		}

		if c == 0 {
			subIndex++
		}

		supIndex++
	}
}

func (s ordered[T]) Union(cmp func(T, T) int, x ordered[T]) ordered[T] {
	if len(s) == 0 {
		return slices.Clone(x)
	}

	if len(x) == 0 {
		return slices.Clone(s)
	}

	sIndex, xIndex := 0, 0
	sLen, xLen := len(s), len(x)

	union := make(ordered[T], 0, max(sLen, xLen))

	for {
		if sIndex >= sLen {
			union = append(union, x[xIndex:]...)
			return union
		}

		if xIndex >= xLen {
			union = append(union, s[sIndex:]...)
			return union
		}

		sMember := s[sIndex]
		xMember := x[xIndex]

		c := cmp(sMember, xMember)

		if c < 0 {
			union = append(union, sMember)
			sIndex++
		} else if c > 0 {
			union = append(union, xMember)
			xIndex++
		} else {
			union = append(union, sMember)
			sIndex++
			xIndex++
		}
	}
}

func (s ordered[T]) Select(pred func(T) bool) ordered[T] {
	var subset ordered[T]

	for _, m := range s {
		if pred(m) {
			subset = append(subset, m)
		}
	}

	return subset
}

func (s ordered[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, m := range s {
			if !yield(m) {
				return
			}
		}
	}
}

func (s ordered[T]) Reverse() iter.Seq[T] {
	return func(yield func(T) bool) {
		for i := len(s) - 1; i >= 0; i-- {
			if !yield(s[i]) {
				return
			}
		}
	}
}
