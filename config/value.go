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
	ctx *normalizationContext,
	v *Value[T],
	f *Fidelity,
) {
	typeName, typeNameOK := v.TypeName.TryGet()
	value, valueOK := v.Value.TryGet()

	if !typeNameOK {
		*f |= Incomplete
	}

	if valueOK {
		actualTypeName := typename.Get(reflect.TypeOf(value))
		if typeNameOK && typeName != actualTypeName {
			ctx.Fail(TypeNameMismatchError{actualTypeName, typeName})
		}

		v.TypeName = optional.Some(actualTypeName)
	} else if ctx.Options.RequireValues {
		ctx.Fail(ImplementationUnavailableError{reflect.TypeFor[T]()})
	}
}
