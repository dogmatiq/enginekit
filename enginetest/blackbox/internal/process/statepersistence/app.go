package statepersistence

import (
	"context"
	"encoding/json"

	"github.com/dogmatiq/dogma"
)

// app drives the state persistence test:
// [firstIntegration] records [firstEvent] → [process] sets root.Value →
// [secondIntegration] records [secondEvent] → [process] reads root.Value →
// sends to handler.Got.
type app struct {
	Handler handler
}

func (a *app) Configure(c dogma.ApplicationConfigurer) {
	c.Identity("app", "b5a00000-0000-0000-0000-000000000000")
	c.Routes(
		dogma.ViaIntegration(&firstIntegration{}),
		dogma.ViaIntegration(&secondIntegration{}),
		dogma.ViaProcess(&a.Handler),
	)
}

// firstIntegration handles [firstCommand] and records [firstEvent].
type firstIntegration struct{}

func (*firstIntegration) Configure(c dogma.IntegrationConfigurer) {
	c.Identity("first-integration", "b5b00000-0000-0000-0000-000000000000")
	c.Routes(
		dogma.HandlesCommand[*firstCommand](),
		dogma.RecordsEvent[*firstEvent](),
	)
}

func (*firstIntegration) HandleCommand(_ context.Context, s dogma.IntegrationCommandScope, _ dogma.Command) error {
	s.RecordEvent(&firstEvent{})
	return nil
}

// secondIntegration handles [secondCommand] and records [secondEvent].
type secondIntegration struct{}

func (*secondIntegration) Configure(c dogma.IntegrationConfigurer) {
	c.Identity("second-integration", "b5c00000-0000-0000-0000-000000000000")
	c.Routes(
		dogma.HandlesCommand[*secondCommand](),
		dogma.RecordsEvent[*secondEvent](),
	)
}

func (*secondIntegration) HandleCommand(_ context.Context, s dogma.IntegrationCommandScope, _ dogma.Command) error {
	s.RecordEvent(&secondEvent{})
	return nil
}

// handler is the process under test. It sets root.Value when handling
// [firstEvent] and sends root.Value to Got when handling [secondEvent]. If
// state is persisted correctly, the second invocation sees the value set by
// the first.
type handler struct {
	dogma.NoTimeoutMessagesBehavior[*root]
	Got chan string
}

func (*handler) New() *root { return &root{} }

func (*handler) Configure(c dogma.ProcessConfigurer) {
	c.Identity("process", "b5d00000-0000-0000-0000-000000000000")
	c.Routes(
		dogma.HandlesEvent[*firstEvent](),
		dogma.HandlesEvent[*secondEvent](),
		dogma.ExecutesCommand[*firstCommand](),
	)
}

func (*handler) RouteEventToInstance(context.Context, dogma.Event) (string, bool, error) {
	return "instance", true, nil
}

func (h *handler) HandleEvent(_ context.Context, r *root, s dogma.ProcessEventScope[*root], e dogma.Event) error {
	switch e.(type) {
	case *firstEvent:
		r.Value = "persisted-value"
	case *secondEvent:
		h.Got <- r.Value
	}
	return nil
}

// root is the process root. It implements real MarshalBinary / UnmarshalBinary
// so that the engine's persistence path (snapshot or event replay) is
// exercised. The test does not assert that any specific persistence mechanism
// was used — only that the state is correct.
type root struct {
	Value string
}

func (*root) ProcessInstanceDescription(bool) string { return "" }

func (r *root) MarshalBinary() ([]byte, error) {
	return json.Marshal(r)
}

func (r *root) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, r)
}
