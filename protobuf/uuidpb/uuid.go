package uuidpb

import (
	"crypto/rand"
	"crypto/sha1"
	"errors"
	"fmt"
)

// Generate returns a new randomly generated (v4) UUID.
func Generate() *UUID {
	var data [16]byte
	if _, err := rand.Read(data[:]); err != nil {
		panic(err)
	}

	data[6] = (data[6] & 0x0f) | 0x40 // Version 4
	data[8] = (data[8] & 0x3f) | 0x80 // Variant is 10 (RFC 9562)

	return FromByteArray(data)
}

// Derive returns a new SHA-1 derived (v5) UUID from the given namespace and
// names.
func Derive[T ~string | ~[]byte](ns *UUID, names ...T) *UUID {
	hash := sha1.New()
	var data [20]byte // storage for derived UUID and intermediate hash

	CopyBytes(ns, data[:16])

	for _, name := range names {
		hash.Reset()
		hash.Write(data[:16])
		hash.Write([]byte(name))
		hash.Sum(data[:0])

		data[6] = (data[6] & 0x0f) | 0x50 // Version 5
		data[8] = (data[8] & 0x3f) | 0x80 // Variant is 10 (RFC 9562)
	}

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

// Parse parses an RFC 9562 "hex-and-dash" UUID string.
func Parse(str string) (*UUID, error) {
	if len(str) != 36 {
		return nil, errors.New("invalid UUID format, expected 36 characters")
	}

	uuid := &UUID{}
	target := &uuid.Upper
	shift := 60

	for index := range 36 {
		switch index {
		case 18:
			target = &uuid.Lower
			shift = 60
			fallthrough
		case 8, 13, 23:
			if str[index] != '-' {
				return nil, fmt.Errorf("invalid UUID format, expected hyphen at position %d", index)
			}
		default:
			value, err := fromHex(str, index)
			if err != nil {
				return nil, err
			}

			*target |= uint64(value) << shift
			shift -= 4
		}
	}

	return uuid, nil
}

// MustParse parses an RFC 9562 "hex-and-dash" UUID string, or panics if unable
// to do so.
func MustParse(str string) *UUID {
	uuid, err := Parse(str)
	if err != nil {
		panic(err)
	}
	return uuid
}

// ParseAsByteArray parses an RFC 9562 "hex-and-dash" UUID string directly to
// its byte array representation.
//
// This is equivalent to calling AsByteArray() on the result of Parse(str),
// but avoids all allocations.
func ParseAsByteArray(str string) ([16]byte, error) {
	var uuid [16]byte
	return uuid, ParseIntoBytes(str, uuid[:])
}

// MustParseAsByteArray parses an RFC 9562 "hex-and-dash" UUID string directly
// to its byte array representation, or panics if unable to do so.
//
// This is equivalent to calling MustParse(str).AsByteArray(), but avoids
// all allocations.
func MustParseAsByteArray(str string) [16]byte {
	uuid, err := ParseAsByteArray(str)
	if err != nil {
		panic(err)
	}
	return uuid
}

// ParseIntoBytes parses an RFC 9562 "hex-and-dash" UUID string into the given
// byte slice.
//
// This is equivalent to calling [CopyBytes] on the result of [Parse], but
// avoids all allocations.
func ParseIntoBytes(str string, dst []byte) error {
	if len(dst) < 16 {
		return fmt.Errorf(
			"destination slice must have at least 16 bytes, got %d",
			len(dst),
		)
	}

	if len(str) != 36 {
		return errors.New("invalid UUID format, expected 36 characters")
	}

	read := 0
	write := 0

	for read < 36 {
		switch read {
		case 8, 13, 18, 23:
			if str[read] != '-' {
				return fmt.Errorf("invalid UUID format, expected hyphen at position %d", read)
			}

			read++
		}

		high, err := fromHex(str, read)
		if err != nil {
			return err
		}

		read++

		low, err := fromHex(str, read)
		if err != nil {
			return err
		}

		read++

		dst[write] = (high << 4) | low
		write++
	}

	return nil
}

// MustParseIntoBytes parses an RFC 9562 "hex-and-dash" UUID string into the
// given byte slice, or panics if unable to do so.
//
// This is equivalent to calling [CopyBytes] on the result of [MustParse], but
// avoids all allocations.
func MustParseIntoBytes(str string, dst []byte) {
	if err := ParseIntoBytes(str, dst); err != nil {
		panic(err)
	}
}

// ParseAsBytes parses an RFC 9562 "hex-and-dash" UUID string directly to its
// byte slice representation.
//
// This is equivalent to calling AsBytes() on the result of Parse(str),
// but avoids allocation of the intermediate [UUID] struct.
func ParseAsBytes(str string) ([]byte, error) {
	uuid, err := ParseAsByteArray(str)
	if err != nil {
		return nil, err
	}
	return uuid[:], nil
}

// MustParseAsBytes parses an RFC 9562 "hex-and-dash" UUID string directly to
// its byte slice representation, or panics if unable to do so.
//
// This is equivalent to calling MustParse(str).AsBytes(), but avoids allocation
// of the intermediate [UUID] struct.
func MustParseAsBytes(str string) []byte {
	uuid := MustParseAsByteArray(str)
	return uuid[:]
}

// FromBytes returns a UUID from a byte slice.
//
// It returns an error if the length of the slice is not exactly 16 bytes.
func FromBytes[T ~[]B, B byte](data T) (*UUID, error) {
	if len(data) != 16 {
		return nil, fmt.Errorf("slice must be exactly 16 bytes, got %d", len(data))
	}

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
	}, nil
}

