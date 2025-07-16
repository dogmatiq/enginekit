package uuidpb_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/dogmatiq/dapper"
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

func TestParse(t *testing.T) {
	t.Parallel()

	t.Run("when the string is a valid UUID", func(t *testing.T) {
		t.Parallel()

		cases := []struct {
			Desc   string
			String string
		}{
			{"lowercase", "a967a8b9-3f9c-4918-9a41-19577be5fec5"},
			{"uppercase", "A967A8B9-3F9C-4918-9A41-19577BE5FEC5"},
		}

		for _, c := range cases {
			c := c // capture loop variable
			t.Run(c.Desc, func(t *testing.T) {
				t.Parallel()

				expect := &UUID{
					Upper: 0xa967a8b93f9c4918,
					Lower: 0x9a4119577be5fec5,
				}
				actual, err := Parse(c.String)
				if err != nil {
					t.Fatal(err)
				}

				if !proto.Equal(actual, expect) {
					t.Fatalf("got %q, want %q", actual, expect)
				}
			})
		}
	})

	t.Run("when the string is not a valid UUID", func(t *testing.T) {
		t.Parallel()

		cases := []struct {
			Desc   string
			String string
		}{
			{"empty string", ""},
			{"too short", "3493af5d-e4d0-4f3b-a73d-048e6b08496"},
			{"too long", "3493af5d-e4d0-4f3b-a73d-048e6b08496ab"},
			{"no hyphens", "a967a8b93f9c49189a4119577be5fec5"},
			{"no hyphens at position 8", "7e770248_7336-4cee-881d-f24013e6c1bf"},
			{"no hyphens at position 13", "7e770248-7336_4cee-881d-f24013e6c1bf"},
			{"no hyphens at position 18", "7e770248-7336-4cee_881d-f24013e6c1bf"},
			{"no hyphens at position 23", "7e770248-7336-4cee-881d_f24013e6c1bf"},
			{"non-hex character", "26c4a622-e4b4-4Xe7-8454-e3dc90f7d1d8"},
		}

		for _, c := range cases {
			c := c // capture loop variable
			t.Run(c.Desc, func(t *testing.T) {
				t.Parallel()

				_, err := Parse(c.String)
				if err == nil {
					t.Fatal("expected an error")
				}
			})
		}
	})
}

func TestMustParse(t *testing.T) {
	t.Run("when the string is a valid UUID", func(t *testing.T) {
		t.Parallel()

		expect := &UUID{
			Upper: 0xa967a8b93f9c4918,
			Lower: 0x9a4119577be5fec5,
		}
		actual := MustParse("a967a8b9-3f9c-4918-9a41-19577be5fec5")

		if !proto.Equal(actual, expect) {
			t.Fatalf("got %q, want %q", actual, expect)
		}
	})

	t.Run("when the string is not a valid UUID", func(t *testing.T) {
		t.Parallel()

		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expected a panic")
			}
		}()

		MustParse("invalid")
	})
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

func TestCopyBytes(t *testing.T) {
	t.Parallel()

	subject := &UUID{
		Upper: 0xa967a8b93f9c4918,
		Lower: 0x9a4119577be5fec5,
	}

	t.Run("fills the target slice with the UUID's bytes", func(t *testing.T) {
		var actual [16]byte
		n := CopyBytes(subject, actual[:])
		expect := []byte{
			0xa9, 0x67, 0xa8, 0xb9,
			0x3f, 0x9c, 0x49, 0x18,
			0x9a, 0x41, 0x19, 0x57,
			0x7b, 0xe5, 0xfe, 0xc5,
		}

		if !bytes.Equal(actual[:], expect) {
			t.Fatalf("got %q, want %q", actual, expect)
		}

		if n != 16 {
			t.Fatalf("got %d, want 16", n)
		}
	})

	t.Run("it truncates to the target slice's length", func(t *testing.T) {
		var actual [8]byte
		n := CopyBytes(subject, actual[:])
		expect := []byte{
			0xa9, 0x67, 0xa8, 0xb9,
			0x3f, 0x9c, 0x49, 0x18,
		}

		if !bytes.Equal(actual[:], expect) {
			t.Fatalf("got %q, want %q", actual, expect)
		}

		if n != 8 {
			t.Fatalf("got %d, want 8", n)
		}
	})
}

