package atomicity

import "github.com/dogmatiq/dogma"

// command is handled by the aggregate under test.
type command struct{}

func (*command) MessageDescription() string                  { return "command" }
func (*command) Validate(dogma.CommandValidationScope) error { return nil }
func (*command) MarshalBinary() ([]byte, error)              { return nil, nil }
func (*command) UnmarshalBinary([]byte) error                { return nil }

// begun is the first event recorded in the scope. Both this event and
// [finished] must be visible downstream for the atomicity guarantee to hold.
type begun struct{}

func (*begun) MessageDescription() string                { return "begun" }
func (*begun) Validate(dogma.EventValidationScope) error { return nil }
func (*begun) MarshalBinary() ([]byte, error)            { return nil, nil }
func (*begun) UnmarshalBinary([]byte) error              { return nil }

// finished is the second event recorded in the same scope. The observer
// returns satisfied == true on this event.
type finished struct{}

func (*finished) MessageDescription() string                { return "finished" }
func (*finished) Validate(dogma.EventValidationScope) error { return nil }
func (*finished) MarshalBinary() ([]byte, error)            { return nil, nil }
func (*finished) UnmarshalBinary([]byte) error              { return nil }

func init() {
	dogma.RegisterCommand[*command]("a3000001-0000-0000-0000-000000000000")
	dogma.RegisterEvent[*begun]("a3000002-0000-0000-0000-000000000000")
	dogma.RegisterEvent[*finished]("a3000003-0000-0000-0000-000000000000")
}
