package testapp

import (
	"bytes"
	"context"
	sync "sync"

	"github.com/dogmatiq/dogma"
	action "github.com/dogmatiq/enginekit/enginetest/internal/action"
)

// EventProjection tracks all events produced by the test application.
type EventProjection struct {
	dogma.NoCompactBehavior

	m         sync.Mutex
	resources map[string][]byte
	wait      chan struct{}
	events    []dogma.Event
}

// Range calls fn for each event that the application records until fn returns
// false or ctx is canceled.
func (h *EventProjection) Range(
	ctx context.Context,
	fn func(dogma.Event) bool,
) error {
	index := 0

	for {
		h.m.Lock()
		if h.wait == nil {
			h.wait = make(chan struct{})
		}
		wait := h.wait
		events := h.events[index:]
		index += len(events)
		h.m.Unlock()

		for _, e := range events {
			if fn(e) {
				return nil
			}
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-wait:
		}
	}
}

// Configure describes the handler's configuration to the engine.
func (h *EventProjection) Configure(c dogma.ProjectionConfigurer) {
	c.Identity("events", "aa0d1129-4713-42e9-aad5-bf95f024c6aa")
	c.Routes(
		dogma.HandlesEvent[*GenericEvent](),
		dogma.HandlesEvent[*IntegrationEventA](),
		dogma.HandlesEvent[*IntegrationEventB](),
		dogma.HandlesEvent[*ProcessEventA](),
	)
}

// HandleEvent updates the projection to reflect the occurrence of an event.
func (h *EventProjection) HandleEvent(
	_ context.Context,
	r, c, n []byte,
	_ dogma.ProjectionEventScope,
	e dogma.Event,
) (ok bool, err error) {
	h.m.Lock()
	defer h.m.Unlock()

	v := h.resources[string(r)]
	if !bytes.Equal(v, c) {
		return false, nil
	}

	if h.resources == nil {
		h.resources = map[string][]byte{}
	}
	h.resources[string(r)] = n
	h.events = append(h.events, e)

	if h.wait != nil {
		close(h.wait)
		h.wait = nil
	}

	return true, nil
}

// ResourceVersion returns the current version of a resource.
func (h *EventProjection) ResourceVersion(_ context.Context, r []byte) ([]byte, error) {
	h.m.Lock()
	v := h.resources[string(r)]
	h.m.Unlock()

	return v, nil
}

// CloseResource informs the handler that the engine has no further use for a
// resource.
func (h *EventProjection) CloseResource(_ context.Context, r []byte) error {
	h.m.Lock()
	delete(h.resources, string(r))
	h.m.Unlock()

	return nil
}

// MessageDescription returns a human-readable description of the message.
func (x *DoActions) MessageDescription() string {
	return "performing an action"
}

// Validate returns an error if the message is invalid.
func (x *DoActions) Validate(dogma.CommandValidationScope) error {
	return nil
}

// MessageDescription returns a human-readable description of the message.
func (x *GenericEvent) MessageDescription() string {
	return x.Value
}

// Validate returns an error if the message is invalid.
func (x *GenericEvent) Validate(dogma.EventValidationScope) error {
	return nil
}

type actionExecutor struct{}

func (h *actionExecutor) Configure(c dogma.IntegrationConfigurer) {
	c.Identity("action-executor", "ea8dcfb4-d37e-45e1-ac92-c0775d5cf277")
	c.Routes(
		dogma.HandlesCommand[*DoActions](),
		dogma.RecordsEvent[*GenericEvent](),
		dogma.RecordsEvent[*ProcessEventA](),
	)
}

func (h *actionExecutor) HandleCommand(
	_ context.Context,
	s dogma.IntegrationCommandScope,
	c dogma.Command,
) error {
	switch c := c.(type) {
	case *DoActions:
		return action.Run(s, c)
	default:
		panic(dogma.UnexpectedMessage)
	}
}
