package config

import (
	"reflect"
	"sync"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/marshalkit"
)

// MessageTypeContainer is an interface for containers of message types.
type MessageTypeContainer interface {
	// Has returns true if t is in the container.
	Has(t MessageType) bool

	// HasM returns true if TypeOf(m) is in the container.
	HasM(m dogma.Message) bool

	// Each invokes fn once for each type in the container.
	//
	// Iteration stops when fn returns false or once fn has been invoked for all
	// types in the container.
	//
	// It returns true if fn returned true for all types.
	Each(fn func(MessageType) bool) bool
}

// MessageType is a value that identifies the type of a message.
type MessageType interface {
	// ReflectType returns the reflect.Type for this message type.
	ReflectType() reflect.Type

	// String returns a human-readable name for the message type.
	// Note that this representation may not be globally unique.
	String() string
}

// NewMessageType returns the message type for the Go type reprsented by rt.
//
// If rt does not implement dogma.Message then mt is nil, and ok is false.
func NewMessageType(rt reflect.Type) (mt MessageType, ok bool) {
	// The current implementation always returns true, as the dogma.Message
	// interface is empty, and hence all types satisfy it.
	return newMessageType(rt), true
}

// MessageTypeOf returns the message type of m.
func MessageTypeOf(m dogma.Message) MessageType {
	rt := reflect.TypeOf(m)
	return newMessageType(rt)
}

// MarshalMessageType marshals a message type to its portable representation.
func MarshalMessageType(ma *marshalkit.Marshaler, mt MessageType) (string, error) {
	return ma.MarshalType(mt.ReflectType())
}

// UnmarshalMessageType unmarshals a message type from its portable
// representation.
func UnmarshalMessageType(ma *marshalkit.Marshaler, mt string) (MessageType, error) {
	rt, err := ma.UnmarshalType(mt)
	if err != nil {
		return nil, err
	}

	return newMessageType(rt), nil
}

// UnmarshalMessageTypeFromMediaType unmarshals a message type from a MIME
// media-type.
func UnmarshalMessageTypeFromMediaType(ma *marshalkit.Marshaler, mt string) (MessageType, error) {
	rt, err := ma.UnmarshalTypeFromMediaType(mt)
	if err != nil {
		return nil, err
	}

	return newMessageType(rt), nil
}

// newMessageType returns the message type for the Go type reprsented by t.
//
// It is assumed that t implements dogma.Message.
func newMessageType(rt reflect.Type) MessageType {
	// This is a compile time assertion that the dogma.Message interface
	// contains no methods.
	//
	// If this line fails to compile due to missing methods, then the contents
	// of this function should be moved into the exposed NewMessageType()
	// function, which returns a boolean to indicate the type's compatibility
	// with dogma.Message.
	//
	// This function should then be removed, and any callers updated to use
	// NewMessageType() instead.
	var _ dogma.Message = interface{}(nil)

	// try to load first, to avoid building the string if it's already stored
	v, loaded := mtypes.Load(rt)

	if !loaded {
		mt := newmtype(rt)

		// try to store the new mt, but if another goroutine has stored it since, use
		// that so that we get the same pointer value.
		v, loaded = mtypes.LoadOrStore(rt, mt)
		if !loaded {
			// if we stored out mt, create the reverse mapping as well
			rtypes.Store(mt, rt)
		}
	}

	return v.(*mtype)
}

var mtypes, rtypes sync.Map

type mtype string

func newmtype(rt reflect.Type) *mtype {
	var n, p string

	for rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
		p += "*"
	}

	if rt.Name() == "" {
		n = "<" + p + "anonymous>"
	} else {
		n = p + rt.String()
	}

	mt := mtype(n)

	return &mt
}

func (mt *mtype) ReflectType() reflect.Type {
	v, _ := rtypes.Load(mt)
	return v.(reflect.Type)
}

func (mt *mtype) String() string {
	return string(*mt)
}
