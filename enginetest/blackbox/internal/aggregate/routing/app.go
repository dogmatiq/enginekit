package routing

import "github.com/dogmatiq/dogma"

// app routes [command] to its [handler].
type app struct{}

func (*app) Configure(c dogma.ApplicationConfigurer) {
	c.Identity("app", "a1a00000-0000-0000-0000-000000000000")
	c.Routes(dogma.ViaAggregate(&handler{}))
}

// handler handles [command] and records [recorded].
type handler struct{}

type root struct{ dogma.NoSnapshotBehavior }

func (*root) AggregateInstanceDescription() string { return "" }
func (*root) ApplyEvent(dogma.Event)               {}

func (*handler) New() *root { return &root{} }

func (*handler) Configure(c dogma.AggregateConfigurer) {
	c.Identity("handler", "a1b00000-0000-0000-0000-000000000000")
	c.Routes(
		dogma.HandlesCommand[*command](),
		dogma.RecordsEvent[*recorded](),
	)
}

func (*handler) RouteCommandToInstance(dogma.Command) string { return "instance" }

func (*handler) HandleCommand(_ *root, s dogma.AggregateCommandScope[*root], _ dogma.Command) {
	s.RecordEvent(&recorded{})
}
