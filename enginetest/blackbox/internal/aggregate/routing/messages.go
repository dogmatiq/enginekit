package routing

import "github.com/dogmatiq/dogma"

// command is routed to the aggregate handler under test.
type command struct{}

func (*command) MessageDescription() string                  { return "command" }
func (*command) Validate(dogma.CommandValidationScope) error { return nil }
func (*command) MarshalBinary() ([]byte, error)              { return nil, nil }
func (*command) UnmarshalBinary([]byte) error                { return nil }

// recorded is the event produced by the handler to confirm the command arrived.
type recorded struct{}

func (*recorded) MessageDescription() string                { return "recorded" }
func (*recorded) Validate(dogma.EventValidationScope) error { return nil }
func (*recorded) MarshalBinary() ([]byte, error)            { return nil, nil }
func (*recorded) UnmarshalBinary([]byte) error              { return nil }

func init() {
	dogma.RegisterCommand[*command]("a1000001-0000-0000-0000-000000000000")
	dogma.RegisterEvent[*recorded]("a1000002-0000-0000-0000-000000000000")
}
