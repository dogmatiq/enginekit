package identitypb_test

import (
	"fmt"
	"strings"
	"testing"

	. "github.com/dogmatiq/enginekit/protobuf/identitypb"
	uuidpb "github.com/dogmatiq/enginekit/protobuf/uuidpb"
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
				"invalid key: UUID must use version 4",
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

func TestIdentity_Format(t *testing.T) {
	t.Parallel()

	subject := New(
		"<name>",
		&uuidpb.UUID{
			Upper: 0xa967a8b93f9c4918,
			Lower: 0x9a4119577be5fec5,
		},
	)

	expect := "<name>/a967a8b9-3f9c-4918-9a41-19577be5fec5"
	actual := fmt.Sprintf("%s", subject)

	if actual != expect {
		t.Fatalf("got %q, want %q", actual, expect)
	}
}
