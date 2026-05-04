// Package routing tests that the engine routes commands to the correct
// integration handler.
package routing

import (
	"context"
	"testing"
	"time"

	"github.com/dogmatiq/dogma"
)

// command is the message dispatched to the integration handler under test.
// Its Content field lets the test verify that what the handler receives is
// equal in value to what was dispatched.
type command struct{ Content string }

func (c *command) MessageDescription() string                { return "command" }
func (*command) Validate(dogma.CommandValidationScope) error { return nil }
func (c *command) MarshalBinary() ([]byte, error)            { return []byte(c.Content), nil }
func (c *command) UnmarshalBinary(data []byte) error         { c.Content = string(data); return nil }

func init() {
	dogma.RegisterCommand[*command]("f6d22b37-ecc2-4ff7-a860-5d0e7f5321f1")
}

// handler is the integration handler under test. It forwards the received
// command to Received so the test can inspect it.
type handler struct {
	Received chan *command
}

func (*handler) Configure(c dogma.IntegrationConfigurer) {
	c.Identity("handler", "4c6bddfb-c821-47c2-81d3-5c7611020786")
	c.Routes(dogma.HandlesCommand[*command]())
}

func (h *handler) HandleCommand(_ context.Context, _ dogma.IntegrationCommandScope, cmd dogma.Command) error {
	h.Received <- cmd.(*command)
	return nil
}

// app wires the handler into a dogma application.
type app struct{ Handler handler }

func (a *app) Configure(c dogma.ApplicationConfigurer) {
	c.Identity("app", "6c2f7cd6-9d1c-4ef0-a87c-210e21745d0b")
	c.Routes(dogma.ViaIntegration(&a.Handler))
}

// Run runs the integration routing tests.
func Run(
	t *testing.T,
	setup func(t *testing.T, app dogma.Application) dogma.CommandExecutor,
) {
	t.Helper()

	t.Run("it routes the command to the handler", func(t *testing.T) {
		// The handler captures the command it receives. The test verifies:
		//   1. the command's content matches what was dispatched (correct routing)
		//   2. the command is a distinct copy, not the original pointer (Dogma ADR 32)
		a := &app{Handler: handler{Received: make(chan *command, 1)}}
		x := setup(t, a)

		sent := &command{Content: "hello"}
		if err := x.ExecuteCommand(t.Context(), sent); err != nil {
			t.Fatal(err)
		}

		select {
		case received := <-a.Handler.Received:
			if received.Content != sent.Content {
				t.Fatalf("got command content %q, want %q", received.Content, sent.Content)
			}
			if received == sent {
				t.Fatal("engine passed the original command pointer to the handler; expected a distinct copy (Dogma ADR 32)")
			}
		case <-time.After(5 * time.Second):
			t.Fatal("timed out waiting for handler to be called")
		}
	})
}
