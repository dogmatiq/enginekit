package marshaling

import (
	"reflect"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/message"
)

// MarshalMessageType marshals a message type to its portable representation.
func MarshalMessageType(ma *Marshaler, mt message.Type) (string, error) {
	return ma.MarshalType(mt.ReflectType())
}

// MustMarshalMessageType marshals a message type to its portable representation.
// It panics if marshaling fails.
func MustMarshalMessageType(ma *Marshaler, mt message.Type) string {
	s, err := MarshalMessageType(ma, mt)
	if err != nil {
		panic(PanicSentinel{err})
	}

	return s
}

// UnmarshalMessageType unmarshals a message type from its portable
// representation.
func UnmarshalMessageType(ma *Marshaler, n string) (message.Type, error) {
	rt, err := ma.UnmarshalType(n)
	if err != nil {
		return nil, err
	}

	return toMessageType(rt), nil
}

// MustUnmarshalMessageType unmarshals a message type from its portable
// representation. It panics if unmarshaling fails.
func MustUnmarshalMessageType(ma *Marshaler, n string) message.Type {
	mt, err := UnmarshalMessageType(ma, n)
	if err != nil {
		panic(PanicSentinel{err})
	}

	return mt
}

// UnmarshalMessageTypeFromMediaType unmarshals a message type from a MIME
// media-type.
func UnmarshalMessageTypeFromMediaType(ma *Marshaler, mt string) (message.Type, error) {
	rt, err := ma.UnmarshalTypeFromMediaType(mt)
	if err != nil {
		return nil, err
	}

	return toMessageType(rt), nil
}

// MustUnmarshalMessageTypeFromMediaType unmarshals a message type from a MIME
// media-type. It panics if unmarshaling fails.
func MustUnmarshalMessageTypeFromMediaType(ma *Marshaler, mt string) message.Type {
	t, err := UnmarshalMessageTypeFromMediaType(ma, mt)
	if err != nil {
		panic(PanicSentinel{err})
	}

	return t
}

// MarshalMessage returns a binary representation of a message.
func MarshalMessage(ma *Marshaler, m dogma.Message) (Packet, error) {
	return ma.Marshal(m)
}

// MustMarshalMessage returns a binary representation of a message.
// It panics if marshaling fails.
func MustMarshalMessage(ma *Marshaler, m dogma.Message) Packet {
	p, err := ma.Marshal(m)
	if err != nil {
		panic(PanicSentinel{err})
	}

	return p
}

// UnmarshalMessage returns a message from its binary representation.
func UnmarshalMessage(ma *Marshaler, p Packet) (dogma.Message, error) {
	// Note: Unmarshal() returns interface{}, which works at the moment because
	// dogma.Message is also empty.
	//
	// If this fails to compile in the future, a branch needs to be added to
	// return a meaningful error if the unmarshaled value does not implement
	// dogma.Message.
	return ma.Unmarshal(p)
}

// MustUnmarshalMessage returns a message from its binary representation.
// It panics if un marshaling fails.
func MustUnmarshalMessage(ma *Marshaler, p Packet) dogma.Message {
	m, err := UnmarshalMessage(ma, p)
	if err != nil {
		panic(PanicSentinel{err})
	}

	return m
}

// toMessageType converts a reflect.Type to a message.Type.
//
// TODO: Remove this function. Blocked by
// https://github.com/dogmatiq/enginekit/issues/8.
func toMessageType(rt reflect.Type) message.Type {
	return message.TypeOf(
		reflect.Zero(rt).Interface(),
	)
}
