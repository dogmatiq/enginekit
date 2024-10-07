package config

import (
	"github.com/dogmatiq/enginekit/optional"
)

// Implementation contains information about the implementation of the T
// interface.
type Implementation[T any] struct {
	// TypeName is the fully-qualified name of the Go type that implements T.
	TypeName string

	// Source is the value that produced the configuration, if available.
	Source optional.Optional[T]
}
