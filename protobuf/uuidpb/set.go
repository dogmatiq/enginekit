package uuidpb

// Set is a set of UUIDs.
type Set map[MapKey]struct{}

// Add adds id to the set.
func (s Set) Add(id *UUID) {
	s[id.AsMapKey()] = struct{}{}
}

// Has returns true if id is in the set.
func (s Set) Has(id *UUID) bool {
	_, ok := s[id.AsMapKey()]
	return ok
}

// Delete removes id from the set.
func (s Set) Delete(id *UUID) {
	delete(s, id.AsMapKey())
}
