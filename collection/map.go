package collection

import (
	"iter"
	"slices"
)

// OrderedMap is an map of keys to values that maintains the order of the keys.
type OrderedMap[K Ordered[K], V any] struct {
	elements []pair[K, V]
}

type pair[K, V any] struct {
	K K
	V V
}

// Set sets the value associated with k to v.
func (m *OrderedMap[K, V]) Set(k K, v V) {
	if i, ok := m.search(k); ok {
		m.elements[i].V = v
	} else {
		p := pair[K, V]{k, v}
		m.elements = slices.Insert(m.elements, i, p)
	}
}

// Remove removes all of the given keys from the map.
func (m *OrderedMap[K, V]) Remove(keys ...K) {
	for _, k := range keys {
		if i, ok := m.search(k); ok {
			m.elements = slices.Delete(m.elements, i, i+1)
		}
	}
}

// Clear removes all keys from the map.
func (m *OrderedMap[K, V]) Clear() {
	clear(m.elements)
	m.elements = m.elements[:0]
}

// Get returns the value associated with the given key. It returns the zero
// value if the key is not in the map.
func (m OrderedMap[K, V]) Get(k K) V {
	v, _ := m.TryGet(k)
	return v
}

// TryGet returns the value associated with the given key, or false if the key
// is not in the map.
func (m OrderedMap[K, V]) TryGet(k K) (V, bool) {
	if i, ok := m.search(k); ok {
		return m.elements[i].V, true
	}

	var zero V
	return zero, false
}

// Has returns true if the map contains all of the given keys.
func (m OrderedMap[K, V]) Has(keys ...K) bool {
	for _, k := range keys {
		if _, ok := m.search(k); !ok {
			return false
		}
	}
	return true
}

// Len returns the number of elements in the map.
func (m OrderedMap[K, V]) Len() int {
	return len(m.elements)
}

// Elements returns an iterator that yields all the key/value pairs in the map,
// in key order.
func (m OrderedMap[K, V]) Elements() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for _, p := range m.elements {
			if !yield(p.K, p.V) {
				return
			}
		}
	}
}

// Keys returns an iterator that yields all the keys in the map, in order.
func (m OrderedMap[K, V]) Keys() iter.Seq[K] {
	return func(yield func(K) bool) {
		for _, p := range m.elements {
			if !yield(p.K) {
				return
			}
		}
	}
}

// Values returns an iterator that yields all the values in the map in the order
// of their keys.
func (m OrderedMap[K, V]) Values() iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, p := range m.elements {
			if !yield(p.V) {
				return
			}
		}
	}
}

func (m OrderedMap[K, V]) search(k K) (int, bool) {
	return slices.BinarySearchFunc(
		m.elements,
		k,
		func(p pair[K, V], k K) int {
			return p.K.Compare(k)
		},
	)
}
