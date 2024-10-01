package maps

import (
	"iter"
	"slices"
)

// ordered is a helper for implementing ordered sets based on pre-sorted slices
// of key/value pairs.
type ordered[K, V any] []Pair[K, V]

func (s *ordered[K, V]) Set(cmp func(K, K) int, k K, v V) {
	if i, ok := s.search(cmp, k); ok {
		(*s)[i].Value = v
	} else {
		*s = slices.Insert(*s, i, Pair[K, V]{k, v})
	}
}

func (s *ordered[K, V]) Remove(cmp func(K, K) int, keys ...K) {
	for _, k := range keys {
		if i, ok := s.search(cmp, k); ok {
			*s = slices.Delete(*s, i, i+1)
		}
	}
}

func (s *ordered[K, V]) Clear() {
	clear(*s)
	*s = (*s)[:0]
}

func (s ordered[K, V]) Len() int {
	return len(s)
}

func (s ordered[K, V]) Has(cmp func(K, K) int, keys ...K) bool {
	for _, k := range keys {
		if _, ok := s.search(cmp, k); !ok {
			return false
		}
	}
	return true
}

func (s ordered[K, V]) TryGet(cmp func(K, K) int, k K) (V, bool) {
	if i, ok := s.search(cmp, k); ok {
		return s[i].Value, true
	}

	var zero V
	return zero, false
}

func (s ordered[K, V]) Merge(cmp func(K, K) int, x ordered[K, V]) ordered[K, V] {
	if len(s) == 0 {
		return slices.Clone(x)
	}

	if len(x) == 0 {
		return slices.Clone(s)
	}

	sIndex, xIndex := 0, 0
	sLen, xLen := len(s), len(x)

	merged := make(ordered[K, V], 0, max(sLen, xLen))

	for {
		if sIndex >= sLen {
			merged = append(merged, x[xIndex:]...)
			return merged
		}

		if xIndex >= xLen {
			merged = append(merged, s[sIndex:]...)
			return merged
		}

		sPair := s[sIndex]
		xPair := x[xIndex]

		c := cmp(sPair.Key, xPair.Key)

		if c < 0 {
			merged = append(merged, sPair)
			sIndex++
		} else if c > 0 {
			merged = append(merged, xPair)
			xIndex++
		} else {
			merged = append(merged, xPair)
			sIndex++
			xIndex++
		}
	}
}

func (s ordered[K, V]) Select(pred func(K, V) bool) ordered[K, V] {
	var x ordered[K, V]

	for _, p := range s {
		if pred(p.Key, p.Value) {
			x = append(x, p)
		}
	}

	return x
}

func (s ordered[K, V]) Project(
	cmp func(K, K) int,
	transform func(K, V) (K, V, bool),
) ordered[K, V] {
	var x ordered[K, V]

	for _, p := range s {
		if k, v, ok := transform(p.Key, p.Value); ok {
			x.Set(cmp, k, v)
		}
	}

	return x
}

func (s ordered[K, V]) All() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for _, p := range s {
			if !yield(p.Key, p.Value) {
				return
			}
		}
	}
}

func (s ordered[K, V]) Reverse() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for i := len(s) - 1; i >= 0; i-- {
			p := s[i]
			if !yield(p.Key, p.Value) {
				return
			}
		}
	}
}

func (s ordered[K, V]) Keys() iter.Seq[K] {
	return func(yield func(K) bool) {
		for _, p := range s {
			if !yield(p.Key) {
				return
			}
		}
	}
}

func (s ordered[K, V]) ReverseKeys() iter.Seq[K] {
	return func(yield func(K) bool) {
		for i := len(s) - 1; i >= 0; i-- {
			if !yield(s[i].Key) {
				return
			}
		}
	}
}

func (s ordered[K, V]) Values() iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, p := range s {
			if !yield(p.Value) {
				return
			}
		}
	}
}

func (s ordered[K, V]) ReverseValues() iter.Seq[V] {
	return func(yield func(V) bool) {
		for i := len(s) - 1; i >= 0; i-- {
			if !yield(s[i].Value) {
				return
			}
		}
	}
}

func (s ordered[K, V]) search(cmp func(K, K) int, k K) (int, bool) {
	return slices.BinarySearchFunc(
		s,
		k,
		func(p Pair[K, V], k K) int {
			return cmp(p.Key, k)
		},
	)
}
