package atomicity

import "github.com/dogmatiq/dogma"

// app routes [command] to [handler].
type app struct{}

func (*app) Configure(c dogma.ApplicationConfigurer) {
	c.Identity("app", "a3a00000-0000-0000-0000-000000000000")
	c.Routes(dogma.ViaAggregate(&handler{}))
}

// handler handles [command] and records [begun] and [finished] in the same
// scope, exercising the atomicity guarantee.
type handler struct{}

type root struct{ dogma.NoSnapshotBehavior }

func (*root) AggregateInstanceDescription() string { return "" }
func (*root) ApplyEvent(dogma.Event)               {}

func (*handler) New() *root { return &root{} }

func (*handler) Configure(c dogma.AggregateConfigurer) {
	c.Identity("handler", "a3b00000-0000-0000-0000-000000000000")
	c.Routes(
		dogma.HandlesCommand[*command](),
		dogma.RecordsEvent[*begun](),
		dogma.RecordsEvent[*finished](),
	)
}

func (*handler) RouteCommandToInstance(dogma.Command) string { return "instance" }

func (*handler) HandleCommand(_ *root, s dogma.AggregateCommandScope[*root], _ dogma.Command) {
	s.RecordEvent(&begun{})
	s.RecordEvent(&finished{})
}
