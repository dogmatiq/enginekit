package config_test

import (
	"testing"

	. "github.com/dogmatiq/enginekit/config"
	. "github.com/dogmatiq/enginekit/internal/test"
	"github.com/dogmatiq/enginekit/optional"
)

func TestIdentity_String(t *testing.T) {
	cases := []struct {
		Name     string
		Want     string
		Identity *Identity
	}{
		{
			"valid, canonical",
			`identity:name/2da5eec5-374e-4716-b1c7-f24abd8df57f`,
			&Identity{
				AsConfigured: IdentityAsConfigured{
					Name: optional.Some("name"),
					Key:  optional.Some("2da5eec5-374e-4716-b1c7-f24abd8df57f"),
				},
			},
		},
		{
			"valid, non-canonical",
			`identity:name/2da5eec5-374e-4716-b1c7-f24abd8df57f`,
			&Identity{
				AsConfigured: IdentityAsConfigured{
					Name: optional.Some("name"),
					Key:  optional.Some("2DA5EEC5-374E-4716-B1C7-F24ABD8DF57F"),
				},
			},
		},
		{
			"valid with name containing non-ASCII characters",
			`identity:😀/79f63053-1ca6-4537-974f-dd0121eb5195`,
			&Identity{
				AsConfigured: IdentityAsConfigured{
					Name: optional.Some("😀"),
					Key:  optional.Some("79f63053-1ca6-4537-974f-dd0121eb5195"),
				},
			},
		},
		{
			"empty",
			`identity`,
			&Identity{},
		},
		{
			"empty name",
			`identity:""/c79d01bb-b289-4e5d-b2fd-9779f33b3a19`,
			&Identity{
				AsConfigured: IdentityAsConfigured{
					Name: optional.Some(""),
					Key:  optional.Some("c79d01bb-b289-4e5d-b2fd-9779f33b3a19"),
				},
			},
		},
		{
			"missing name",
			`identity:""/c79d01bb-b289-4e5d-b2fd-9779f33b3a19`,
			&Identity{
				AsConfigured: IdentityAsConfigured{
					Key: optional.Some("c79d01bb-b289-4e5d-b2fd-9779f33b3a19"),
				},
			},
		},
		{
			"name containing spaces",
			`identity:"the name"/c405f1e2-b309-4a43-84bf-5a1f8e7656b8`,
			&Identity{
				AsConfigured: IdentityAsConfigured{
					Name: optional.Some("the name"),
					Key:  optional.Some("c405f1e2-b309-4a43-84bf-5a1f8e7656b8"),
				},
			},
		},
		{
			"name containing non-printable characters",
			`identity:"name\n"/79f63053-1ca6-4537-974f-dd0121eb5195`,
			&Identity{
				AsConfigured: IdentityAsConfigured{
					Name: optional.Some("name\n"),
					Key:  optional.Some("79f63053-1ca6-4537-974f-dd0121eb5195"),
				},
			},
		},
		{
			"empty key",
			`identity:name/""`,
			&Identity{
				AsConfigured: IdentityAsConfigured{
					Name: optional.Some("name"),
					Key:  optional.Some(""),
				},
			},
		},
		{
			"missing key",
			`identity:name/""`,
			&Identity{
				AsConfigured: IdentityAsConfigured{
					Name: optional.Some("name"),
				},
			},
		},
		{
			"non-UUID key",
			`identity:name/_b4ac052-68b1-4877-974e-c437aceb7f3f`,
			&Identity{
				AsConfigured: IdentityAsConfigured{
					Name: optional.Some("name"),
					Key:  optional.Some("_b4ac052-68b1-4877-974e-c437aceb7f3f"),
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

func TestIdentity_validation(t *testing.T) {
	cases := []validationTestCase{
		{
			Name:   "valid",
			Expect: ``,
			Component: &Identity{
				AsConfigured: IdentityAsConfigured{
					Name: optional.Some("name"),
					Key:  optional.Some("2da5eec5-374e-4716-b1c7-f24abd8df57f"),
				},
			},
		},
		{
			Name:   "valid with name containing non-ASCII characters",
			Expect: ``,
			Component: &Identity{
				AsConfigured: IdentityAsConfigured{
					Name: optional.Some("😀"),
					Key:  optional.Some("79f63053-1ca6-4537-974f-dd0121eb5195"),
				},
			},
		},
		{
			Name:      "empty",
			Expect:    `identity is invalid: configuration includes values that could not be evaluated`,
			Component: &Identity{},
		},
		{
			Name:   "partial",
			Expect: `identity:name/e6b691dd-731c-4c14-8e1c-1622381202dc is invalid: could not evaluate entire configuration`,
			Component: &Identity{
				AsConfigured: IdentityAsConfigured{
					Name:     optional.Some("name"),
					Key:      optional.Some("e6b691dd-731c-4c14-8e1c-1622381202dc"),
					Fidelity: Fidelity{IsPartial: true},
				},
			},
		},
		{
			Name:   "spectulative",
			Expect: `identity:name/e6b691dd-731c-4c14-8e1c-1622381202dc is invalid: conditions for the component's inclusion in the configuration could not be evaluated`,
			Component: &Identity{
				AsConfigured: IdentityAsConfigured{
					Name:     optional.Some("name"),
					Key:      optional.Some("e6b691dd-731c-4c14-8e1c-1622381202dc"),
					Fidelity: Fidelity{IsSpeculative: true},
				},
			},
		},
		{
			Name:   "unresolved",
			Expect: `identity:name/e6b691dd-731c-4c14-8e1c-1622381202dc is invalid: configuration includes values that could not be evaluated`,
			Component: &Identity{
				// It's possibly non-sensical to have an identity that contains
				// both it's name and key be considered unresolved, but this
				// allows us to represent a case where the name and key are
				// build dynamically and we don't have the _entire_ string.
				AsConfigured: IdentityAsConfigured{
					Name:     optional.Some("name"),
					Key:      optional.Some("e6b691dd-731c-4c14-8e1c-1622381202dc"),
					Fidelity: Fidelity{IsUnresolved: true},
				},
			},
		},
		{
			Name:   "empty name",
			Expect: `identity:""/c79d01bb-b289-4e5d-b2fd-9779f33b3a19 is invalid: invalid name (""), expected a non-empty, printable UTF-8 string with no whitespace`,
			Component: &Identity{
				AsConfigured: IdentityAsConfigured{
					Name: optional.Some(""),
					Key:  optional.Some("c79d01bb-b289-4e5d-b2fd-9779f33b3a19"),
				},
			},
		},
		{
			Name:   "name containing spaces",
			Expect: `identity:"the name"/c405f1e2-b309-4a43-84bf-5a1f8e7656b8 is invalid: invalid name ("the name"), expected a non-empty, printable UTF-8 string with no whitespace`,
			Component: &Identity{
				AsConfigured: IdentityAsConfigured{
					Name: optional.Some("the name"),
					Key:  optional.Some("c405f1e2-b309-4a43-84bf-5a1f8e7656b8"),
				},
			},
		},
		{
			Name:   "name containing non-printable characters",
			Expect: `identity:"name\n"/79f63053-1ca6-4537-974f-dd0121eb5195 is invalid: invalid name ("name\n"), expected a non-empty, printable UTF-8 string with no whitespace`,
			Component: &Identity{
				AsConfigured: IdentityAsConfigured{
					Name: optional.Some("name\n"),
					Key:  optional.Some("79f63053-1ca6-4537-974f-dd0121eb5195"),
				},
			},
		},
		{
			Name:   "empty key",
			Expect: `identity:name/"" is invalid: invalid key (""), expected an RFC 4122/9562 UUID`,
			Component: &Identity{
				AsConfigured: IdentityAsConfigured{
					Name: optional.Some("name"),
					Key:  optional.Some(""),
				},
			},
		},
		{
			Name:   "non-UUID key",
			Expect: `identity:name/_b4ac052-68b1-4877-974e-c437aceb7f3f is invalid: invalid key ("_b4ac052-68b1-4877-974e-c437aceb7f3f"), expected an RFC 4122/9562 UUID`,
			Component: &Identity{
				AsConfigured: IdentityAsConfigured{
					Name: optional.Some("name"),
					Key:  optional.Some("_b4ac052-68b1-4877-974e-c437aceb7f3f"),
				},
			},
		},
	}

	runValidationTests(t, cases)
}

func TestIdentity_normalize(t *testing.T) {
	id := &Identity{
		AsConfigured: IdentityAsConfigured{
			Name: optional.Some("name"),
			Key:  optional.Some("0EB1E0A1-B067-4625-A7DC-D7D260B0AFAB"),
		},
	}

	got, err := Normalize(id)
	if err != nil {
		t.Fatal(err)
	}

	want := &Identity{
		AsConfigured: IdentityAsConfigured{
			Name: optional.Some("name"),
			Key:  optional.Some("0eb1e0a1-b067-4625-a7dc-d7d260b0afab"),
		},
	}

	Expect(
		t,
		"unexpected identity",
		got,
		want,
	)

	id.AsConfigured.Name = optional.None[string]()

	if _, err = Normalize(id); err == nil {
		t.Fatal("expected an error")
	}
}
