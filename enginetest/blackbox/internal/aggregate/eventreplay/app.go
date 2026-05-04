package eventreplay

import (
	"encoding/json"

	"github.com/dogmatiq/dogma"
)

// app routes [writeCommand] and [checkCommand] to [handler].
type app struct {
	Handler handler
}

func (a *app) Configure(c dogma.ApplicationConfigurer) {
	c.Identity("app", "a2a00000-0000-0000-0000-000000000000")
	c.Routes(dogma.ViaAggregate(&a.Handler))
}

// handler handles [writeCommand] and [checkCommand] against [root].
type handler struct {
	Got chan string
}

func (*handler) New() *root { return &root{} }

func (*handler) Configure(c dogma.AggregateConfigurer) {
	c.Identity("handler", "a2b00000-0000-0000-0000-000000000000")
	c.Routes(
		dogma.HandlesCommand[*writeCommand](),
		dogma.HandlesCommand[*checkCommand](),
		dogma.RecordsEvent[*valueWritten](),
	)
}

func (*handler) RouteCommandToInstance(dogma.Command) string { return "instance" }

func (h *handler) HandleCommand(r *root, s dogma.AggregateCommandScope[*root], c dogma.Command) {
	switch cmd := c.(type) {
	case *writeCommand:
		s.RecordEvent(&valueWritten{Value: cmd.Value})
	case *checkCommand:
		h.Got <- r.Value
	}
}

// root is the aggregate root. It stores a Value string and implements real
// MarshalBinary / UnmarshalBinary as a snapshot smoke test.
type root struct {
	Value string
}

func (r *root) AggregateInstanceDescription() string { return "" }

func (r *root) ApplyEvent(e dogma.Event) {
	if ev, ok := e.(*valueWritten); ok {
		r.Value = ev.Value
	}
}

func (r *root) MarshalBinary() ([]byte, error) {
	return json.Marshal(r)
}

func (r *root) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, r)
}
