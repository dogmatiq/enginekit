package testapp

import (
	"context"
	"errors"
	"fmt"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/enginetest/internal/action"
)

type integrationA struct{}

func (h *integrationA) Configure(c dogma.IntegrationConfigurer) {
	c.Identity("integration-a", "e5a20886-23bc-4948-badd-0f6930b7130a")
	c.Routes(
		dogma.HandlesCommand[*IntegrationCommandA](),
		dogma.RecordsEvent[*IntegrationEventA](),
	)
}

func (h *integrationA) HandleCommand(
	_ context.Context,
	s dogma.IntegrationCommandScope,
	c dogma.Command,
) error {
	switch c := c.(type) {
	case *IntegrationCommandA:
		return action.Run(s, c)
	default:
		panic(dogma.UnexpectedMessage)
	}
}

type integrationB struct{}

func (h *integrationB) Configure(c dogma.IntegrationConfigurer) {
	c.Identity("integration-b", "45794cc7-216b-4f15-9abc-dcac6a1eb3a5")
	c.Routes(
		dogma.HandlesCommand[*IntegrationCommandB](),
		dogma.RecordsEvent[*IntegrationEventB](),
	)
}

func (h *integrationB) HandleCommand(
	_ context.Context,
	s dogma.IntegrationCommandScope,
	c dogma.Command,
) error {
	switch c := c.(type) {
	case *IntegrationCommandB:
		return action.Run(s, c)
	default:
		panic(dogma.UnexpectedMessage)
	}
}

// ErrInvalidIntegrationMessage is returned by the IntegrationCommand.Validate()
// if the IsInvalid flag is true.
var ErrInvalidIntegrationMessage = errors.New("integration message is invalid")

// MessageDescription returns a human-readable description of the message.
func (x *IntegrationCommandA) MessageDescription() string {
	return "integration command A"
}

// Validate returns an error if the message is invalid.
func (x *IntegrationCommandA) Validate(dogma.CommandValidationScope) error {
	if x.IsInvalid {
		return ErrInvalidIntegrationMessage
	}
	return nil
}

// MessageDescription returns a human-readable description of the message.
func (x *IntegrationCommandB) MessageDescription() string {
	return "integration command B"
}

// Validate returns an error if the message is invalid.
func (x *IntegrationCommandB) Validate(dogma.CommandValidationScope) error {
	if x.IsInvalid {
		return ErrInvalidIntegrationMessage
	}
	return nil
}

// MessageDescription returns a human-readable description of the message.
func (x *IntegrationEventA) MessageDescription() string {
	return fmt.Sprintf("integration event A: %s", x.Value)
}

// Validate returns an error if the message is invalid.
func (x *IntegrationEventA) Validate(dogma.EventValidationScope) error {
	return nil
}

// MessageDescription returns a human-readable description of the message.
func (x *IntegrationEventB) MessageDescription() string {
	return fmt.Sprintf("integration event B: %s", x.Value)
}

// Validate returns an error if the message is invalid.
func (x *IntegrationEventB) Validate(dogma.EventValidationScope) error {
	return nil
}

func init() {
	dogma.RegisterCommand[*IntegrationCommandA]("7f0d9dfc-e936-460e-b32e-a80913ed8211")
	dogma.RegisterCommand[*IntegrationCommandB]("d47d4e83-ab5f-434d-a9b8-4a2306f231fb")
	dogma.RegisterEvent[*IntegrationEventA]("f3d33db7-9061-4d17-aae0-a072cd78b3c4")
	dogma.RegisterEvent[*IntegrationEventB]("c0743533-63af-4cc2-9f61-59edd8baaa0b")
}
