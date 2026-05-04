package routing

import "github.com/dogmatiq/dogma"

// triggerCommand causes [triggerIntegration] to emit [observed].
type triggerCommand struct{}

func (*triggerCommand) MessageDescription() string                  { return "triggerCommand" }
func (*triggerCommand) Validate(dogma.CommandValidationScope) error { return nil }
func (*triggerCommand) MarshalBinary() ([]byte, error)              { return nil, nil }
func (*triggerCommand) UnmarshalBinary([]byte) error                { return nil }

// observed is the event that [projection] is registered to handle.
type observed struct{}

func (*observed) MessageDescription() string                { return "observed" }
func (*observed) Validate(dogma.EventValidationScope) error { return nil }
func (*observed) MarshalBinary() ([]byte, error)            { return nil, nil }
func (*observed) UnmarshalBinary([]byte) error              { return nil }

func init() {
	dogma.RegisterCommand[*triggerCommand]("c1000001-0000-0000-0000-000000000000")
	dogma.RegisterEvent[*observed]("c1000002-0000-0000-0000-000000000000")
}
