package uuidpb_test

import (
	"fmt"
	"testing"

	. "github.com/dogmatiq/enginekit/protobuf/uuidpb"
	"google.golang.org/protobuf/proto"
)

func TestGenerate(t *testing.T) {
	t.Parallel()

	uuid := Generate()

	if uuid.IsNil() {
		t.Fatal("got 'nil' UUID (all zeroes), want random UUID")
	}

	if uuid.IsOmni() {
		t.Fatal("got 'omni' UUID (all ones), want random UUID")
	}

	expect := uint64(0x40)
	actual := (uuid.Upper >> 8) & 0xf0

	if actual != expect {
		t.Fatalf("got version %d, want %d", actual, expect)
	}

	expect = 0x80
	actual = (uuid.Lower >> 56) & 0xc0

	if actual != expect {
		t.Fatalf("got variant %d, want %d (RFC 4122)", actual, expect)
	}
}

func TestFromByteArray(t *testing.T) {
	t.Parallel()

	data := [16]byte{
		0xa9, 0x67, 0xa8, 0xb9,
		0x3f, 0x9c, 0x49, 0x18,
		0x9a, 0x41, 0x19, 0x57,
		0x7b, 0xe5, 0xfe, 0xc5,
	}

	actual := FromByteArray(data)
	expect := &UUID{
		Upper: 0xa967a8b93f9c4918,
		Lower: 0x9a4119577be5fec5,
	}

	if !proto.Equal(actual, expect) {
		t.Fatalf("got %s, want %s", actual, expect)
	}
}

func TestToByteArray(t *testing.T) {
	t.Parallel()

	uuid := &UUID{
		Upper: 0xa967a8b93f9c4918,
		Lower: 0x9a4119577be5fec5,
	}

	actual := ToByteArray[byte](uuid)
	expect := [16]byte{
		0xa9, 0x67, 0xa8, 0xb9,
		0x3f, 0x9c, 0x49, 0x18,
		0x9a, 0x41, 0x19, 0x57,
		0x7b, 0xe5, 0xfe, 0xc5,
	}

	if actual != expect {
		t.Fatalf("got %#v, want %#v", actual, expect)
	}
}

func TestUUID_ToString(t *testing.T) {
	t.Parallel()

	uuid := &UUID{
		Upper: 0xa967a8b93f9c4918,
		Lower: 0x9a4119577be5fec5,
	}

	expect := "a967a8b9-3f9c-4918-9a41-19577be5fec5"
	actual := uuid.ToString()

	if actual != expect {
		t.Fatalf("got %q, want %q", actual, expect)
	}
}

func TestUUID_Format(t *testing.T) {
	t.Parallel()

	uuid := &UUID{
		Upper: 0xa967a8b93f9c4918,
		Lower: 0x9a4119577be5fec5,
	}

	expect := uuid.ToString()
	actual := fmt.Sprintf("%s", uuid)

	if actual != expect {
		t.Fatalf("got %q, want %q", actual, expect)
	}
}

func TestUUID_IsNil(t *testing.T) {
	t.Parallel()

	cases := []struct {
		Description string
		UUID        *UUID
		Expect      bool
	}{
		{"nil pointer", nil, true},
		{"zero value", &UUID{}, true},
		{"non-zero upper component", &UUID{Upper: 1}, false},
		{"non-zero lower component", &UUID{Lower: 1}, false},
		{"result of Nil()", Nil(), true},
	}

	for _, c := range cases {
		c := c // capture loop variable
		t.Run(c.Description, func(t *testing.T) {
			t.Parallel()

			actual := c.UUID.IsNil()
			if actual != c.Expect {
				t.Fatalf("got %t, want %t", actual, c.Expect)
			}
		})
	}
}

func TestUUID_IsOmni(t *testing.T) {
	t.Parallel()

	cases := []struct {
		Description string
		UUID        *UUID
		Expect      bool
	}{
		{"nil pointer", nil, false},
		{"zero value", &UUID{}, false},
		{"all ones in upper component", &UUID{Upper: 0xffffffffffffffff}, false},
		{"all ones in lower component", &UUID{Lower: 0xffffffffffffffff}, false},
		{"all ones", &UUID{Upper: 0xffffffffffffffff, Lower: 0xffffffffffffffff}, true},
		{"result of Omni()", Omni(), true},
	}

	for _, c := range cases {
		c := c // capture loop variable
		t.Run(c.Description, func(t *testing.T) {
			t.Parallel()

			actual := c.UUID.IsOmni()
			if actual != c.Expect {
				t.Fatalf("got %t, want %t", actual, c.Expect)
			}
		})
	}
}
