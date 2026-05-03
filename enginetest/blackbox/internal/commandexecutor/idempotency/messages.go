package idempotency

import "github.com/dogmatiq/dogma"

// command is submitted with an idempotency key. On the first call the handler
// is invoked; on subsequent calls with the same key it is deduplicated.
type command struct{}

func (*command) MessageDescription() string                  { return "command" }
func (*command) Validate(dogma.CommandValidationScope) error { return nil }
func (*command) MarshalBinary() ([]byte, error)              { return nil, nil }
func (*command) UnmarshalBinary([]byte) error                { return nil }

// handled is recorded when [command] is handled for the first time. On
// deduplicated calls no events are produced, so an observer waiting for this
// event returns ErrEventObserverNotSatisfied.
type handled struct{}

func (*handled) MessageDescription() string                { return "handled" }
func (*handled) Validate(dogma.EventValidationScope) error { return nil }
func (*handled) MarshalBinary() ([]byte, error)            { return nil, nil }
func (*handled) UnmarshalBinary([]byte) error              { return nil }

func init() {
	dogma.RegisterCommand[*command]("63ed0986-344a-4591-a56d-ccea4e380e2a")
	dogma.RegisterEvent[*handled]("cb53f606-eccb-41f9-91bb-c000e79c8499")
}
