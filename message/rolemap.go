package message

import (
	"github.com/dogmatiq/dogma"
)

// RoleMap is a map of message type to role.
// It implements the TypeContainer interface.
type RoleMap map[Type]Role

// Has returns true if rm contains t.
func (rm RoleMap) Has(t Type) bool {
	_, ok := rm[t]
	return ok
}

// HasM returns true if rm contains TypeOf(m).
func (rm RoleMap) HasM(m dogma.Message) bool {
	return rm.Has(TypeOf(m))
}

// Add maps t to r.
//
// It returns true if the mapping was added, or false if the map already
// contained the type.
func (rm RoleMap) Add(t Type, r Role) bool {
	if _, ok := rm[t]; ok {
		return false
	}

	rm[t] = r
	return true
}

// AddM adds TypeOf(m) to rm.
//
// It returns true if the mapping was added, or false if the map already
// contained the type.
func (rm RoleMap) AddM(m dogma.Message, r Role) bool {
	return rm.Add(TypeOf(m), r)
}

// Remove removes t from rm.
//
// It returns true if the type was removed, or false if the set did not contain
// the type.
func (rm RoleMap) Remove(t Type) bool {
	if _, ok := rm[t]; ok {
		delete(rm, t)
		return true
	}

	return false
}

// RemoveM removes TypeOf(m) from rm.
//
// It returns true if the type was removed, or false if the set did not contain
// the type.
func (rm RoleMap) RemoveM(m dogma.Message) bool {
	return rm.Remove(TypeOf(m))
}
