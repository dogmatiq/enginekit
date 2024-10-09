package config

import (
	"fmt"
	"iter"
	"maps"
	"slices"
)

// renderList returns a human-readable list of items.
func renderList[T any](items []T) string {
	var s string

	for i, item := range items {
		if i == len(items)-1 {
			s += " and "
		} else if i > 0 {
			s += ", "
		}
		s += fmt.Sprint(item)
	}

	return s
}

type conflictDetector[T comparable, S any] struct {
	m map[T]map[int]S
}

func (t *conflictDetector[T, S]) Add(
	v1 T, i int, src1 S,
	v2 T, j int, src2 S,
) bool {
	if v1 != v2 {
		return false
	}

	if t.m == nil {
		t.m = map[T]map[int]S{}
	}

	if t.m[v1] == nil {
		t.m[v2] = map[int]S{}
	}

	t.m[v1][i] = src1
	t.m[v1][i+1+j] = src2

	return true
}

func (t *conflictDetector[T, S]) All() iter.Seq2[T, []S] {
	return func(yield func(T, []S) bool) {
		for v, indices := range t.m {
			sortedIndices := slices.Sorted(maps.Keys(indices))
			sortedComponents := make([]S, len(sortedIndices))

			for i, j := range sortedIndices {
				sortedComponents[i] = indices[j]
			}

			if !yield(v, sortedComponents) {
				return
			}
		}
	}
}
