package end

import "github.com/dogmatiq/dogma"

// --- messages for the "panics after End" test ---

// endCommand triggers the end sequence.
type endCommand struct{}

func (*endCommand) MessageDescription() string                  { return "endCommand" }
func (*endCommand) Validate(dogma.CommandValidationScope) error { return nil }
func (*endCommand) MarshalBinary() ([]byte, error)              { return nil, nil }
func (*endCommand) UnmarshalBinary([]byte) error                { return nil }

// endTrigger is produced by [endIntegration] and consumed by [endProcess].
type endTrigger struct{}

func (*endTrigger) MessageDescription() string                { return "endTrigger" }
func (*endTrigger) Validate(dogma.EventValidationScope) error { return nil }
func (*endTrigger) MarshalBinary() ([]byte, error)            { return nil, nil }
func (*endTrigger) UnmarshalBinary([]byte) error              { return nil }

// continueCommand is declared as an outbound route so that
// ProcessScope.ExecuteCommand does not panic due to an unknown type.
// It is never actually dispatched (the call panics before queuing it).
type continueCommand struct{}

func (*continueCommand) MessageDescription() string                  { return "continueCommand" }
func (*continueCommand) Validate(dogma.CommandValidationScope) error { return nil }
func (*continueCommand) MarshalBinary() ([]byte, error)              { return nil, nil }
func (*continueCommand) UnmarshalBinary([]byte) error                { return nil }

// --- messages for the "future events ignored after End" test ---

// replayCommand triggers the replay-test chain.
type replayCommand struct{}

func (*replayCommand) MessageDescription() string                  { return "replayCommand" }
func (*replayCommand) Validate(dogma.CommandValidationScope) error { return nil }
func (*replayCommand) MarshalBinary() ([]byte, error)              { return nil, nil }
func (*replayCommand) UnmarshalBinary([]byte) error                { return nil }

// replayTrigger is produced by [replayIntegration] and consumed by
// [replayProcess].
type replayTrigger struct{}

func (*replayTrigger) MessageDescription() string                { return "replayTrigger" }
func (*replayTrigger) Validate(dogma.EventValidationScope) error { return nil }
func (*replayTrigger) MarshalBinary() ([]byte, error)            { return nil, nil }
func (*replayTrigger) UnmarshalBinary([]byte) error              { return nil }

func init() {
	dogma.RegisterCommand[*endCommand]("b4000001-0000-0000-0000-000000000000")
	dogma.RegisterEvent[*endTrigger]("b4000002-0000-0000-0000-000000000000")
	dogma.RegisterCommand[*continueCommand]("b4000003-0000-0000-0000-000000000000")
	dogma.RegisterCommand[*replayCommand]("b4000004-0000-0000-0000-000000000000")
	dogma.RegisterEvent[*replayTrigger]("b4000005-0000-0000-0000-000000000000")
}
