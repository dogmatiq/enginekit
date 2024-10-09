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

type conflictDetector[T comparable, C Component] struct {
	m map[T]map[int]C
}

func (t *conflictDetector[T, C]) Add(
	v1 T, i int, component1 C,
	v2 T, j int, component2 C,
) bool {
	if v1 != v2 {
		return false
	}

	if t.m == nil {
		t.m = map[T]map[int]C{}
	}

	if t.m[v1] == nil {
		t.m[v2] = map[int]C{}
	}

	t.m[v1][i] = component1
	t.m[v1][i+1+j] = component2

	return true
}

func (t *conflictDetector[T, C]) All() iter.Seq2[T, []C] {
	return func(yield func(T, []C) bool) {
		for v, indices := range t.m {
			sortedIndices := slices.Sorted(maps.Keys(indices))
			sortedComponents := make([]C, len(sortedIndices))

			for i, j := range sortedIndices {
				sortedComponents[i] = indices[j]
			}

			if !yield(v, sortedComponents) {
				return
			}
		}
	}
}
