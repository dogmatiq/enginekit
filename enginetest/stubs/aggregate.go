package stubs

import (
	"encoding/json"

	"github.com/dogmatiq/dogma"
)

// AggregateRootStub is a test implementation of [dogma.AggregateRoot].
type AggregateRootStub struct {
	AppliedEvents  []dogma.Event     `json:"applied_events,omitempty"`
	ApplyEventFunc func(dogma.Event) `json:"-"`
}

var _ dogma.AggregateRoot = &AggregateRootStub{}

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
type AggregateMessageHandlerStub struct {
	NewFunc                    func() dogma.AggregateRoot
	ConfigureFunc              func(dogma.AggregateConfigurer)
	RouteCommandToInstanceFunc func(dogma.Command) string
	HandleCommandFunc          func(dogma.AggregateRoot, dogma.AggregateCommandScope, dogma.Command)
}

var _ dogma.AggregateMessageHandler = &AggregateMessageHandlerStub{}

// Configure describes the handler's configuration to the engine.
func (h *AggregateMessageHandlerStub) Configure(c dogma.AggregateConfigurer) {
	if h.ConfigureFunc != nil {
		h.ConfigureFunc(c)
	}
}

// New returns an aggregate root instance in its initial state.
func (h *AggregateMessageHandlerStub) New() dogma.AggregateRoot {
	if h.NewFunc != nil {
		return h.NewFunc()
	}
	return &AggregateRootStub{}
}

// RouteCommandToInstance returns the ID of the instance that handles a specific
// command.
func (h *AggregateMessageHandlerStub) RouteCommandToInstance(c dogma.Command) string {
	if h.RouteCommandToInstanceFunc == nil {
		panic(dogma.UnexpectedMessage)
	}
	return h.RouteCommandToInstanceFunc(c)
}

// HandleCommand executes business logic in response to a command.
func (h *AggregateMessageHandlerStub) HandleCommand(
	r dogma.AggregateRoot,
	s dogma.AggregateCommandScope,
	c dogma.Command,
) {
	if h.HandleCommandFunc != nil {
		h.HandleCommandFunc(r, s, c)
	}
}
