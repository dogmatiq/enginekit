package identitypb

import (
	"errors"
	"fmt"

	"github.com/dogmatiq/enginekit/internal/fmtbackport"
	uuidpb "github.com/dogmatiq/enginekit/protobuf/uuidpb"
)

// New returns a new identity with the given name and key.
func New(name string, key *uuidpb.UUID) *Identity {
	x := &Identity{
		Name: name,
		Key:  key,
	}

	if err := x.Validate(); err != nil {
		panic(x)
	}

	return x
}

// Validate returns an error if x is invalid.
//
// It does not perform UTF-8 validation on the name. This should be validated by
// the engine when the identity is configured.
func (x *Identity) Validate() error {
	name := x.GetName()

	if len(name) == 0 || len(name) > 255 {
		return errors.New("invalid name: must be between 1 and 255 bytes")
	}

	if err := x.GetKey().Validate(); err != nil {
		return fmt.Errorf("invalid key: %w", err)
	}

	return nil
}

// Format implements the fmt.Formatter interface, allowing UUIDs to be formatted
// with functions from the fmt package.
func (x *Identity) Format(f fmt.State, verb rune) {
	fmt.Fprintf(
		f,
		fmtbackport.FormatString(f, verb),
		fmt.Sprintf("%s/%s", x.GetName(), x.GetKey()),
	)
}
