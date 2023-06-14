package uuidpb_test

import (
	"bytes"
	"fmt"
	"testing"

	. "github.com/dogmatiq/enginekit/protobuf/uuidpb"
	"google.golang.org/protobuf/proto"
)

func TestGenerate(t *testing.T) {
	t.Parallel()

	uuid := Generate()

	if err := uuid.Validate(); err != nil {
		t.Fatal(err)
	}
}

func TestFromByteArray(t *testing.T) {
	t.Parallel()

	subject := [16]byte{
		0xa9, 0x67, 0xa8, 0xb9,
		0x3f, 0x9c, 0x49, 0x18,
		0x9a, 0x41, 0x19, 0x57,
		0x7b, 0xe5, 0xfe, 0xc5,
	}

	actual := FromByteArray(subject)
	expect := &UUID{
		Upper: 0xa967a8b93f9c4918,
		Lower: 0x9a4119577be5fec5,
	}

	if !proto.Equal(actual, expect) {
		t.Fatalf("got %s, want %s", actual, expect)
	}
}

func TestAsByteArray(t *testing.T) {
	t.Parallel()

	subject := &UUID{
		Upper: 0xa967a8b93f9c4918,
		Lower: 0x9a4119577be5fec5,
	}

	actual := AsByteArray[[16]byte](subject)
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

func TestUUID_AsBytes(t *testing.T) {
	t.Parallel()

	subject := &UUID{
		Upper: 0xa967a8b93f9c4918,
		Lower: 0x9a4119577be5fec5,
	}

	actual := subject.AsBytes()
	expect := []byte{
		0xa9, 0x67, 0xa8, 0xb9,
		0x3f, 0x9c, 0x49, 0x18,
		0x9a, 0x41, 0x19, 0x57,
		0x7b, 0xe5, 0xfe, 0xc5,
	}

	if !bytes.Equal(actual, expect) {
		t.Fatalf("got %q, want %q", actual, expect)
	}
}

func TestUUID_AsString(t *testing.T) {
	t.Parallel()

	subject := &UUID{
		Upper: 0xa967a8b93f9c4918,
		Lower: 0x9a4119577be5fec5,
	}

	expect := "a967a8b9-3f9c-4918-9a41-19577be5fec5"
	actual := subject.AsString()

	if actual != expect {
		t.Fatalf("got %q, want %q", actual, expect)
	}
}

func TestUUID_Format(t *testing.T) {
	t.Parallel()

	subject := &UUID{
		Upper: 0xa967a8b93f9c4918,
		Lower: 0x9a4119577be5fec5,
	}

	expect := subject.AsString()
	actual := fmt.Sprintf("%s", subject)

	if actual != expect {
		t.Fatalf("got %q, want %q", actual, expect)
	}
}

func TestUUID_Validate(t *testing.T) {
	t.Parallel()

	t.Run("when the UUID is valid", func(t *testing.T) {
		t.Parallel()

		subject := &UUID{
			Upper: 0xa967a8b93f9c4918,
			Lower: 0x9a4119577be5fec5,
		}
		if err := subject.Validate(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("when the UUID is invalid", func(t *testing.T) {
		t.Parallel()

		cases := []struct {
			Desc    string
			Subject *UUID
			Expect  string
		}{
			{
				"nil UUID",
				&UUID{},
				"UUID must use version 4",
			},
			{
				"omni UUID",
				&UUID{
					Upper: 0xffffffffffffffff,
					Lower: 0xffffffffffffffff,
				},
				"UUID must use version 4",
			},
			{
				"wrong version",
				&UUID{
					Upper: 0xa967a8b93f9c_f0_18,
					Lower: 0x9a4119577be5fec5,
				},
				"UUID must use version 4",
			},
			{
				"wrong variant",
				&UUID{
					Upper: 0xa967a8b93f9c4918,
					Lower: 0xc0_4119577be5fec5,
				},
				"UUID must use RFC 4122 variant",
			},
		}

		for _, c := range cases {
			c := c // capture loop variable
			t.Run(c.Desc, func(t *testing.T) {
				t.Parallel()

				err := c.Subject.Validate()
				if err == nil {
					t.Fatal("expected an error")
				}
				if err.Error() != c.Expect {
					t.Fatalf("got %q, want %q", err, c.Expect)
				}
			})
		}
	})
}
