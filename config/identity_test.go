package config_test

import (
	"testing"

	. "github.com/dogmatiq/enginekit/config"
	. "github.com/dogmatiq/enginekit/internal/test"
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
			if _, err := Normalize(c.Identity); err != nil {
				got = err.Error()
			}

			if got != c.Want {
				t.Fatalf("unexpected error: got %q, want %q", got, c.Want)
			}
		})
	}
}

func TestIdentity_normalize(t *testing.T) {
	id := Identity{
		Name: "name",
		Key:  "0EB1E0A1-B067-4625-A7DC-D7D260B0AFAB",
	}

	got, err := Normalize(id)
	if err != nil {
		t.Fatal(err)
	}

	want := Identity{
		Name: "name",
		Key:  "0eb1e0a1-b067-4625-a7dc-d7d260b0afab",
	}

	Expect(
		t,
		"unexpected identity",
		got,
		want,
	)

	id.Name = ""

	if _, err = Normalize(id); err == nil {
		t.Fatal("expected an error")
	}
}

func TestIdentity_String(t *testing.T) {
	cases := []struct {
		Name     string
		Identity Identity
		Want     string
	}{
		{
			"valid, canonical",
			Identity{
				Name: "name",
				Key:  "2da5eec5-374e-4716-b1c7-f24abd8df57f",
			},
			"name/2da5eec5-374e-4716-b1c7-f24abd8df57f",
		},
		{
			"valid, non-canonical",
			Identity{
				Name: "name",
				Key:  "2DA5EEC5-374E-4716-B1C7-F24ABD8DF57F",
			},
			"name/2da5eec5-374e-4716-b1c7-f24abd8df57f",
		},
		{
			"valid with name containing non-ASCII characters",
			Identity{
				Name: "😀",
				Key:  "79f63053-1ca6-4537-974f-dd0121eb5195",
			},
			"😀/79f63053-1ca6-4537-974f-dd0121eb5195",
		},
		{
			"empty",
			Identity{},
			`?/?`,
		},
		{
			"empty name",
			Identity{
				Key: "c79d01bb-b289-4e5d-b2fd-9779f33b3a19",
			},
			`?/c79d01bb-b289-4e5d-b2fd-9779f33b3a19`,
		},
		{
			"name containing spaces",
			Identity{
				Name: "the name",
				Key:  "c405f1e2-b309-4a43-84bf-5a1f8e7656b8",
			},
			`"the name"/c405f1e2-b309-4a43-84bf-5a1f8e7656b8`,
		},
		{
			"name containing non-printable characters",
			Identity{
				Name: "name\n",
				Key:  "79f63053-1ca6-4537-974f-dd0121eb5195",
			},
			`"name\n"/79f63053-1ca6-4537-974f-dd0121eb5195`,
		},
		{
			"empty key",
			Identity{
				Name: "name",
			},
			`name/?`,
		},
		{
			"non-UUID key",
			Identity{
				Name: "name",
				Key:  "_b4ac052-68b1-4877-974e-c437aceb7f3f",
			},
			`name/"_b4ac052-68b1-4877-974e-c437aceb7f3f"`,
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			got := c.Identity.String()

			if got != c.Want {
				t.Fatalf("unexpected string: got %q, want %q", got, c.Want)
			}
		})
	}
}
