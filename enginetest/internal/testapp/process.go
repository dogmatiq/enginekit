package testapp

import (
	"context"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/enginetest/internal/action"
)

type processA struct {
	dogma.NoTimeoutMessagesBehavior
	dogma.StatelessProcessBehavior
}

func (h *processA) Configure(c dogma.ProcessConfigurer) {
	c.Identity("process-a", "999cdaa2-9323-49fe-ba07-95920ad1837b")
	c.Routes(
		dogma.HandlesEvent[*ProcessEventA](),
		dogma.ExecutesCommand[*DoActions](),
	)
}

func (h *processA) RouteEventToInstance(
	_ context.Context,
	e dogma.Event,
) (id string, ok bool, err error) {
	switch e := e.(type) {
	case *ProcessEventA:
		return e.GetInstanceId(), e.GetInstanceId() != "", nil
	default:
		panic(dogma.UnexpectedMessage)
	}
}

func (h *processA) HandleEvent(
	_ context.Context,
	_ dogma.ProcessRoot,
	s dogma.ProcessEventScope,
	e dogma.Event,
) error {
	switch e := e.(type) {
	case *ProcessEventA:
		return action.Run(s, e)
	default:
		panic(dogma.UnexpectedMessage)
	}
}

// MessageDescription returns a human-readable description of the message.
func (x *ProcessEventA) MessageDescription() string {
	return "process event A"
}

// Validate returns an error if the message is invalid.
func (x *ProcessEventA) Validate(dogma.EventValidationScope) error {
	return nil
}
