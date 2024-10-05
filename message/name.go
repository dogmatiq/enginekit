package message

import (
	"reflect"

	"github.com/dogmatiq/dogma"
)

// Name is the fully-qualified name of a Go type that implements
// [dogma.Command], [dogma.Event] or [dogma.Timeout].
type Name string

// NameFor returns the [Name] of T.
//
// It panics if T does not implement [dogma.Command], [dogma.Event] or
// [dogma.Timeout].
func NameFor[T dogma.Message]() Name {
	return nameFromReflect(reflect.TypeFor[T]())
}

// NameOf returns the fully-qualified type name of m.
//
// It panics if m does not implement [dogma.Command], [dogma.Event] or
// [dogma.Timeout].
func NameOf(m dogma.Message) Name {
	return nameFromReflect(reflect.TypeOf(m))
}

func nameFromReflect(r reflect.Type) Name {
	guardAgainstNonMessage(r)
	return Name(r.PkgPath() + "." + r.Name())
}
