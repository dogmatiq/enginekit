package uuidpb

import (
	"iter"
)

// Map is a map of [*UUID] to values of type V, where key equality is based on
// the UUID value, not the pointer address.
type Map[V any] struct {
	m map[plain]V
}

// Len returns the number of key/value pairs in the map.
func (m *Map[V]) Len() int {
	return len(m.m)
}

// All returns an iterator that yields each key/value pair in the map, in no
// particular order.
func (m *Map[V]) All() iter.Seq2[*UUID, V] {
	return func(yield func(*UUID, V) bool) {
		for k, v := range m.m {
			if !yield(k.proto(), v) {
				return
			}
		}
	}
}

// Keys returns an iterator that yields each key in the map, in no particular
// order.
func (m *Map[V]) Keys() iter.Seq[*UUID] {
	return func(yield func(*UUID) bool) {
		for k := range m.m {
			if !yield(k.proto()) {
				return
			}
		}
	}
}

// Values returns an iterator that yields each value in the map, in no
// particular order.
func (m *Map[V]) Values() iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, v := range m.m {
			if !yield(v) {
				return
			}
		}
	}
}

// Set associates v with k.
func (m *Map[V]) Set(k *UUID, v V) {
	if m.m == nil {
		m.m = map[plain]V{}
	}
	m.m[k.plain()] = v
}

// Get returns the value associated with k, or the zero-value if k is not
// present in the map.
func (m *Map[V]) Get(k *UUID) V {
	return m.m[k.plain()]
}

// Has returns true if k is present in the map.
func (m *Map[V]) Has(k *UUID) bool {
	_, ok := m.m[k.plain()]
	return ok
}

// TryGet returns the value associated with k, or false if k is not present in
// the map.
func (m *Map[V]) TryGet(k *UUID) (V, bool) {
	v, ok := m.m[k.plain()]
	return v, ok
}

// Delete removes the value associated with k.
func (m *Map[V]) Delete(k *UUID) {
	delete(m.m, k.plain())
}
