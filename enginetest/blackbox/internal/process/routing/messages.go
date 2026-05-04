package routing

import "github.com/dogmatiq/dogma"

// startCommand triggers the causal chain under test.
type startCommand struct{}

func (*startCommand) MessageDescription() string                  { return "startCommand" }
func (*startCommand) Validate(dogma.CommandValidationScope) error { return nil }
func (*startCommand) MarshalBinary() ([]byte, error)              { return nil, nil }
func (*startCommand) UnmarshalBinary([]byte) error                { return nil }

// workflowStarted is produced by [startIntegration] and consumed by [process].
type workflowStarted struct{}

func (*workflowStarted) MessageDescription() string                { return "workflowStarted" }
func (*workflowStarted) Validate(dogma.EventValidationScope) error { return nil }
func (*workflowStarted) MarshalBinary() ([]byte, error)            { return nil, nil }
func (*workflowStarted) UnmarshalBinary([]byte) error              { return nil }

// actionCommand is executed by [process] and handled by [completionIntegration].
type actionCommand struct{}

func (*actionCommand) MessageDescription() string                  { return "actionCommand" }
func (*actionCommand) Validate(dogma.CommandValidationScope) error { return nil }
func (*actionCommand) MarshalBinary() ([]byte, error)              { return nil, nil }
func (*actionCommand) UnmarshalBinary([]byte) error                { return nil }

// workflowCompleted is produced by [completionIntegration] to mark the end of
// the chain. The test observes this event to confirm routing succeeded.
type workflowCompleted struct{}

func (*workflowCompleted) MessageDescription() string                { return "workflowCompleted" }
func (*workflowCompleted) Validate(dogma.EventValidationScope) error { return nil }
func (*workflowCompleted) MarshalBinary() ([]byte, error)            { return nil, nil }
func (*workflowCompleted) UnmarshalBinary([]byte) error              { return nil }

func init() {
	dogma.RegisterCommand[*startCommand]("b1000001-0000-0000-0000-000000000000")
	dogma.RegisterEvent[*workflowStarted]("b1000002-0000-0000-0000-000000000000")
	dogma.RegisterCommand[*actionCommand]("b1000003-0000-0000-0000-000000000000")
	dogma.RegisterEvent[*workflowCompleted]("b1000004-0000-0000-0000-000000000000")
}
