package observer

import "github.com/dogmatiq/dogma"

// command is submitted by the test to initiate the scenario under test.
type command struct{}

func (*command) MessageDescription() string                  { return "command" }
func (*command) Validate(dogma.CommandValidationScope) error { return nil }
func (*command) MarshalBinary() ([]byte, error)              { return nil, nil }
func (*command) UnmarshalBinary([]byte) error                { return nil }

// triggered is recorded by the aggregate when it handles [command]. The process
// consumes it to continue the causal chain.
type triggered struct{}

func (*triggered) MessageDescription() string                { return "triggered" }
func (*triggered) Validate(dogma.EventValidationScope) error { return nil }
func (*triggered) MarshalBinary() ([]byte, error)            { return nil, nil }
func (*triggered) UnmarshalBinary([]byte) error              { return nil }

// relayed is the command dispatched by the process in response to [triggered].
type relayed struct{}

func (*relayed) MessageDescription() string                  { return "relayed" }
func (*relayed) Validate(dogma.CommandValidationScope) error { return nil }
func (*relayed) MarshalBinary() ([]byte, error)              { return nil, nil }
func (*relayed) UnmarshalBinary([]byte) error                { return nil }

// observed is recorded by the integration handler in response to [relayed]. It
// marks the end of the causal chain and is the event the observer waits for.
type observed struct{}

func (*observed) MessageDescription() string                { return "observed" }
func (*observed) Validate(dogma.EventValidationScope) error { return nil }
func (*observed) MarshalBinary() ([]byte, error)            { return nil, nil }
func (*observed) UnmarshalBinary([]byte) error              { return nil }

func init() {
	dogma.RegisterCommand[*command]("4a573023-836c-4281-ade8-ceb497894d95")
	dogma.RegisterEvent[*triggered]("59f3c164-ad1b-4277-86dc-71e6f66f8536")
	dogma.RegisterCommand[*relayed]("3babc712-32d0-40a6-ba54-ff3c07bb6cfc")
	dogma.RegisterEvent[*observed]("2f4eaf0b-7191-40cc-9175-98e69796f9f5")
}
