package config

import (
	"uuid"

	"github.com/dogmatiq/protobuf/uuidpb"
)

// UnparsedUUID is a string that is expected to contain an RFC 9562 UUID, but
// that hasn't yet been parsed, so it may not be valid. Identity keys and message
// type IDs use this type.
//
// Dogma applications supply these values as strings, and a configuration may
// be invalid. A valid UUID uses the hyphenated form, such as
// "f81d4fae-7dec-11d0-a765-00a0c91e6bf6", in any mix of upper and lower case.
type UnparsedUUID string

// String returns s exactly as it was supplied, which may not be a valid UUID.
func (s UnparsedUUID) String() string {
	return string(s)
}

// AsStd returns s as a [uuid.UUID].
//
// It panics if s isn't a valid UUID.
func (s UnparsedUUID) AsStd() uuid.UUID {
	return uuidpb.MustParseAsByteArray(string(s))
}

// TryAsStd returns s as a [uuid.UUID], or false if s isn't a valid UUID.
func (s UnparsedUUID) TryAsStd() (uuid.UUID, bool) {
	id, err := uuidpb.ParseAsByteArray(string(s))
	return id, err == nil
}

// AsProto returns s as a [uuidpb.UUID].
//
// It panics if s isn't a valid UUID.
func (s UnparsedUUID) AsProto() *uuidpb.UUID {
	return uuidpb.MustParse(string(s))
}

// TryAsProto returns s as a [uuidpb.UUID], or false if s isn't a valid UUID.
func (s UnparsedUUID) TryAsProto() (*uuidpb.UUID, bool) {
	id, err := uuidpb.Parse(string(s))
	return id, err == nil
}
