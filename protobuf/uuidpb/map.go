package uuidpb

import "fmt"

// Map is a map data structure that uses a UUID as its key.
type Map[V any] map[MapKey]V

// Set associates v with k.
func (m Map[V]) Set(k *UUID, v V) {
	m[k.AsMapKey()] = v
}

// Get returns the value associated with k, or the zero-value if k is not
// present in the map.
func (m Map[V]) Get(k *UUID) V {
	return m[k.AsMapKey()]
}

// TryGet returns the value associated with k, or false if k is not present in
// the map.
func (m Map[V]) TryGet(k *UUID) (V, bool) {
	v, ok := m[k.AsMapKey()]
	return v, ok
}

// Delete removes the value associated with k.
func (m Map[V]) Delete(k *UUID) {
	delete(m, k.AsMapKey())
}

// MapKey is an opaque representation of a UUID that can be used as a map key.
type MapKey struct {
	upper, lower uint64
}

// AsUUID returns the UUID represented by the key.
func (k MapKey) AsUUID() *UUID {
	return &UUID{
		Upper: k.upper,
		Lower: k.lower,
	}
}

// Format implements the fmt.Formatter interface, allowing UUIDs to be formatted
// with functions from the fmt package.
func (k MapKey) Format(f fmt.State, verb rune) {
	k.AsUUID().Format(f, verb)
}

// String returns a string representation of the UUID.
func (k MapKey) String() string {
	return k.AsUUID().AsString()
}

// DapperString implements [github.com/dogmatiq/dapper.Stringer].
func (k MapKey) DapperString() string {
	return k.String()
}

// AsMapKey returns an opaque representation of the UUID that can be used as a
// map key.
func (x *UUID) AsMapKey() MapKey {
	return MapKey{x.GetUpper(), x.GetLower()}
}
