package fixtures

import (
	"context"
	"time"

	"github.com/dogmatiq/dogma"
)

// ProjectionMessageHandler is a test implementation of dogma.ProjectionMessageHandler.
type ProjectionMessageHandler struct {
	ConfigureFunc   func(dogma.ProjectionConfigurer)
	HandleEventFunc func(context.Context, dogma.ProjectionEventScope, dogma.Message, []byte, []byte) error
	RecoverFunc     func(context.Context, []byte) (v []byte, ok bool, err error)
	DiscardFunc     func(context.Context, []byte) error
	TimeoutHintFunc func(m dogma.Message) time.Duration
}

var _ dogma.ProjectionMessageHandler = &ProjectionMessageHandler{}

// Configure configures the behavior of the engine as it relates to this
// handler.
//
// c provides access to the various configuration options, such as specifying
// which types of event messages are routed to this handler.
//
// If h.ConfigureFunc is non-nil, it calls h.ConfigureFunc(c)
func (h *ProjectionMessageHandler) Configure(c dogma.ProjectionConfigurer) {
	if h.ConfigureFunc != nil {
		h.ConfigureFunc(c)
	}
}

// HandleEvent handles a domain event message that has been routed to this
// handler.
//
// s provides access to the operations available within the scope of handling
// m.
//
// It panics with the UnexpectedMessage value if m is not one of the event
// types that is routed to this handler via Configure().
//
// If h.HandleEventFunc is non-nil it calls h.HandleEventFunc(ctx, s, m, k, v).
func (h *ProjectionMessageHandler) HandleEvent(
	ctx context.Context,
	s dogma.ProjectionEventScope,
	m dogma.Message,
	k, v []byte,
) error {
	if h.HandleEventFunc != nil {
		return h.HandleEventFunc(ctx, s, m, k, v)
	}

	return nil
}

// Recover returns the value component of a key/value association persisted
// by a call to HandleEvent().
//
// If h.RecoverFunc is non-nil it calls h.RecoverFunc(ctx, k).
func (h *ProjectionMessageHandler) Recover(ctx context.Context, k []byte) (v []byte, ok bool, err error) {
	if h.RecoverFunc != nil {
		return h.RecoverFunc(ctx, k)
	}

	return nil, true, nil
}

// Discard informs the projection that a specific key/value association is
// no longer required.
//
// If h.DiscardFunc is non-nil it calls h.DiscardFunc(ctx, k).
func (h *ProjectionMessageHandler) Discard(ctx context.Context, k []byte) error {
	if h.DiscardFunc != nil {
		return h.DiscardFunc(ctx, k)
	}

	return nil
}

// TimeoutHint returns a duration that is suitable for computing a deadline
// for the handling of the given message by this handler.
//
// If h.TimeoutHintFunc is non-nil it calls h.TimeoutHintFunc(m).
func (h *ProjectionMessageHandler) TimeoutHint(m dogma.Message) time.Duration {
	if h.TimeoutHintFunc != nil {
		return h.TimeoutHintFunc(m)
	}

	return 0
}
