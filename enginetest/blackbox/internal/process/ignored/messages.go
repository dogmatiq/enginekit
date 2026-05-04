package ignored

import "github.com/dogmatiq/dogma"

// triggerCommand causes [triggerIntegration] to emit [processTrigger].
type triggerCommand struct{}

func (*triggerCommand) MessageDescription() string                  { return "triggerCommand" }
func (*triggerCommand) Validate(dogma.CommandValidationScope) error { return nil }
func (*triggerCommand) MarshalBinary() ([]byte, error)              { return nil, nil }
func (*triggerCommand) UnmarshalBinary([]byte) error                { return nil }

// processTrigger is the event the process is registered to handle.
type processTrigger struct{}

func (*processTrigger) MessageDescription() string                { return "processTrigger" }
func (*processTrigger) Validate(dogma.EventValidationScope) error { return nil }
func (*processTrigger) MarshalBinary() ([]byte, error)            { return nil, nil }
func (*processTrigger) UnmarshalBinary([]byte) error              { return nil }

func init() {
	dogma.RegisterCommand[*triggerCommand]("b2000001-0000-0000-0000-000000000000")
	dogma.RegisterEvent[*processTrigger]("b2000002-0000-0000-0000-000000000000")
}
