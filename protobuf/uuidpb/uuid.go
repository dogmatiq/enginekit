package uuidpb

import (
	"crypto/rand"
	"errors"
	"fmt"
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

// FromString returns a UUID from an RFC 4122 string
func FromString(str string) (*UUID, error) {
	if len(str) != 36 {
		return nil, errors.New("invalid UUID format, expected 36 characters")
	}

	uuid := &UUID{}
	target := &uuid.Upper
	shift := 60

	for index := 0; index < 36; index++ {
		char := str[index]

		switch index {
		case 18:
			target = &uuid.Lower
			shift = 60
			fallthrough
		case 8, 13, 23:
			if char != '-' {
				return nil, errors.New("invalid UUID format, expected hyphen")
			}
		default:
			value := fromHex[char]
			if value == bad {
				return nil, errors.New("invalid UUID format, expected hex digit")
			}

			*target |= uint64(value) << shift
			shift -= 4
		}
	}

	return uuid, nil
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

// AsBytes returns the UUID as a byte slice.
func (x *UUID) AsBytes() []byte {
	data := AsByteArray[[16]byte](x)
	return data[:]
}

// AsString returns the UUID as an RFC 4122 string.
func (x *UUID) AsString() string {
	if x == nil {
		return "00000000-0000-0000-0000-000000000000"
	}

	var str [36]byte

	source := &x.Upper
	shift := 60

	for index := 0; index < 36; index++ {
		switch index {
		case 18:
			source = &x.Lower
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

// DapperString implements [github.com/dogmatiq/dapper.Stringer].
func (x *UUID) DapperString() string {
	return x.AsString()
}

// Format implements the fmt.Formatter interface, allowing UUIDs to be formatted
// with functions from the fmt package.
func (x *UUID) Format(f fmt.State, verb rune) {
	fmt.Fprintf(
		f,
		fmt.FormatString(f, verb),
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
	if x.GetUpper() > id.GetUpper() {
		return 1
	} else if x.GetUpper() < id.GetUpper() {
		return -1
	} else if x.GetLower() > id.GetLower() {
		return 1
	} else if x.GetLower() < id.GetLower() {
		return -1
	}
	return 0
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

	fromHex = [256]byte{
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
