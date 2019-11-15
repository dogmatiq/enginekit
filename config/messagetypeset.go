package config

import (
	"github.com/dogmatiq/dogma"
)

// MessageTypeSet is a collection of distinct message types.
// It implements the MessageTypeContainer interface.
type MessageTypeSet map[MessageType]struct{}

// NewMessageTypeSet returns a MessageTypeSet containing the given types.
func NewMessageTypeSet(types ...MessageType) MessageTypeSet {
	s := MessageTypeSet{}

	for _, t := range types {
		s[t] = struct{}{}
	}

	return s
}

// MessageTypesOf returns a type set containing the types of the given messages.
func MessageTypesOf(messages ...dogma.Message) MessageTypeSet {
	s := MessageTypeSet{}

	for _, m := range messages {
		s[MessageTypeOf(m)] = struct{}{}
	}

	return s
}

// Has returns true if s contains t.
func (s MessageTypeSet) Has(t MessageType) bool {
	_, ok := s[t]
	return ok
}

// HasM returns true if s contains MessageTypeOf(m).
func (s MessageTypeSet) HasM(m dogma.Message) bool {
	return s.Has(MessageTypeOf(m))
}

// Add adds t to s.
//
// It returns true if the type was added, or false if the set already contained
// the type.
func (s MessageTypeSet) Add(t MessageType) bool {
	if _, ok := s[t]; ok {
		return false
	}

	s[t] = struct{}{}
	return true
}

// AddM adds MessageTypeOf(m) to s.
//
// It returns true if the type was added, or false if the set already contained
// the type.
func (s MessageTypeSet) AddM(m dogma.Message) bool {
	return s.Add(MessageTypeOf(m))
}

// Remove removes t from s.
//
// It returns true if the type was removed, or false if the set did not contain
// the type.
func (s MessageTypeSet) Remove(t MessageType) bool {
	if _, ok := s[t]; ok {
		delete(s, t)
		return true
	}

	return false
}

// RemoveM removes MessageTypeOf(m) from s.
//
// It returns true if the type was removed, or false if the set did not contain
// the type.
func (s MessageTypeSet) RemoveM(m dogma.Message) bool {
	return s.Remove(MessageTypeOf(m))
}

// Each invokes fn once for each type in the container.
//
// Iteration stops when fn returns false or once fn has been invoked for all
// types in the container.
//
// It returns true if fn returned true for all types.
func (s MessageTypeSet) Each(fn func(MessageType) bool) bool {
	for t := range s {
		if !fn(t) {
			return false
		}
	}

	return true
}
