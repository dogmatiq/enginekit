package atomicity

import "github.com/dogmatiq/dogma"

// command is handled by the integration handler under test.
type command struct{}

func (*command) MessageDescription() string                  { return "command" }
func (*command) Validate(dogma.CommandValidationScope) error { return nil }
func (*command) MarshalBinary() ([]byte, error)              { return nil, nil }
func (*command) UnmarshalBinary([]byte) error                { return nil }

// started is the first event recorded by the handler in a single scope. Both
// this event and [completed] must be visible downstream for the atomicity
// guarantee to hold.
type started struct{}

func (*started) MessageDescription() string                { return "started" }
func (*started) Validate(dogma.EventValidationScope) error { return nil }
func (*started) MarshalBinary() ([]byte, error)            { return nil, nil }
func (*started) UnmarshalBinary([]byte) error              { return nil }

// completed is the second event recorded by the handler in the same scope. The
// observer returns satisfied == true on this event.
type completed struct{}

func (*completed) MessageDescription() string                { return "completed" }
func (*completed) Validate(dogma.EventValidationScope) error { return nil }
func (*completed) MarshalBinary() ([]byte, error)            { return nil, nil }
func (*completed) UnmarshalBinary([]byte) error              { return nil }

func init() {
	dogma.RegisterCommand[*command]("0f2f2c51-8df9-45e0-857a-c1d72d70bf8b")
	dogma.RegisterEvent[*started]("43ab5acd-44d0-4b6f-a2bf-e3c738d3c738")
	dogma.RegisterEvent[*completed]("8a7ea022-9b89-47a7-895d-e37143856d34")
}