// FromByteArray returns a UUID from a byte array.
func FromByteArray[T ~[16]B, B byte](data T) *UUID {
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
func AsByteArray[T ~[16]B, B ~byte](x *UUID) T {
	var data T

	if x != nil {
		data[0] = B(x.Upper >> 56)
		data[1] = B(x.Upper >> 48)
		data[2] = B(x.Upper >> 40)
		data[3] = B(x.Upper >> 32)
		data[4] = B(x.Upper >> 24)
		data[5] = B(x.Upper >> 16)
		data[6] = B(x.Upper >> 8)
		data[7] = B(x.Upper)

		data[8] = B(x.Lower >> 56)
		data[9] = B(x.Lower >> 48)
		data[10] = B(x.Lower >> 40)
		data[11] = B(x.Lower >> 32)
		data[12] = B(x.Lower >> 24)
		data[13] = B(x.Lower >> 16)
		data[14] = B(x.Lower >> 8)
		data[15] = B(x.Lower)
	}

	return data
}

// CopyBytes copies the binary representation of x into dst.
//
// It returns the number of elements copied, which will be the minimum of 16 and
// len(dst).
func CopyBytes[T ~[]B, B byte](x *UUID, dst T) int {
	data := AsByteArray[[16]B](x)
	return copy(dst, data[:])
}

// AsBytes returns the UUID as a byte slice.
func (x *UUID) AsBytes() []byte {
	data := AsByteArray[[16]byte](x)
	return data[:]
}

// AsByteArray returns the UUID as a byte array.
func (x *UUID) AsByteArray() [16]byte {
	return AsByteArray[[16]byte](x)
}

// AsString returns the UUID as an RFC 9562 string.
func (x *UUID) AsString() string {
	return asString(x.GetUpper(), x.GetLower())
}

// DapperString implements [github.com/dogmatiq/dapper.Stringer].
func (x *UUID) DapperString() string {
	return x.AsString()
}

// Format implements the [fmt.Formatter] interface, allowing UUIDs to be
// formatted with functions from the [fmt] package.
//
// This method takes precedence over the [UUID.String] method, which is
// generated by the Protocol Buffers compiler.
func (x *UUID) Format(f fmt.State, verb rune) {
	format := fmt.FormatString(f, verb)

	// If we're formatting as a string, use the RFC 9562 format.
	if verb == 's' || verb == 'q' {
		fmt.Fprintf(f, format, x.AsString())
		return
	}

	// If we're formatting the Go syntax, output something more useful than the
	// protobuf internals.
	if verb == 'v' && f.Flag('#') {
		fmt.Fprintf(f, "uuidpb.MustParse(%q)", x.AsString())
		return
	}

	// Otherwise, fall-back to the default behavior. In order to avoid infinite
	// recursion into this method, we define a new type that does not have any
	// methods.

	// First, we create an alias to the _real_ type so that we can base our new
	// type on it without causing a recursive type definition.
	type realType = UUID

	// Then, we create a new type with the structure of the real type, but none
	// of its methods. We use the same name as the real type so that any format
	// verbs that include the type name (such as "%T") will still print the
	// correct name.
	type UUID realType

	fmt.Fprintf(f, format, (*UUID)(x))
}

// Validate returns an error if x is not a valid Version 4 (random) or Version 5
// (name-based) UUID.
func (x *UUID) Validate() error {
	switch (x.GetUpper() >> 8) & 0xf0 {
	case 0x40, 0x50:
	default:
		return errors.New("UUID must use version 4 or 5")
	}

	if variant := (x.GetLower() >> 56) & 0xc0; variant != 0x80 {
		return fmt.Errorf("UUID must use RFC 9562 variant")
	}

	return nil
}

// Equal returns true if x and id are equal.
func (x *UUID) Equal(id *UUID) bool {
	return x.GetUpper() == id.GetUpper() && x.GetLower() == id.GetLower()
}

// Less returns true if x is less than id.
func (x *UUID) Less(id *UUID) bool {
	return x.Compare(id) < 0
}

// Compare performs a three-way comparison between x and id.
func (x *UUID) Compare(id *UUID) int {
	return compare(
		x.GetUpper(),
		x.GetLower(),
		id.GetUpper(),
		id.GetLower(),
	)
}

// bad is a sentinel value that indicates an invalid hexadecimal digit.
const bad = 255

var (
	toHex = [16]byte{
		'0', '1', '2', '3',
		'4', '5', '6', '7',
		'8', '9', 'a', 'b',
		'c', 'd', 'e', 'f',
	}

	fromHexMap = [256]byte{
		bad, bad, bad, bad, bad, bad, bad, bad, bad, bad, bad, bad, bad, bad, bad, bad,
		bad, bad, bad, bad, bad, bad, bad, bad, bad, bad, bad, bad, bad, bad, bad, bad,
		bad, bad, bad, bad, bad, bad, bad, bad, bad, bad, bad, bad, bad, bad, bad, bad,
		0x0, 0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0x8, 0x9, bad, bad, bad, bad, bad, bad,
		bad, 0xa, 0xb, 0xc, 0xd, 0xe, 0xf, bad, bad, bad, bad, bad, bad, bad, bad, bad,
		bad, bad, bad, bad, bad, bad, bad, bad, bad, bad, bad, bad, bad, bad, bad, bad,
		bad, 0xa, 0xb, 0xc, 0xd, 0xe, 0xf, bad, bad, bad, bad, bad, bad, bad, bad, bad,
		bad, bad, bad, bad, bad, bad, bad, bad, bad, bad, bad, bad, bad, bad, bad, bad,
		bad, bad, bad, bad, bad, bad, bad, bad, bad, bad, bad, bad, bad, bad, bad, bad,
		bad, bad, bad, bad, bad, bad, bad, bad, bad, bad, bad, bad, bad, bad, bad, bad,
		bad, bad, bad, bad, bad, bad, bad, bad, bad, bad, bad, bad, bad, bad, bad, bad,
		bad, bad, bad, bad, bad, bad, bad, bad, bad, bad, bad, bad, bad, bad, bad, bad,
		bad, bad, bad, bad, bad, bad, bad, bad, bad, bad, bad, bad, bad, bad, bad, bad,
		bad, bad, bad, bad, bad, bad, bad, bad, bad, bad, bad, bad, bad, bad, bad, bad,
		bad, bad, bad, bad, bad, bad, bad, bad, bad, bad, bad, bad, bad, bad, bad, bad,
		bad, bad, bad, bad, bad, bad, bad, bad, bad, bad, bad, bad, bad, bad, bad, bad,
	}
)

// fromHex returns the numerical value of a hexadecimal digit at the given index
// in str.
func fromHex(str string, index int) (byte, error) {
	v := fromHexMap[str[index]]
	if v == bad {
		return 0, fmt.Errorf("invalid UUID format, expected hex digit at position %d", index)
	}
	return v, nil
}

// plain is a [comparable] representation of a [UUID].
type plain struct {
	upper, lower uint64
}

// DapperString implements [github.com/dogmatiq/dapper.Stringer].
func (k plain) DapperString() string {
	return asString(k.upper, k.lower)
}

// asString returns the string representation of a UUID with the given upper and
// lower components.
func asString(upper, lower uint64) string {
	var str [36]byte

	source := &upper
	shift := 60

	for index := 0; index < 36; index++ {
		switch index {
		case 18:
			source = &lower
			shift = 60
			fallthrough
		case 8, 13, 23:
			str[index] = '-'
		default:
			value := (*source >> shift) & 0xf
			str[index] = toHex[value]
			shift -= 4
		}
	}

	return string(str[:])
}

func compare(
	upperA, lowerA uint64,
	upperB, lowerB uint64,
) int {
	if upperA > upperB {
		return 1
	} else if upperA < upperB {
		return -1
	} else if lowerA > lowerB {
		return 1
	} else if lowerA < lowerB {
		return -1
	}
	return 0
}
