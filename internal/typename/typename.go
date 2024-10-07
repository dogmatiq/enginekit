package typename

import "reflect"

// For returns the fully-qualified name of T.
func For[T any]() string {
	return Get(reflect.TypeFor[T]())
}

// Of returns the fully-qualified name of v's type.
func Of(v any) string {
	return Get(reflect.TypeOf(v))
}

// Get returns the fully-qualified name of t.
func Get(t reflect.Type) string {
	if t.Name() != "" {
		return t.PkgPath() + "." + t.Name()
	}

	for t.Kind() == reflect.Ptr {
		return "*" + Get(t.Elem())
	}

	panic("cannot build name of unnamed type")
}
