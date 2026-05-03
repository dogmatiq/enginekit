package routing

import "github.com/dogmatiq/dogma"

// command is routed to the integration handler under test.
type command struct{}

func (*command) MessageDescription() string                  { return "command" }
func (*command) Validate(dogma.CommandValidationScope) error { return nil }
func (*command) MarshalBinary() ([]byte, error)              { return nil, nil }
func (*command) UnmarshalBinary([]byte) error                { return nil }

func init() {
	dogma.RegisterCommand[*command]("f6d22b37-ecc2-4ff7-a860-5d0e7f5321f1")
}
