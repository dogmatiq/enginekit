package identitypb

import (
	"errors"
	"fmt"
	"slices"

	uuidpb "github.com/dogmatiq/enginekit/protobuf/uuidpb"
)

// New returns a new identity with the given name and key.
func New(name string, key *uuidpb.UUID) *Identity {
	x := NewIdentityBuilder().
		WithName(name).
		WithKey(key).
		Build()

	if err := x.Validate(); err != nil {
		panic(x)
	}

	return x
}

// Parse returns a new identity with the given name and key by parsing the key
// as a UUID.
func Parse(name, key string) (*Identity, error) {
	k, err := uuidpb.Parse(key)
	if err != nil {
		return nil, fmt.Errorf("invalid key: %w", err)
	}

	x := NewIdentityBuilder().
		WithName(name).
		WithKey(k).
		Build()

	if err := x.Validate(); err != nil {
		return nil, err
	}

	return x, nil
}

// MustParse returns a new identity with the given name and key by parsing the
// key as a UUID, or panics if unable to do so.
func MustParse(name, key string) *Identity {
	x, err := Parse(name, key)
	if err != nil {
		panic(err)
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

// MarshalText implements the [encoding.TextMarshaler] interface.
//
// The text representation is the UUID key followed by a space and the name,
// e.g. "5195fe85-eb3f-4121-84b0-be72cbc5722f handler-name".
func (x *Identity) MarshalText() ([]byte, error) {
	text, err := x.GetKey().MarshalText()
	if err != nil {
		return nil, err
	}

	text = slices.Grow(text, 1+len(x.xxx_hidden_Name))
	text = append(text, ' ')
	text = append(text, x.xxx_hidden_Name...)

	return text, nil
}

// UnmarshalText implements the [encoding.TextUnmarshaler] interface.
func (x *Identity) UnmarshalText(text []byte) error {
	if len(text) < 38 {
		return errors.New("invalid identity format, expected UUID followed by a space and name")
	}

	if text[36] != ' ' {
		return errors.New("invalid identity format, expected space after UUID")
	}

	var key uuidpb.UUID
	if err := key.UnmarshalText(text[:36]); err != nil {
		return fmt.Errorf("invalid key: %w", err)
	}

	x.SetKey(&key)
	x.SetName(string(text[37:]))

	return x.Validate()
}

// Format implements the fmt.Formatter interface, allowing identities to be
// formatted with functions from the fmt package.
func (x *Identity) Format(f fmt.State, verb rune) {
	format := fmt.FormatString(f, verb)

	// If we're formatting as a string, show the UUID followed by the name. A
	// question mark is used as a placeholder when the name is empty.
	if verb == 's' {
		name := x.GetName()
		if name == "" {
			name = "?"
		}

		fmt.Fprintf(f, "%s %s", x.GetKey(), name)
		return
	}

	// If we're formatting the Go syntax, output something more useful than the
	// protobuf internals.
	if verb == 'v' && f.Flag('#') {
		fmt.Fprintf(
			f,
			"identitypb.New(%#v, %#v)",
			x.GetName(),
			x.GetKey(),
		)
		return
	}

	// Otherwise, fall-back to the default behavior. In order to avoid infinite
	// recursion into this method, we define a new type that does not have any
	// methods.

	// First, we create an alias to the _real_ type so that we can base our new
	// type on it without causing a recursive type definition.
	type realType = Identity

	// Then, we create a new type with the structure of the real type, but none
	// of its methods. We use the same name as the real type so that any format
	// verbs that include the type name (such as "%T") will still print the
	// correct name.
	type Identity realType

	fmt.Fprintf(f, format, (*Identity)(x))
}

// Equal returns true if x and id are equal.
func (x *Identity) Equal(id *Identity) bool {
	return x.GetName() == id.GetName() && x.GetKey().Equal(id.GetKey())
}
