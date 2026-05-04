package timeout

import "github.com/dogmatiq/dogma"

// startCommand triggers the causal chain.
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

// workflowTimeout is the timeout message scheduled by [process]. The engine
// must deliver it at or after the scheduled time.
type workflowTimeout struct{}

func (*workflowTimeout) MessageDescription() string                  { return "workflowTimeout" }
func (*workflowTimeout) Validate(dogma.TimeoutValidationScope) error { return nil }
func (*workflowTimeout) MarshalBinary() ([]byte, error)              { return nil, nil }
func (*workflowTimeout) UnmarshalBinary([]byte) error                { return nil }

func init() {
	dogma.RegisterCommand[*startCommand]("b3000001-0000-0000-0000-000000000000")
	dogma.RegisterEvent[*workflowStarted]("b3000002-0000-0000-0000-000000000000")
	dogma.RegisterTimeout[*workflowTimeout]("b3000003-0000-0000-0000-000000000000")
}
