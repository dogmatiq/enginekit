package stubs

import (
	"context"

	"github.com/dogmatiq/dogma"
)

// ProjectionMessageHandlerStub is a test implementation of
// [dogma.ProjectionMessageHandler].
type ProjectionMessageHandlerStub struct {
	ConfigureFunc        func(dogma.ProjectionConfigurer)
	HandleEventFunc      func(context.Context, dogma.ProjectionEventScope, dogma.Event) (uint64, error)
	CheckpointOffsetFunc func(context.Context, string) (uint64, error)
	CompactFunc          func(context.Context, dogma.ProjectionCompactScope) error
}

var _ dogma.ProjectionMessageHandler = &ProjectionMessageHandlerStub{}

// Configure describes the handler's configuration to the engine.
func (h *ProjectionMessageHandlerStub) Configure(c dogma.ProjectionConfigurer) {
	if h.ConfigureFunc != nil {
		h.ConfigureFunc(c)
	}
}

// HandleEvent updates the projection to reflect the occurrence of an event.
func (h *ProjectionMessageHandlerStub) HandleEvent(
	ctx context.Context,
	s dogma.ProjectionEventScope,
	e dogma.Event,
) (uint64, error) {
	if h.HandleEventFunc != nil {
		return h.HandleEventFunc(ctx, s, e)
	}
	return s.Offset() + 1, nil
}

// CheckpointOffset returns the offset at which the handler expects to
// resume handling events from a specific stream.
func (h *ProjectionMessageHandlerStub) CheckpointOffset(
	ctx context.Context,
	id string,
) (uint64, error) {
	if h.CheckpointOffsetFunc != nil {
		return h.CheckpointOffsetFunc(ctx, id)
	}
	return 0, nil
}

// Compact attempts to reduce the size of the projection.
func (h *ProjectionMessageHandlerStub) Compact(
	ctx context.Context,
	s dogma.ProjectionCompactScope,
) error {
	if h.CompactFunc != nil {
		return h.CompactFunc(ctx, s)
	}
	return nil
}
