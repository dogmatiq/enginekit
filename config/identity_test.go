package config_test

import (
	"testing"

	. "github.com/dogmatiq/enginekit/config"
	. "github.com/dogmatiq/enginekit/internal/test"
)

func TestIdentity_validation(t *testing.T) {
	cases := []struct {
		Name     string
		Want     string
		Identity Identity
	}{

		{
			"valid",
			``, // no error
			Identity{
				AsConfigured: IdentityProperties{
					Name: "name",
					Key:  "2da5eec5-374e-4716-b1c7-f24abd8df57f",
				},
			},
		},
		{
			"valid with name containing non-ASCII characters",
			``, // no error
			Identity{
				AsConfigured: IdentityProperties{
					Name: "😀",
					Key:  "79f63053-1ca6-4537-974f-dd0121eb5195",
				},
			},
		},
		{
			"empty",
			`identity is invalid:` +
				"\n" + `- invalid name (""), expected a non-empty, printable UTF-8 string with no whitespace` +
				"\n" + `- invalid key (""), expected an RFC 4122/9562 UUID`,
			Identity{},
		},
		{
			"partial",
			`partial identity:name/1e6264cd-8df7-49e0-b246-c7e17e70ccdf is invalid: some configuration is potentially missing`,
			Identity{
				AsConfigured: IdentityProperties{
					// NOTE(jmalloc): In practice it's nonsensical to have a
					// partial identity that has non-empty values for both name
					// and key. I suppose it's possible via static analysis to
					// have access to _portions_ of a name or key, but at time
					// of writing we don't do anything so sophisticated.
					Name:     "name",
					Key:      "1e6264cd-8df7-49e0-b246-c7e17e70ccdf",
					Fidelity: Fidelity{IsPartial: true},
				},
			},
		},
		{
			"spectulative",
			`speculative identity:name/e6b691dd-731c-4c14-8e1c-1622381202dc is invalid: conditions for the component's inclusion in the configuration could not be evaluated`,
			Identity{
				AsConfigured: IdentityProperties{
					Name:     "name",
					Key:      "e6b691dd-731c-4c14-8e1c-1622381202dc",
					Fidelity: Fidelity{IsSpeculative: true},
				},
			},
		},
		{
			"unresolved",
			`unresolved identity:name/e6b691dd-731c-4c14-8e1c-1622381202dc is invalid: configuration includes values that could not be evaluated`,
			Identity{
				AsConfigured: IdentityProperties{
					Name:     "name",
					Key:      "e6b691dd-731c-4c14-8e1c-1622381202dc",
					Fidelity: Fidelity{IsUnresolved: true},
				},
			},
		},
		{
			"empty name",
			`identity:""/c79d01bb-b289-4e5d-b2fd-9779f33b3a19 is invalid: invalid name (""), expected a non-empty, printable UTF-8 string with no whitespace`,
			Identity{
				AsConfigured: IdentityProperties{
					Key: "c79d01bb-b289-4e5d-b2fd-9779f33b3a19",
				},
			},
		},
		{
			"name containing spaces",
			`identity:"the name"/c405f1e2-b309-4a43-84bf-5a1f8e7656b8 is invalid: invalid name ("the name"), expected a non-empty, printable UTF-8 string with no whitespace`,
			Identity{
				AsConfigured: IdentityProperties{
					Name: "the name",
					Key:  "c405f1e2-b309-4a43-84bf-5a1f8e7656b8",
				},
			},
		},
		{
			"name containing non-printable characters",
			`identity:"name\n"/79f63053-1ca6-4537-974f-dd0121eb5195 is invalid: invalid name ("name\n"), expected a non-empty, printable UTF-8 string with no whitespace`,
			Identity{
				AsConfigured: IdentityProperties{
					Name: "name\n",
					Key:  "79f63053-1ca6-4537-974f-dd0121eb5195",
				},
			},
		},
		{
			"empty key",
			`identity:name/"" is invalid: invalid key (""), expected an RFC 4122/9562 UUID`,
			Identity{
				AsConfigured: IdentityProperties{
					Name: "name",
				},
			},
		},
		{
			"non-UUID key",
			`identity:name/_b4ac052-68b1-4877-974e-c437aceb7f3f is invalid: invalid key ("_b4ac052-68b1-4877-974e-c437aceb7f3f"), expected an RFC 4122/9562 UUID`,
			Identity{
				AsConfigured: IdentityProperties{
					Name: "name",
					Key:  "_b4ac052-68b1-4877-974e-c437aceb7f3f",
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			got := ""
			if _, err := Normalize(c.Identity); err != nil {
				got = err.Error()
			}

			if c.Want != got {
				t.Log("unexpected error:")
				t.Log("  got:  ", got)
				t.Log("  want: ", c.Want)
				t.FailNow()
			}
		})
	}
}

func TestIdentity_normalize(t *testing.T) {
	id := Identity{
		AsConfigured: IdentityProperties{
			Name: "name",
			Key:  "0EB1E0A1-B067-4625-A7DC-D7D260B0AFAB",
		},
	}

	got, err := Normalize(id)
	if err != nil {
		t.Fatal(err)
	}

	want := Identity{
		AsConfigured: IdentityProperties{
			Name: "name",
			Key:  "0eb1e0a1-b067-4625-a7dc-d7d260b0afab",
		},
	}

	Expect(
		t,
		"unexpected identity",
		got,
		want,
	)

	id.AsConfigured.Name = ""

	if _, err = Normalize(id); err == nil {
		t.Fatal("expected an error")
	}
}

func TestIdentity_String(t *testing.T) {
	cases := []struct {
		Name     string
		Want     string
		Identity Identity
	}{
		{
			"valid, canonical",
			`identity:name/2da5eec5-374e-4716-b1c7-f24abd8df57f`,
			Identity{
				AsConfigured: IdentityProperties{
					Name: "name",
					Key:  "2da5eec5-374e-4716-b1c7-f24abd8df57f",
				},
			},
		},
		{
			"valid, non-canonical",
			`identity:name/2da5eec5-374e-4716-b1c7-f24abd8df57f`,
			Identity{
				AsConfigured: IdentityProperties{
					Name: "name",
					Key:  "2DA5EEC5-374E-4716-B1C7-F24ABD8DF57F",
				},
			},
		},
		{
			"valid with name containing non-ASCII characters",
			`identity:😀/79f63053-1ca6-4537-974f-dd0121eb5195`,
			Identity{
				AsConfigured: IdentityProperties{
					Name: "😀",
					Key:  "79f63053-1ca6-4537-974f-dd0121eb5195",
				},
			},
		},
		{
			"empty",
			`identity`,
			Identity{},
		},
		{
			"empty name",
			`identity:""/c79d01bb-b289-4e5d-b2fd-9779f33b3a19`,
			Identity{
				AsConfigured: IdentityProperties{
					Key: "c79d01bb-b289-4e5d-b2fd-9779f33b3a19",
				},
			},
		},
		{
			"name containing spaces",
			`identity:"the name"/c405f1e2-b309-4a43-84bf-5a1f8e7656b8`,
			Identity{
				AsConfigured: IdentityProperties{
					Name: "the name",
					Key:  "c405f1e2-b309-4a43-84bf-5a1f8e7656b8",
				},
			},
		},
		{
			"name containing non-printable characters",
			`identity:"name\n"/79f63053-1ca6-4537-974f-dd0121eb5195`,
			Identity{
				AsConfigured: IdentityProperties{
					Name: "name\n",
					Key:  "79f63053-1ca6-4537-974f-dd0121eb5195",
				},
			},
		},
		{
			"empty key",
			`identity:name/""`,
			Identity{
				AsConfigured: IdentityProperties{
					Name: "name",
				},
			},
		},
		{
			"non-UUID key",
			`identity:name/_b4ac052-68b1-4877-974e-c437aceb7f3f`,
			Identity{
				AsConfigured: IdentityProperties{
					Name: "name",
					Key:  "_b4ac052-68b1-4877-974e-c437aceb7f3f",
				},
			},
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
