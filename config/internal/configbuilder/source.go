package configbuilder

import (
	"github.com/dogmatiq/enginekit/config"
	"github.com/dogmatiq/enginekit/internal/typename"
	"github.com/dogmatiq/enginekit/optional"
)

func setSourceTypeName[T any](
	src *config.Value[T],
	typeName string,
) {
	if typeName == "" {
		panic("type name must not be empty")
	}

	*src = config.Value[T]{
		TypeName: optional.Some(typeName),
	}
}

func setSource[T any](
	src *config.Value[T],
	value T,
) {
	if any(value) == nil {
		panic("value must not be nil")
	}

	*src = config.Value[T]{
		TypeName: optional.Some(typename.Of(value)),
		Value:    optional.Some(value),
	}
}
