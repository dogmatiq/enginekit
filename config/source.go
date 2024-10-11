package config

import (
	"github.com/dogmatiq/enginekit/optional"
)

// Source contains information about the type and value that implements
// a [Component] of type T.
type Source[T any] struct {
	// TypeName is the fully-qualified name of the Go type that implements T.
	TypeName string

	// Interface is the value of type T that produced the configuration, if
	// available.
	Interface optional.Optional[T]
}
