package uuidpb

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/dogmatiq/enginekit/internal/fmtbackport"
)

// Generate returns a new randonly generated UUID.
func Generate() *UUID {
	var data [16]byte
	if _, err := rand.Read(data[:]); err != nil {
		panic(err)
	}

	data[6] = (data[6] & 0x0f) | 0x40 // Version 4
	data[8] = (data[8] & 0x3f) | 0x80 // Variant is 10 (RFC 4122)

	return FromByteArray(data)
}

// FromByteArray returns a UUID from a byte array.
func FromByteArray[T ~byte](data [16]T) *UUID {
	return &UUID{
		Upper: uint64(data[0])<<56 |
			uint64(data[1])<<48 |
			uint64(data[2])<<40 |
			uint64(data[3])<<32 |
			uint64(data[4])<<24 |
			uint64(data[5])<<16 |
			uint64(data[6])<<8 |
			uint64(data[7]),
		Lower: uint64(data[8])<<56 |
			uint64(data[9])<<48 |
			uint64(data[10])<<40 |
			uint64(data[11])<<32 |
			uint64(data[12])<<24 |
			uint64(data[13])<<16 |
			uint64(data[14])<<8 |
			uint64(data[15]),
	}
}

// AsByteArray returns the UUID as a byte array.
func AsByteArray[T ~byte](x *UUID) [16]T {
	var data [16]T

	if x != nil {
		data[0] = T(x.Upper >> 56)
		data[1] = T(x.Upper >> 48)
		data[2] = T(x.Upper >> 40)
		data[3] = T(x.Upper >> 32)
		data[4] = T(x.Upper >> 24)
		data[5] = T(x.Upper >> 16)
		data[6] = T(x.Upper >> 8)
		data[7] = T(x.Upper)

		data[8] = T(x.Lower >> 56)
		data[9] = T(x.Lower >> 48)
		data[10] = T(x.Lower >> 40)
		data[11] = T(x.Lower >> 32)
		data[12] = T(x.Lower >> 24)
		data[13] = T(x.Lower >> 16)
		data[14] = T(x.Lower >> 8)
		data[15] = T(x.Lower)
	}

	return data
}

// AsBytes returns the UUID as a byte slice.
func (x *UUID) AsBytes() []byte {
	data := AsByteArray[byte](x)
	return data[:]
}

// AsString returns the UUID as an RFC 4122 string.
func (x *UUID) AsString() string {
	var str [36]byte
	uuid := AsByteArray[byte](x)

	hex.Encode(str[:], uuid[:4])
	str[8] = '-'
	hex.Encode(str[9:], uuid[4:6])
	str[13] = '-'
	hex.Encode(str[14:], uuid[6:8])
	str[18] = '-'
	hex.Encode(str[19:], uuid[8:10])
	str[23] = '-'
	hex.Encode(str[24:], uuid[10:])

	return string(str[:])
}

// Format implements the fmt.Formatter interface, allowing UUIDs to be formatted
// with functions from the fmt package.
func (x *UUID) Format(f fmt.State, verb rune) {
	fmt.Fprintf(
		f,
		fmtbackport.FormatString(f, verb),
		x.AsString(),
	)
}

// Validate returns an error if x is not a valid Version 4 (random) UUID.
func (x *UUID) Validate() error {
	if version := (x.GetUpper() >> 8) & 0xf0; version != 0x40 {
		return errors.New("UUID must use version 4")
	}

	if variant := (x.GetLower() >> 56) & 0xc0; variant != 0x80 {
		return fmt.Errorf("UUID must use RFC 4122 variant")
	}

	return nil
}
