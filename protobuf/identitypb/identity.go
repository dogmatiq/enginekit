package identitypb

import (
	"errors"
	"fmt"
)

// Validate returns an error if x is invalid.
//
// It does not perform UTF-8 validation on the name. This should be validated by
// the engine when the identity is configured.
func (x *Identity) Validate() error {
	name := x.GetName()

	if len(name) == 0 || len(name) > 255 {
		return errors.New("invalid identity name: must be between 1 and 255 bytes")
	}

	if x.GetKey().IsNil() {
		return fmt.Errorf("invalid identity key: must not be the nil UUID (all zeroes)")
	}

	if x.GetKey().IsOmni() {
		return fmt.Errorf("invalid identity key: must not be the omni UUID (all ones)")
	}

	return nil
}
