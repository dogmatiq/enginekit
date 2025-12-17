package identitypb_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/dogmatiq/enginekit/internal/test"
	. "github.com/dogmatiq/enginekit/protobuf/identitypb"
	uuidpb "github.com/dogmatiq/enginekit/protobuf/uuidpb"
	proto "google.golang.org/protobuf/proto"
)

func TestIdentity_Validate(t *testing.T) {
	t.Parallel()

	t.Run("when the identity is valid", func(t *testing.T) {
		t.Parallel()

		subject := New("<name>", uuidpb.Generate())
		if err := subject.Validate(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("when the identity is invalid", func(t *testing.T) {
		t.Parallel()

		cases := []struct {
			Desc    string
			Subject *Identity
			Expect  string
		}{
			{
				"too short",
				&Identity{
					Name: "",
					Key:  uuidpb.Generate(),
				},
				"invalid name: must be between 1 and 255 bytes",
			},
			{
				"too long",
				&Identity{
					Name: strings.Repeat("*", 256),
					Key:  uuidpb.Generate(),
				},
				"invalid name: must be between 1 and 255 bytes",
			},
			{
				"invalid key",
				&Identity{
					Name: "<name>",
					Key:  &uuidpb.UUID{},
				},
				"invalid key: UUID must use version 4 or 5",
			},
		}

		for _, c := range cases {
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

func TestIdentity_ParseAndMustParse(t *testing.T) {
	t.Parallel()

	t.Run("when the identity is valid", func(t *testing.T) {
		t.Run("Parse() returns the identity", func(t *testing.T) {
			t.Parallel()

			actual, err := Parse(
				"<name>",
				"a967a8b9-3f9c-4918-9a41-19577be5fec5",
			)
			if err != nil {
				t.Fatal(err)
			}

			expect := &Identity{
				Name: "<name>",
				Key: &uuidpb.UUID{
					Upper: 0xa967a8b93f9c4918,
					Lower: 0x9a4119577be5fec5,
				},
			}

			if !actual.Equal(expect) {
				t.Fatalf("unexpected identiy: got %s, want %s", actual, expect)
			}
		})

		t.Run("MustParse() returns the identity", func(t *testing.T) {
			t.Parallel()

			actual := MustParse(
				"<name>",
				"a967a8b9-3f9c-4918-9a41-19577be5fec5",
			)

			expect := &Identity{
				Name: "<name>",
				Key: &uuidpb.UUID{
					Upper: 0xa967a8b93f9c4918,
					Lower: 0x9a4119577be5fec5,
				},
			}

			if !actual.Equal(expect) {
				t.Fatalf("unexpected identiy: got %s, want %s", actual, expect)
			}
		})
	})

	t.Run("when the identity is invalid", func(t *testing.T) {
		t.Parallel()

		cases := []struct {
			Desc   string
			Name   string
			Key    string
			Expect string
		}{
			{
				"too short",
				"",
				"a967a8b9-3f9c-4918-9a41-19577be5fec5",
				"invalid name: must be between 1 and 255 bytes",
			},
			{
				"too long",
				strings.Repeat("*", 256),
				"a967a8b9-3f9c-4918-9a41-19577be5fec5",
				"invalid name: must be between 1 and 255 bytes",
			},
			{
				"malformed key",
				"<name>",
				"not-a-uuid",
				"invalid key: invalid UUID format, expected 36 characters",
			},
			{
				"invalid key",
				"<name>",
				"00000000-0000-0000-0000-000000000000",
				"invalid key: UUID must use version 4 or 5",
			},
		}

		for _, c := range cases {
			t.Run(c.Desc, func(t *testing.T) {
				t.Parallel()

				t.Run("Parse() returns an error", func(t *testing.T) {
					t.Parallel()

					_, err := Parse(c.Name, c.Key)
					if err == nil {
						t.Fatal("expected an error")
					}
					if err.Error() != c.Expect {
						t.Fatalf("got %q, want %q", err, c.Expect)
					}
				})

				t.Run("MustParse() panics", func(t *testing.T) {
					t.Parallel()

					test.ExpectPanic(
						t,
						c.Expect,
						func() {
							MustParse(c.Name, c.Key)
						},
					)
				})
			})
		}
	})
}

func TestIdentity_Format(t *testing.T) {
	t.Parallel()

	subject := New(
		"<name>",
		&uuidpb.UUID{
			Upper: 0xa967a8b93f9c4918,
			Lower: 0x9a4119577be5fec5,
		},
	)

	cases := []struct {
		Desc   string
		Format string
		Want   string
	}{
		{
			"string",
			"%s",
			`<name>/a967a8b9-3f9c-4918-9a41-19577be5fec5`,
		},
		{
			"go string",
			"%#v",
			`identitypb.New("<name>", uuidpb.MustParse("a967a8b9-3f9c-4918-9a41-19577be5fec5"))`,
		},
		{
			"fallback",
			"%v",
			`&{{{} [] [] <nil>} <name> &{{{} [] [] <nil>} 12206910828600641816 11115193218858614469 [] 0} [] 0}`, // depends on protobuf internals, unfortunately
		},
	}

	for _, c := range cases {
		t.Run(c.Desc, func(t *testing.T) {
			t.Parallel()

			actual := fmt.Sprintf(c.Format, subject)
			if actual != c.Want {
				t.Fatalf("got %q, want %q", actual, c.Want)
			}
		})
	}
}

func TestIdentity_Equal(t *testing.T) {
	t.Parallel()

	a := &Identity{
		Name: "<a>",
		Key:  uuidpb.Generate(),
	}

	b1 := &Identity{
		Name: "<b>",
		Key:  uuidpb.Generate(),
	}

	b2 := &Identity{
		Name: b1.Name,
		Key:  proto.Clone(b1.Key).(*uuidpb.UUID),
	}

	if !a.Equal(a) {
		t.Fatal("expected a == a")
	}

	if a.Equal(b1) {
		t.Fatal("did not expect a == b1")
	}

	if !b1.Equal(b2) {
		t.Fatal("expected b1 == b2")
	}
}