func TestUUID_AsBytes(t *testing.T) {
	t.Parallel()

	cases := []struct {
		Desc    string
		Subject *UUID
		Expect  []byte
	}{
		{
			"nil",
			nil,
			[]byte{
				0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00,
			},
		},
		{
			"zero",
			&UUID{},
			[]byte{
				0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00,
			},
		},
		{
			"non-zero",
			&UUID{
				Upper: 0xa967a8b93f9c4918,
				Lower: 0x9a4119577be5fec5,
			},
			[]byte{
				0xa9, 0x67, 0xa8, 0xb9,
				0x3f, 0x9c, 0x49, 0x18,
				0x9a, 0x41, 0x19, 0x57,
				0x7b, 0xe5, 0xfe, 0xc5,
			},
		},
	}

	for _, c := range cases {
		c := c // capture loop variable
		t.Run(c.Desc, func(t *testing.T) {
			t.Parallel()

			actual := c.Subject.AsBytes()

			if !bytes.Equal(actual, c.Expect) {
				t.Fatalf("got %q, want %q", actual, c.Expect)
			}
		})
	}
}

func TestUUID_AsString(t *testing.T) {
	t.Parallel()

	cases := []struct {
		Desc    string
		Subject *UUID
		Expect  string
	}{
		{"nil", nil, "00000000-0000-0000-0000-000000000000"},
		{"zero", &UUID{}, "00000000-0000-0000-0000-000000000000"},
		{
			"non-zero",
			&UUID{
				Upper: 0xa967a8b93f9c4918,
				Lower: 0x9a4119577be5fec5,
			},
			"a967a8b9-3f9c-4918-9a41-19577be5fec5",
		},
	}

	for _, c := range cases {
		c := c // capture loop variable
		t.Run(c.Desc, func(t *testing.T) {
			t.Parallel()

			actual := c.Subject.AsString()

			if actual != c.Expect {
				t.Fatalf("got %q, want %q", actual, c.Expect)
			}
		})
	}
}

func TestUUID_DapperString(t *testing.T) {
	t.Parallel()

	subject := &UUID{
		Upper: 0xa967a8b93f9c4918,
		Lower: 0x9a4119577be5fec5,
	}

	actual := dapper.Format(subject)
	expect := "*github.com/dogmatiq/enginekit/protobuf/uuidpb.UUID [a967a8b9-3f9c-4918-9a41-19577be5fec5]"

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
				"UUID must use RFC 9562 variant",
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

func TestUUID_Equal(t *testing.T) {
	t.Parallel()

	a := &UUID{
		Upper: 0xa967a8b93f9c4918,
		Lower: 0x9a4119577be5fec5,
	}

	b := &UUID{
		Upper: 0x3f9c4918a967a8b9,
		Lower: 0x7be5fec59a411957,
	}

	if a.Equal(b) {
		t.Fatal("did not expect a == b")
	}

	if !a.Equal(a) {
		t.Fatal("did not expect a != b")
	}
}

func TestUUID_Less(t *testing.T) {
	t.Parallel()

	a := &UUID{
		Upper: 0xa967a8b93f9c4918,
		Lower: 0x9a4119577be5fec5,
	}

	b := &UUID{
		Upper: 0x3f9c4918a967a8b9,
		Lower: 0x7be5fec59a411957,
	}

	if a.Less(b) {
		t.Fatal("did not expect a < b")
	}

	if !b.Less(a) {
		t.Fatal("did not expect b >= a")
	}

	if a.Less(a) {
		t.Fatal("did not expect a < b")
	}
}

func TestUUID_Compare(t *testing.T) {
	t.Parallel()

	t.Run("compares both the upper and lower parts", func(t *testing.T) {
		a := &UUID{
			Upper: 0xa967a8b93f9c4918,
			Lower: 0x9a4119577be5fec5,
		}

		b := &UUID{
			Upper: 0x3f9c4918a967a8b9,
			Lower: 0x7be5fec59a411957,
		}

		c := &UUID{
			Upper: 0x3f9c4918a967a8b9,
			Lower: 0x0000fec59a411957,
		}

		if a.Compare(b) < 0 {
			t.Fatal("did not expect a < b")
		}

		if b.Compare(a) >= 0 {
			t.Fatal("did not expect b >= a")
		}

		if a.Compare(a) != 0 {
			t.Fatal("did not expect a != b")
		}

		if b.Compare(c) <= 0 {
			t.Fatal("did not expect b <= c")
		}
	})

	t.Run("subtraction overflow regression", func(t *testing.T) {
		a, err := Parse("e6929f89-acb1-4994-8248-d097ad20da5f")
		if err != nil {
			t.Fatal(err)
		}

		b, err := Parse("46471d80-5b50-44a0-a8a8-261941336291")
		if err != nil {
			t.Fatal(err)
		}

		if a.Compare(b) <= 0 {
			t.Fatal("did not expect a <= b")
		}
	})
}
