package maps

import (
	"cmp"
	"slices"
)

// Ordered is an ordered map of keys of type K to values of type V.
type Ordered[K cmp.Ordered, V any] struct {
	ordered[K, V]
}

// NewOrdered returns an [Ordered] containing the given key/value pairs.
func NewOrdered[K cmp.Ordered, V any](pairs ...Pair[K, V]) Ordered[K, V] {
	var s Ordered[K, V]
	for _, p := range pairs {
		s.Set(p.Key, p.Value)
	}
	return s
}

// Set sets the value associated with the given key.
func (m *Ordered[K, V]) Set(k K, v V) {
	m.ordered.Set(cmp.Compare, k, v)
}

// Remove removes the given keys from the map.
func (m *Ordered[K, V]) Remove(keys ...K) {
	m.ordered.Remove(cmp.Compare, keys...)
}

// Has returns true if all of the given keys are in the map.
func (m Ordered[K, V]) Has(keys ...K) bool {
	return m.ordered.Has(cmp.Compare, keys...)
}

// Get returns the value associated with the given key. It returns the zero
// value if the key is not in the map.
func (m Ordered[K, V]) Get(k K) V {
	v, _ := m.TryGet(k)
	return v
}

// TryGet returns the value associated with the given key, or false if the key
// is not in the map.
func (m Ordered[K, V]) TryGet(k K) (V, bool) {
	return m.ordered.TryGet(cmp.Compare, k)
}

// Clone returns a shallow copy of the map.
func (m Ordered[K, V]) Clone() Ordered[K, V] {
	return Ordered[K, V]{slices.Clone(m.ordered)}
}

// Merge returns a new map containing all key/value pairs from s and x.
//
// If a key is present in both maps, the value from x is used.
func (m Ordered[K, V]) Merge(x Ordered[K, V]) Ordered[K, V] {
	return Ordered[K, V]{m.ordered.Merge(cmp.Compare, x.ordered)}
}

// Select returns a new map containing all key/value pairs from m for which the
// given predicate returns true.
func (m Ordered[K, V]) Select(pred func(K, V) bool) Ordered[K, V] {
	return Ordered[K, V]{m.ordered.Select(pred)}
}

// Project constructs a new map by applying the given transform function to each
// key/value pair in the map. If the transform function returns false, the key
// is omitted from the resulting map.
func (m Ordered[K, V]) Project(transform func(K, V) (K, V, bool)) Ordered[K, V] {
	return Ordered[K, V]{m.ordered.Project(cmp.Compare, transform)}
}
