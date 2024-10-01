package maps

import (
	"slices"

	"github.com/dogmatiq/enginekit/collection/constraints"
)

// OrderedByComparator is an an ordered map of keys of type K to values of type
// V with ordering defined by a separate comparitor.
type OrderedByComparator[K, V any, C constraints.Comparator[K]] struct {
	ordered[K, V]
	Comparator C
}

// NewOrderedByComparator returns an [OrderedByComparator] containing the
// given key/value pairs.
func NewOrderedByComparator[K, V any, C constraints.Comparator[K]](
	cmp C,
	pairs ...Pair[K, V],
) OrderedByComparator[K, V, C] {
	s := OrderedByComparator[K, V, C]{
		Comparator: cmp,
	}

	for _, p := range pairs {
		s.Set(p.Key, p.Value)
	}

	return s
}

// Set sets the value associated with the given key.
func (m *OrderedByComparator[K, V, C]) Set(k K, v V) {
	m.ordered.Set(m.Comparator.Compare, k, v)
}

// Remove removes the given keys from the map.
func (m *OrderedByComparator[K, V, C]) Remove(keys ...K) {
	m.ordered.Remove(m.Comparator.Compare, keys...)
}

// Has returns true if all of the given keys are in the map.
func (m OrderedByComparator[K, V, C]) Has(keys ...K) bool {
	return m.ordered.Has(m.Comparator.Compare, keys...)
}

// Get returns the value associated with the given key. It returns the zero
// value if the key is not in the map.
func (m OrderedByComparator[K, V, C]) Get(k K) V {
	v, _ := m.TryGet(k)
	return v
}

// TryGet returns the value associated with the given key, or false if the key
// is not in the map.
func (m OrderedByComparator[K, V, C]) TryGet(k K) (V, bool) {
	return m.ordered.TryGet(m.Comparator.Compare, k)
}

// Clone returns a shallow copy of the map.
func (m OrderedByComparator[K, V, C]) Clone() OrderedByComparator[K, V, C] {
	return OrderedByComparator[K, V, C]{
		slices.Clone(m.ordered),
		m.Comparator,
	}
}

// Merge returns a new map containing all key/value pairs from s and x.
//
// If a key is present in both maps, the value from x is used.
func (m OrderedByComparator[K, V, C]) Merge(x OrderedByComparator[K, V, C]) OrderedByComparator[K, V, C] {
	return OrderedByComparator[K, V, C]{
		m.ordered.Merge(m.Comparator.Compare, x.ordered),
		m.Comparator,
	}
}

// Select returns a new map containing all key/value pairs from m for which the
// given predicate returns true.
func (m OrderedByComparator[K, V, C]) Select(pred func(K, V) bool) OrderedByComparator[K, V, C] {
	return OrderedByComparator[K, V, C]{
		m.ordered.Select(pred),
		m.Comparator,
	}
}

// Project constructs a new map by applying the given transform function to each
// key/value pair in the map. If the transform function returns false, the key
// is omitted from the resulting map.
func (m OrderedByComparator[K, V, C]) Project(transform func(K, V) (K, V, bool)) OrderedByComparator[K, V, C] {
	return OrderedByComparator[K, V, C]{
		m.ordered.Project(m.Comparator.Compare, transform),
		m.Comparator,
	}
}
