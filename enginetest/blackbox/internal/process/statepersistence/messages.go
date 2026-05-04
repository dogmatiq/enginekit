package statepersistence

import "github.com/dogmatiq/dogma"

// firstCommand triggers [firstIntegration] to record [firstEvent].
type firstCommand struct{}

func (*firstCommand) MessageDescription() string                  { return "firstCommand" }
func (*firstCommand) Validate(dogma.CommandValidationScope) error { return nil }
func (*firstCommand) MarshalBinary() ([]byte, error)              { return nil, nil }
func (*firstCommand) UnmarshalBinary([]byte) error                { return nil }

// firstEvent is the first event delivered to the process. The process handler
// uses it to set a value on the root.
type firstEvent struct{}

func (*firstEvent) MessageDescription() string                { return "firstEvent" }
func (*firstEvent) Validate(dogma.EventValidationScope) error { return nil }
func (*firstEvent) MarshalBinary() ([]byte, error)            { return nil, nil }
func (*firstEvent) UnmarshalBinary([]byte) error              { return nil }

// secondCommand triggers [secondIntegration] to record [secondEvent].
type secondCommand struct{}

func (*secondCommand) MessageDescription() string                  { return "secondCommand" }
func (*secondCommand) Validate(dogma.CommandValidationScope) error { return nil }
func (*secondCommand) MarshalBinary() ([]byte, error)              { return nil, nil }
func (*secondCommand) UnmarshalBinary([]byte) error                { return nil }

// secondEvent is the second event delivered to the process. The process
// handler reads the root value set during firstEvent to verify persistence.
type secondEvent struct{}

func (*secondEvent) MessageDescription() string                { return "secondEvent" }
func (*secondEvent) Validate(dogma.EventValidationScope) error { return nil }
func (*secondEvent) MarshalBinary() ([]byte, error)            { return nil, nil }
func (*secondEvent) UnmarshalBinary([]byte) error              { return nil }

func init() {
	dogma.RegisterCommand[*firstCommand]("b5000001-0000-0000-0000-000000000000")
	dogma.RegisterEvent[*firstEvent]("b5000002-0000-0000-0000-000000000000")
	dogma.RegisterCommand[*secondCommand]("b5000003-0000-0000-0000-000000000000")
	dogma.RegisterEvent[*secondEvent]("b5000004-0000-0000-0000-000000000000")
}
