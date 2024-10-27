package config

import (
	"github.com/dogmatiq/enginekit/optional"
)

// Identity is a [Component] that that represents the unique identity of an
// [Entity].
type Identity struct {
	ComponentCommon

	Name optional.Optional[string]
	Key  optional.Optional[string]
}
