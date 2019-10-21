package message

import (
	"github.com/dogmatiq/dogma"
)

// TypeSet is a collection of distinct message types.
// It implements the TypeContainer interface.
type TypeSet map[Type]struct{}

// NewTypeSet returns a TypeSet containing the given types.
func NewTypeSet(types ...Type) TypeSet {
	s := TypeSet{}

	for _, t := range types {
		s[t] = struct{}{}
	}

	return s
}

// TypesOf returns a type set containing the types of the given messages.
func TypesOf(messages ...dogma.Message) TypeSet {
	s := TypeSet{}

	for _, m := range messages {
		s[TypeOf(m)] = struct{}{}
	}

	return s
}

// Has returns true if s contains t.
func (s TypeSet) Has(t Type) bool {
	_, ok := s[t]
	return ok
}

// HasM returns true if s contains TypeOf(m).
func (s TypeSet) HasM(m dogma.Message) bool {
	return s.Has(TypeOf(m))
}

// Add adds t to s.
//
// It returns true if the type was added, or false if the set already contained
// the type.
func (s TypeSet) Add(t Type) bool {
	if _, ok := s[t]; ok {
		return false
	}

	s[t] = struct{}{}
	return true
}

// AddM adds TypeOf(m) to s.
//
// It returns true if the type was added, or false if the set already contained
// the type.
func (s TypeSet) AddM(m dogma.Message) bool {
	return s.Add(TypeOf(m))
}

// Remove removes t from s.
//
// It returns true if the type was removed, or false if the set did not contain
// the type.
func (s TypeSet) Remove(t Type) bool {
	if _, ok := s[t]; ok {
		delete(s, t)
		return true
	}

	return false
}

// RemoveM removes TypeOf(m) from s.
//
// It returns true if the type was removed, or false if the set did not contain
// the type.
func (s TypeSet) RemoveM(m dogma.Message) bool {
	return s.Remove(TypeOf(m))
}

// Each invokes fn once for each type in the container.
//
// Iteration stops when fn returns false or once fn has been invoked for all
// types in the container.
//
// It returns true if fn returned true for all types.
func (s TypeSet) Each(fn func(Type) bool) bool {
	for t := range s {
		if !fn(t) {
			return false
		}
	}

	return true
}
