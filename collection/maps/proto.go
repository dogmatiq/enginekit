package maps

import (
	"iter"
	"reflect"

	"google.golang.org/protobuf/proto"
)

// Proto is a map Protocol Buffers messages of type K to values of type V.
//
// K must be a pointer type that implements [proto.Message].
//
// Key equality is determined based on the serialized form of the message, and
// so is subject to the caveats described by
// https://protobuf.dev/programming-guides/encoding/#implications.
//
// At time of writing, the Go implementation provides deterministic output for
// the same input within the same binary/process, which is sufficient for the
// purposes of this type.
type Proto[K proto.Message, V any] struct {
	elements Map[string, V]
}

// NewProto returns a [Proto] containing the given key/value pairs.
func NewProto[K proto.Message, V any](pairs ...Pair[K, V]) Proto[K, V] {
	var m Proto[K, V]
	for _, p := range pairs {
		m.Set(p.Key, p.Value)
	}
	return m
}

// Set sets the value associated with the given key.
func (m *Proto[K, V]) Set(k K, v V) {
	m.elements.Set(m.marshal(k), v)
}

// Remove removes the given keys from the map.
func (m *Proto[K, V]) Remove(keys ...K) {
	for _, k := range keys {
		m.elements.Remove(m.marshal(k))
	}
}

// Clear removes all keys from the map.
func (m *Proto[K, V]) Clear() {
	m.elements.Clear()
}

// Len returns the number of elements in the map.
func (m Proto[K, V]) Len() int {
	return m.elements.Len()
}

// Has returns true if all of the given keys are in the map.
func (m Proto[K, V]) Has(keys ...K) bool {
	for _, k := range keys {
		if !m.elements.Has(m.marshal(k)) {
			return false
		}
	}
	return true
}

// Get returns the value associated with the given key. It returns the zero
// value if the key is not in the map.
func (m Proto[K, V]) Get(k K) V {
	v, _ := m.TryGet(k)
	return v
}

// TryGet returns the value associated with the given key, or false if the key
// is not in the map.
func (m Proto[K, V]) TryGet(k K) (V, bool) {
	v, ok := m.elements.TryGet(m.marshal(k))
	return v, ok
}

// Clone returns a shallow copy of the map.
func (m Proto[K, V]) Clone() Proto[K, V] {
	return Proto[K, V]{m.elements.Clone()}
}

// Merge returns a new map containing all key/value pairs from s and x.
//
// If a key is present in both maps, the value from x is used.
func (m Proto[K, V]) Merge(x Proto[K, V]) Proto[K, V] {
	return Proto[K, V]{m.elements.Merge(x.elements)}
}

// Select returns a new map containing all key/value pairs from m for which the
// given predicate returns true.
func (m Proto[K, V]) Select(pred func(K, V) bool) Proto[K, V] {
	return Proto[K, V]{
		m.elements.Select(
			func(s string, v V) bool {
				return pred(m.unmarshal(s), v)
			},
		),
	}
}

// Project constructs a new map by applying the given transform function to each
// key/value pair in the map. If the transform function returns false, the key
// is omitted from the resulting map.
func (m Proto[K, V]) Project(transform func(K, V) (K, V, bool)) Proto[K, V] {
	return Proto[K, V]{
		m.elements.Project(
			func(k string, v V) (string, V, bool) {
				if k, v, ok := transform(m.unmarshal(k), v); ok {
					return m.marshal(k), v, true
				}
				return k, v, false
			},
		),
	}
}

// All returns an iterator that yields all key/value pairs in the map in no
// particular order.
func (m Proto[K, V]) All() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for k, v := range m.elements {
			if !yield(m.unmarshal(k), v) {
				return
			}
		}
	}
}

// Keys returns an iterator that yields all keys in the map in no particular
// order.
func (m Proto[K, V]) Keys() iter.Seq[K] {
	return func(yield func(K) bool) {
		for k := range m.elements {
			if !yield(m.unmarshal(k)) {
				return
			}
		}
	}
}

// Values returns an iterator that yields all values in the map in no particular
// order.
func (m Proto[K, V]) Values() iter.Seq[V] {
	return m.elements.Values()
}

func (m Proto[K, V]) marshal(k K) string {
	data, err := proto.
		MarshalOptions{Deterministic: true}.
		Marshal(k)

	if err != nil {
		panic(err)
	}

	return string(data)
}

func (m Proto[K, V]) unmarshal(data string) K {
	t := reflect.TypeFor[K]().Elem()
	k := reflect.New(t).Interface().(K)

	if err := proto.Unmarshal([]byte(data), k); err != nil {
		panic(err)
	}

	return k
}
