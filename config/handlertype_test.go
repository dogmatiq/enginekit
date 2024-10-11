package config

import "testing"

func TestSwitchByHandlerTypeOf(t *testing.T) {
	cases := []struct {
		Handler Handler
		Want    string
	}{
		{
			&Aggregate{},
			"aggregate",
		},
		{
			&Process{},
			"process",
		},
		{
			&Integration{},
			"integration",
		},
		{
			&Projection{},
			"projection",
		},
	}

	for _, c := range cases {
		var got string

		SwitchByHandlerTypeOf(
			c.Handler,
			func(*Aggregate) { got = "aggregate" },
			func(*Process) { got = "process" },
			func(*Integration) { got = "integration" },
			func(*Projection) { got = "projection" },
		)

		if got != c.Want {
			t.Errorf("unexpected value: got %q, want %q", got, c.Want)
		}
	}
}
