package config

import "github.com/dogmatiq/dogma"

// MessageRoleMap is a map of message type to role.
// It implements the MessageTypeContainer interface.
type MessageRoleMap map[MessageType]MessageRole

// Has returns true if rm contains t.
func (rm MessageRoleMap) Has(t MessageType) bool {
	_, ok := rm[t]
	return ok
}

// HasM returns true if rm contains MessageTypeOf(m).
func (rm MessageRoleMap) HasM(m dogma.Message) bool {
	return rm.Has(MessageTypeOf(m))
}

// Add maps t to r.
//
// It returns true if the mapping was added, or false if the map already
// contained the type.
func (rm MessageRoleMap) Add(t MessageType, r MessageRole) bool {
	if _, ok := rm[t]; ok {
		return false
	}

	rm[t] = r
	return true
}

// AddM adds MessageTypeOf(m) to rm.
//
// It returns true if the mapping was added, or false if the map already
// contained the type.
func (rm MessageRoleMap) AddM(m dogma.Message, r MessageRole) bool {
	return rm.Add(MessageTypeOf(m), r)
}

// Remove removes t from rm.
//
// It returns true if the type was removed, or false if the set did not contain
// the type.
func (rm MessageRoleMap) Remove(t MessageType) bool {
	if _, ok := rm[t]; ok {
		delete(rm, t)
		return true
	}

	return false
}

// RemoveM removes MessageTypeOf(m) from rm.
//
// It returns true if the type was removed, or false if the set did not contain
// the type.
func (rm MessageRoleMap) RemoveM(m dogma.Message) bool {
	return rm.Remove(MessageTypeOf(m))
}

// Each invokes fn once for each type in the container.
//
// Iteration stops when fn returns false or once fn has been invoked for all
// types in the container.
//
// It returns true if fn returned true for all types.
func (rm MessageRoleMap) Each(fn func(MessageType) bool) bool {
	for t := range rm {
		if !fn(t) {
			return false
		}
	}

	return true
}
