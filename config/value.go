package config

import (
	"reflect"

	"github.com/dogmatiq/enginekit/internal/typename"
	"github.com/dogmatiq/enginekit/optional"
)

// Value contains information about an arbitrary value of type T that is used
// within the configuration but is not itself a [Component].
//
// T may be an interface.
type Value[T any] struct {
	// TypeName is the fully-qualified name of the concrete Go type of the
	// value. It must refer to an implementation of T or T itself.
	TypeName optional.Optional[string]

	// Value is the actual value, if available.
	Value optional.Optional[T]
}

func normalizeValue[T any](
	ctx *normalizeContext,
	f Fidelity,
	v optional.Optional[Value[T]],
) (Fidelity, optional.Optional[Value[T]]) {
	inner, ok := v.TryGet()
	if !ok {
		return f, v
	}

	typeName, hasTypeName := inner.TypeName.TryGet()
	value, hasValue := inner.Value.TryGet()

	if !hasTypeName {
		f.IsPartial = true
	}

	if hasValue {
		actualTypeName := typename.Get(reflect.TypeOf(value))
		if hasTypeName && typeName != actualTypeName {
			ctx.Fail(TypeNameMismatchError{actualTypeName, typeName})
		}

		inner.TypeName = optional.Some(actualTypeName)
	} else if ctx.Options.RequireImplementations {
		ctx.Fail(ImplementationUnavailableError{reflect.TypeFor[T]()})
	}

	return f, optional.Some(inner)
}
