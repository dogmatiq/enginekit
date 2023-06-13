package uuidpb

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math"
	"strconv"
	"unicode/utf8"
)

// Generate returns a new randonly generated UUID.
func Generate() *UUID {
	var data [16]byte
	_, err := rand.Read(data[:])
	if err != nil {
		panic(err)
	}

	data[6] = (data[6] & 0x0f) | 0x40 // Version 4
	data[8] = (data[8] & 0x3f) | 0x80 // Variant is 10

	return &UUID{
		Upper: binary.BigEndian.Uint64(data[:8]),
		Lower: binary.BigEndian.Uint64(data[8:]),
	}
}

// Nil returns the "nil UUID", that is, a UUID with all bits set to zero.
//
// This should not be confused with Go's nil value.
func Nil() *UUID {
	return &UUID{}
}

// Omni returns the "omni UUID", that is, a UUID with all bits set to one.
func Omni() *UUID {
	return &UUID{
		Upper: math.MaxUint64,
		Lower: math.MaxUint64,
	}
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

// ToByteArray returns the UUID as a byte array.
func ToByteArray[T ~byte](x *UUID) [16]T {
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

// ToString returns the UUID as an RFC 4122 string.
func (x *UUID) ToString() string {
	var str [36]byte
	uuid := ToByteArray[byte](x)

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
		formatString(f, verb),
		x.ToString(),
	)
}

// IsNil returns true if x is the "nil UUID".
//
// The nil UUID is a UUID with all bits set to zero.
//
// This should not be confused with Go's nil value, although in this
// implementation a nil *UUID is one possible representation of the nil UUID.
func (x *UUID) IsNil() bool {
	return x.GetUpper() == 0 && x.GetLower() == 0
}

// IsOmni returns true if x is the "omni UUID".
//
// The omni UUID is a UUID with all bits set to one.
func (x *UUID) IsOmni() bool {
	return x.GetUpper() == math.MaxUint64 && x.GetLower() == math.MaxUint64
}

// FormatString returns a string representing the fully qualified formatting
// directive captured by the State, followed by the argument verb.
//
// It is a copy of the fmt.FormatString() function in Go v1.20, and may be
// removed once this Go v1.19 support is dropped.
func formatString(state fmt.State, verb rune) string {
	var tmp [16]byte // Use a local buffer.
	b := append(tmp[:0], '%')
	for _, c := range " +-#0" { // All known flags
		if state.Flag(int(c)) { // The argument is an int for historical reasons.
			b = append(b, byte(c))
		}
	}
	if w, ok := state.Width(); ok {
		b = strconv.AppendInt(b, int64(w), 10)
	}
	if p, ok := state.Precision(); ok {
		b = append(b, '.')
		b = strconv.AppendInt(b, int64(p), 10)
	}
	b = utf8.AppendRune(b, verb)
	return string(b)
}
