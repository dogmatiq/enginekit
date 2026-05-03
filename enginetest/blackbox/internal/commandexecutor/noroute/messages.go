package noroute

import "github.com/dogmatiq/dogma"

type command struct{}

func (*command) MessageDescription() string                  { return "command" }
func (*command) Validate(dogma.CommandValidationScope) error { return nil }
func (*command) MarshalBinary() ([]byte, error)              { return nil, nil }
func (*command) UnmarshalBinary([]byte) error                { return nil }

// unroutedCommand is not registered as a route in [app], so executing it must
// cause the engine to panic.
type unroutedCommand struct{}

func (*unroutedCommand) MessageDescription() string                  { return "unrouted command" }
func (*unroutedCommand) Validate(dogma.CommandValidationScope) error { return nil }
func (*unroutedCommand) MarshalBinary() ([]byte, error)              { return nil, nil }
func (*unroutedCommand) UnmarshalBinary([]byte) error                { return nil }

func init() {
	dogma.RegisterCommand[*command]("d6baa55d-1c8b-409a-8d4c-616afe15954f")
	dogma.RegisterCommand[*unroutedCommand]("5b23771f-9be0-4f3e-a272-fec16d1a6293")
}
