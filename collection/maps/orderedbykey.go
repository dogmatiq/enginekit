package maps

import (
	"slices"

	"github.com/dogmatiq/enginekit/collection/constraints"
)

// OrderedByKey is an an ordered map of keys of type K to values of type V with
// ordering defined by the K.Compare method.
type OrderedByKey[K constraints.Ordered[K], V any] struct {
	ordered[K, V]
}

// NewOrderedByKey returns an [OrderedByKey] containing the given key/value
// pairs.
func NewOrderedByKey[K constraints.Ordered[K], V any](pairs ...Pair[K, V]) OrderedByKey[K, V] {
	var s OrderedByKey[K, V]
	for _, p := range pairs {
		s.Set(p.Key, p.Value)
	}
	return s
}

// Set sets the value associated with the given key.
func (m *OrderedByKey[K, V]) Set(k K, v V) {
	m.ordered.Set(K.Compare, k, v)
}

// Remove removes the given keys from the map.
func (m *OrderedByKey[K, V]) Remove(keys ...K) {
	m.ordered.Remove(K.Compare, keys...)
}

// Has returns true if all of the given keys are in the map.
func (m OrderedByKey[K, V]) Has(keys ...K) bool {
	return m.ordered.Has(K.Compare, keys...)
}

// Get returns the value associated with the given key. It returns the zero
// value if the key is not in the map.
func (m OrderedByKey[K, V]) Get(k K) V {
	v, _ := m.TryGet(k)
	return v
}

// TryGet returns the value associated with the given key, or false if the key
// is not in the map.
func (m OrderedByKey[K, V]) TryGet(k K) (V, bool) {
	return m.ordered.TryGet(K.Compare, k)
}

// Clone returns a shallow copy of the map.
func (m OrderedByKey[K, V]) Clone() OrderedByKey[K, V] {
	return OrderedByKey[K, V]{slices.Clone(m.ordered)}
}

// Merge returns a new map containing all key/value pairs from s and x.
//
// If a key is present in both maps, the value from x is used.
func (m OrderedByKey[K, V]) Merge(x OrderedByKey[K, V]) OrderedByKey[K, V] {
	return OrderedByKey[K, V]{m.ordered.Merge(K.Compare, x.ordered)}
}

// Select returns a new map containing all key/value pairs from m for which the
// given predicate returns true.
func (m OrderedByKey[K, V]) Select(pred func(K, V) bool) OrderedByKey[K, V] {
	return OrderedByKey[K, V]{m.ordered.Select(pred)}
}

// Project constructs a new map by applying the given transform function to each
// key/value pair in the map. If the transform function returns false, the key
// is omitted from the resulting map.
func (m OrderedByKey[K, V]) Project(transform func(K, V) (K, V, bool)) OrderedByKey[K, V] {
	return OrderedByKey[K, V]{m.ordered.Project(K.Compare, transform)}
}
