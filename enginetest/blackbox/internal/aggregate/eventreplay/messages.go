package eventreplay

import "github.com/dogmatiq/dogma"

// writeCommand carries a value to store in the aggregate instance.
type writeCommand struct{ Value string }

func (*writeCommand) MessageDescription() string                  { return "writeCommand" }
func (*writeCommand) Validate(dogma.CommandValidationScope) error { return nil }
func (*writeCommand) MarshalBinary() ([]byte, error)              { return nil, nil }
func (*writeCommand) UnmarshalBinary([]byte) error                { return nil }

// checkCommand asks the handler to report the instance's current value.
type checkCommand struct{}

func (*checkCommand) MessageDescription() string                  { return "checkCommand" }
func (*checkCommand) Validate(dogma.CommandValidationScope) error { return nil }
func (*checkCommand) MarshalBinary() ([]byte, error)              { return nil, nil }
func (*checkCommand) UnmarshalBinary([]byte) error                { return nil }

// valueWritten is the event produced when [writeCommand] is handled. It carries
// the value for content-preservation verification.
type valueWritten struct{ Value string }

func (*valueWritten) MessageDescription() string                { return "valueWritten" }
func (*valueWritten) Validate(dogma.EventValidationScope) error { return nil }
func (*valueWritten) MarshalBinary() ([]byte, error)            { return nil, nil }
func (*valueWritten) UnmarshalBinary([]byte) error              { return nil }

func init() {
	dogma.RegisterCommand[*writeCommand]("a2000001-0000-0000-0000-000000000000")
	dogma.RegisterCommand[*checkCommand]("a2000002-0000-0000-0000-000000000000")
	dogma.RegisterEvent[*valueWritten]("a2000003-0000-0000-0000-000000000000")
}
