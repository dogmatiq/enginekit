package message

import (
	"reflect"
	"sync"

	"github.com/dogmatiq/dogma"
)

// messageType is the reflect.Type for the dogma.Message interface.
var messageType = reflect.TypeOf((*dogma.Message)(nil)).Elem()

// TypeContainer is an interface for containers of message types.
type TypeContainer interface {
	// Has returns true if t is in the container.
	Has(t Type) bool

	// HasM returns true if TypeOf(m) is in the container.
	HasM(m dogma.Message) bool

	// Each invokes fn once for each type in the container.
	//
	// Iteration stops when fn returns false or once fn has been invoked for all
	// types in the container.
	//
	// It returns true if fn returned true for all types.
	Each(fn func(Type) bool) bool
}

// Type is a value that identifies the type of a message.
type Type interface {
	// ReflectType returns the reflect.Type for this message type.
	ReflectType() reflect.Type

	// String returns a human-readable name for the message type.
	// Note that this representation may not be globally unique.
	String() string
}

// TypeOf returns the message type of m.
func TypeOf(m dogma.Message) Type {
	rt := reflect.TypeOf(m)
	return fromReflectType(rt)
}

// FromReflectType returns the message type for the Go type reprsented by rt.
//
// If rt does not implement dogma.Message then mt is nil, and ok is false.
func FromReflectType(rt reflect.Type) (mt Type, ok bool) {
	// This is a compile time assertion that the dogma.Message interface
	// contains no methods. If this line fails to compile due to missing
	// methods, additional logic must be added to handle the case where rt
	// does not implement dogma.Message.
	var _ dogma.Message = interface{}(nil)

	// The current implementation always returns true, as the dogma.Message
	// interface is empty, and hence all types satisfy it.
	return fromReflectType(rt), true
}

// fromReflectType returns the message type for the Go type reprsented by t.
//
// It is assumed that t implements dogma.Message.
func fromReflectType(rt reflect.Type) Type {
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
