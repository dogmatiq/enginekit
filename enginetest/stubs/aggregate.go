package stubs

import (
	"encoding/json"

	"github.com/dogmatiq/dogma"
)

// AggregateRootStub is a test implementation of [dogma.AggregateRoot].
type AggregateRootStub struct {
	AppliedEvents                    []dogma.Event     `json:"applied_events,omitempty"`
	ApplyEventFunc                   func(dogma.Event) `json:"-"`
	AggregateInstanceDescriptionFunc func() string     `json:"-"`
}

var _ dogma.AggregateRoot = &AggregateRootStub{}

// AggregateInstanceDescription returns a human-readable description of the
// aggregate instance's current state.
func (r *AggregateRootStub) AggregateInstanceDescription() string {
	if r.AggregateInstanceDescriptionFunc != nil {
		return r.AggregateInstanceDescriptionFunc()
	}
	return ""
}

// ApplyEvent updates aggregate instance to reflect the occurrence of an event.
func (r *AggregateRootStub) ApplyEvent(e dogma.Event) {
	r.AppliedEvents = append(r.AppliedEvents, e)

	if r.ApplyEventFunc != nil {
		r.ApplyEventFunc(e)
	}
}

// MarshalBinary implements [encoding.BinaryMarshaler].
func (r *AggregateRootStub) MarshalBinary() ([]byte, error) {
	return json.Marshal(r)
}

// UnmarshalBinary implements [encoding.BinaryUnmarshaler].
func (r *AggregateRootStub) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, r)
}

// AggregateMessageHandlerStub is a test implementation of
// [dogma.AggregateMessageHandler].
type AggregateMessageHandlerStub[R dogma.AggregateRoot] struct {
	NewFunc                    func() R
	ConfigureFunc              func(dogma.AggregateConfigurer)
	RouteCommandToInstanceFunc func(dogma.Command) string
	HandleCommandFunc          func(R, dogma.AggregateCommandScope[R], dogma.Command)
}

var _ dogma.AggregateMessageHandler[*AggregateRootStub] = &AggregateMessageHandlerStub[*AggregateRootStub]{}

// Configure describes the handler's configuration to the engine.
func (h *AggregateMessageHandlerStub[R]) Configure(c dogma.AggregateConfigurer) {
	if h.ConfigureFunc != nil {
		h.ConfigureFunc(c)
	}
}

// New returns an aggregate root instance in its initial state.
func (h *AggregateMessageHandlerStub[R]) New() R {
	if h.NewFunc != nil {
		return h.NewFunc()
	}
	return newRoot[R]()
}

// RouteCommandToInstance returns the ID of the instance that handles a specific
// command.
func (h *AggregateMessageHandlerStub[R]) RouteCommandToInstance(c dogma.Command) string {
	if h.RouteCommandToInstanceFunc == nil {
		panic(dogma.UnexpectedMessage)
	}
	return h.RouteCommandToInstanceFunc(c)
}

// HandleCommand executes business logic in response to a command.
func (h *AggregateMessageHandlerStub[R]) HandleCommand(
	r R,
	s dogma.AggregateCommandScope[R],
	c dogma.Command,
) {
	if h.HandleCommandFunc != nil {
		h.HandleCommandFunc(r, s, c)
	}
}
