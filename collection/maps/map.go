package maps

import (
	"iter"
	"maps"
)

// Map is an unordered map of keys of type K to values of type V.
//
// It provides the common interface implemented by all map types in this package
// to Go's built-in map type.
type Map[K comparable, V any] map[K]V

// New returns a [Map] containing the given key/value pairs.
func New[K comparable, V any](pairs ...Pair[K, V]) Map[K, V] {
	var m Map[K, V]
	for _, p := range pairs {
		m.Set(p.Key, p.Value)
	}
	return m
}

// Set sets the value associated with the given key.
func (m *Map[K, V]) Set(k K, v V) {
	if *m == nil {
		*m = Map[K, V]{}
	}

	(*m)[k] = v
}

// Remove removes the given keys from the map.
func (m *Map[K, V]) Remove(keys ...K) {
	for _, k := range keys {
		delete(*m, k)
	}
}

// Clear removes all keys from the map.
func (m *Map[K, V]) Clear() {
	clear(*m)
}

// Len returns the number of elements in the map.
func (m Map[K, V]) Len() int {
	return len(m)
}

// Has returns true if all of the given keys are in the map.
func (m Map[K, V]) Has(keys ...K) bool {
	for _, k := range keys {
		if _, ok := m[k]; !ok {
			return false
		}
	}
	return true
}

// Get returns the value associated with the given key. It returns the zero
// value if the key is not in the map.
func (m Map[K, V]) Get(k K) V {
	v, _ := m.TryGet(k)
	return v
}

// TryGet returns the value associated with the given key, or false if the key
// is not in the map.
func (m Map[K, V]) TryGet(k K) (V, bool) {
	v, ok := m[k]
	return v, ok
}

// Clone returns a shallow copy of the map.
func (m Map[K, V]) Clone() Map[K, V] {
	return maps.Clone(m)
}

// Merge returns a new map containing all key/value pairs from s and x.
//
// If a key is present in both maps, the value from x is used.
func (m Map[K, V]) Merge(x Map[K, V]) Map[K, V] {
	if m.Len() == 0 {
		return x.Clone()
	}

	if x.Len() == 0 {
		return m.Clone()
	}

	merged := m.Clone()

	for k, v := range x {
		merged[k] = v
	}

	return merged
}

// Select returns a new map containing all key/value pairs from m for which the
// given predicate returns true.
func (m Map[K, V]) Select(pred func(K, V) bool) Map[K, V] {
	var x Map[K, V]

	for k, v := range m {
		if pred(k, v) {
			if x == nil {
				x = Map[K, V]{}
			}
			x[k] = v
		}
	}

	return x
}

// Project constructs a new map by applying the given transform function to each
// key/value pair in the map. If the transform function returns false, the key
// is omitted from the resulting map.
func (m Map[K, V]) Project(transform func(K, V) (K, V, bool)) Map[K, V] {
	var x Map[K, V]

	for k, v := range m {
		if k, v, ok := transform(k, v); ok {
			x.Set(k, v)
		}
	}

	return x
}

// All returns an iterator that yields all key/value pairs in the map in no
// particular order.
func (m Map[K, V]) All() iter.Seq2[K, V] {
	return maps.All(m)
}

// Keys returns an iterator that yields all keys in the map in no particular
// order.
func (m Map[K, V]) Keys() iter.Seq[K] {
	return maps.Keys(m)
}

// Values returns an iterator that yields all values in the map in no particular
// order.
func (m Map[K, V]) Values() iter.Seq[V] {
	return maps.Values(m)
}
