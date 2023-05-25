package testapp

import (
	"bytes"
	"context"
	sync "sync"

	"github.com/dogmatiq/dogma"
)

// EventProjection tracks all events produced by the test application.
type EventProjection struct {
	dogma.NoTimeoutHintBehavior
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
		dogma.HandlesEvent[*IntegrationEventA](),
		dogma.HandlesEvent[*IntegrationEventB](),
	)
}

// HandleEvent updates the projection to reflect the occurrence of an event.
func (h *EventProjection) HandleEvent(
	ctx context.Context,
	r, c, n []byte,
	s dogma.ProjectionEventScope,
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
func (h *EventProjection) ResourceVersion(ctx context.Context, r []byte) ([]byte, error) {
	h.m.Lock()
	v := h.resources[string(r)]
	h.m.Unlock()

	return v, nil
}

// CloseResource informs the handler that the engine has no further use for a
// resource.
func (h *EventProjection) CloseResource(ctx context.Context, r []byte) error {
	h.m.Lock()
	delete(h.resources, string(r))
	h.m.Unlock()

	return nil
}
