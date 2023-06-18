package uuidpb

// Map is a map data structure that uses a UUID as its key.
type Map[V any] map[MapKey]V

// Set associates v with k.
func (m Map[V]) Set(k *UUID, v V) {
	m[asMapKey(k)] = v
}

// Get returns the value associated with k, or the zero-value if k is not
// present in the map.
func (m Map[V]) Get(k *UUID) V {
	return m[asMapKey(k)]
}

// TryGet returns the value associated with k, or false if k is not present in
// the map.
func (m Map[V]) TryGet(k *UUID) (V, bool) {
	v, ok := m[asMapKey(k)]
	return v, ok
}

// Delete removes the value associated with k.
func (m Map[V]) Delete(k *UUID) {
	delete(m, asMapKey(k))
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

// asMapKey returns an opaque representation of the UUID that can be used as a
// map key.
func asMapKey(x *UUID) MapKey {
	return MapKey{x.GetUpper(), x.GetLower()}
}
