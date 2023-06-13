package identitypb_test

import (
	"strings"
	"testing"

	. "github.com/dogmatiq/enginekit/protobuf/identitypb"
	uuidpb "github.com/dogmatiq/enginekit/protobuf/uuidpb"
)

func TestIdentity_Validate(t *testing.T) {
	t.Parallel()

	t.Run("when the identity is valid", func(t *testing.T) {
		t.Parallel()

		valid := &Identity{
			Name: "<name>",
			Key:  uuidpb.Generate(),
		}
		if err := valid.Validate(); err != nil {
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
				"invalid identity name: must be between 1 and 255 bytes",
			},
			{
				"too long",
				&Identity{
					Name: strings.Repeat("*", 256),
					Key:  uuidpb.Generate(),
				},
				"invalid identity name: must be between 1 and 255 bytes",
			},
			{
				"nil-uuid key",
				&Identity{
					Name: "<name>",
					Key:  uuidpb.Nil(),
				},
				"invalid identity key: must not be the nil UUID (all zeroes)",
			},
			{
				"omni-uuid key",
				&Identity{
					Name: "<name>",
					Key:  uuidpb.Omni(),
				},
				"invalid identity key: must not be the omni UUID (all ones)",
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
