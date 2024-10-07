package config_test

import (
	"testing"

	. "github.com/dogmatiq/enginekit/config"
)

func TestIdentity_validate(t *testing.T) {
	cases := []struct {
		Name     string
		Identity Identity
		Want     string
	}{
		{
			"valid",
			Identity{
				Name: "name",
				Key:  "2da5eec5-374e-4716-b1c7-f24abd8df57f",
			},
			"",
		},
		{
			"valid with name containing non-ASCII characters",
			Identity{
				Name: "😀",
				Key:  "79f63053-1ca6-4537-974f-dd0121eb5195",
			},
			"",
		},
		{
			"empty",
			Identity{},
			`invalid identity name (""): names must be non-empty, printable UTF-8 strings with no whitespace` + "\n" +
				`invalid identity key (""): keys must be RFC 4122/9562 UUIDs: invalid UUID format, expected 36 characters`,
		},
		{
			"empty name",
			Identity{
				Key: "c79d01bb-b289-4e5d-b2fd-9779f33b3a19",
			},
			`invalid identity name (""): names must be non-empty, printable UTF-8 strings with no whitespace`,
		},
		{
			"name containing spaces",
			Identity{
				Name: "the name",
				Key:  "c405f1e2-b309-4a43-84bf-5a1f8e7656b8",
			},
			`invalid identity name ("the name"): names must be non-empty, printable UTF-8 strings with no whitespace`,
		},
		{
			"name containing non-printable characters",
			Identity{
				Name: "name\n",
				Key:  "79f63053-1ca6-4537-974f-dd0121eb5195",
			},
			`invalid identity name ("name\n"): names must be non-empty, printable UTF-8 strings with no whitespace`,
		},
		{
			"empty key",
			Identity{
				Name: "name",
			},
			`invalid identity key (""): keys must be RFC 4122/9562 UUIDs: invalid UUID format, expected 36 characters`,
		},
		{
			"non-UUID key",
			Identity{
				Name: "name",
				Key:  "_b4ac052-68b1-4877-974e-c437aceb7f3f",
			},
			`invalid identity key ("_b4ac052-68b1-4877-974e-c437aceb7f3f"): keys must be RFC 4122/9562 UUIDs: invalid UUID format, expected hex digit`,
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			got := ""

			if err := c.Identity.Validate(); err != nil {
				got = err.Error()
			}

			if got != c.Want {
				t.Fatalf("unexpected error: got %q, want %q", got, c.Want)
			}
		})
	}
}
