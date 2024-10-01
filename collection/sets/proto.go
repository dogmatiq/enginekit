package sets

import (
	"iter"
	"reflect"

	"google.golang.org/protobuf/proto"
)

// Proto is an unordered set of unique Protocol Buffers messages of type T.
//
// T must be a pointer type that implements [proto.Message].
//
// Equality is determined based on the serialized form of the message, and so is
// subject to the caveats described by
// https://protobuf.dev/programming-guides/encoding/#implications.
//
// At time of writing, the Go implementation provides deterministic output
// for the same input within the same binary/process, which is sufficient for
// the purposes of this type.
type Proto[T proto.Message] struct {
	members Set[string]
}

// NewProto returns a [Proto] containing the given members.
func NewProto[T proto.Message](members ...T) Proto[T] {
	var s Proto[T]
	s.Add(members...)
	return s
}

// Add adds the given members to the set.
func (s *Proto[T]) Add(members ...T) {
	for _, m := range members {
		s.members.Add(s.marshal(m))
	}
}

// Remove removes the given members from the set.
func (s *Proto[T]) Remove(members ...T) {
	for _, m := range members {
		s.members.Remove(s.marshal(m))
	}
}

// Clear removes all members from the set.
func (s *Proto[T]) Clear() {
	s.members.Clear()
}

// Len returns the number of members in the set.
func (s Proto[T]) Len() int {
	return s.members.Len()
}

// Has returns true if all of the given values are members of the set.
func (s Proto[T]) Has(members ...T) bool {
	for _, m := range members {
		if !s.members.Has(s.marshal(m)) {
			return false
		}
	}
	return true
}

// IsEqual returns true if s and x have the same members.
func (s Proto[T]) IsEqual(x Proto[T]) bool {
	return s.members.IsEqual(x.members)
}

// IsSuperset returns true if s has all of the members of x.
func (s Proto[T]) IsSuperset(x Proto[T]) bool {
	return s.members.IsSuperset(x.members)
}

// IsSubset returns true if x has all of the members of s.
func (s Proto[T]) IsSubset(x Proto[T]) bool {
	return s.members.IsSubset(x.members)
}

// IsStrictSuperset returns true if s has all of the members of x and at least
// one member that is not in x.
func (s Proto[T]) IsStrictSuperset(x Proto[T]) bool {
	return s.members.IsStrictSuperset(x.members)
}

// IsStrictSubset returns true if x has all of the members of s and at least one
// member that is not in s.
func (s Proto[T]) IsStrictSubset(x Proto[T]) bool {
	return s.members.IsStrictSubset(x.members)
}

// Clone returns a shallow copy of the set.
func (s Proto[T]) Clone() Proto[T] {
	return Proto[T]{s.members.Clone()}
}

// Union returns a set containing all members of s and x.
func (s Proto[T]) Union(x Proto[T]) Proto[T] {
	return Proto[T]{s.members.Union(x.members)}
}

// Select returns the subset of s containing members for which the given
// predicate function returns true.
func (s Proto[T]) Select(pred func(T) bool) Proto[T] {
	return Proto[T]{
		members: s.members.Select(func(m string) bool {
			return pred(s.unmarshal(m))
		}),
	}
}

// All returns an iterator that yields all members of the set in no particular
// order.
func (s Proto[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		for m := range s.members.All() {
			if !yield(s.unmarshal(m)) {
				return
			}
		}
	}
}

func (s Proto[T]) marshal(m T) string {
	data, err := proto.
		MarshalOptions{Deterministic: true}.
		Marshal(m)

	if err != nil {
		panic(err)
	}

	return string(data)
}

func (s Proto[T]) unmarshal(data string) T {
	t := reflect.TypeFor[T]().Elem()
	m := reflect.New(t).Interface().(T)

	if err := proto.Unmarshal([]byte(data), m); err != nil {
		panic(err)
	}

	return m
}
