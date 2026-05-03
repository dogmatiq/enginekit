package observer

import (
	"context"

	"github.com/dogmatiq/dogma"
)

// chainApp routes [command] through a three-handler causal chain:
// [chainAggregate] records [triggered] → [chainProcess] executes [relayed] →
// [chainIntegration] records [observed].
type chainApp struct{}

func (*chainApp) Configure(c dogma.ApplicationConfigurer) {
	c.Identity("app", "c8391adc-f2ac-40dd-896e-319678123170")
	c.Routes(
		dogma.ViaAggregate(&chainAggregate{}),
		dogma.ViaProcess(&chainProcess{}),
		dogma.ViaIntegration(&chainIntegration{}),
	)
}

// chainAggregate handles [command] and records [triggered].
type chainAggregate struct{}
type chainAggregateRoot struct{ dogma.NoSnapshotBehavior }

func (*chainAggregateRoot) AggregateInstanceDescription() string { return "" }
func (*chainAggregateRoot) ApplyEvent(dogma.Event)               {}

func (*chainAggregate) New() *chainAggregateRoot { return &chainAggregateRoot{} }

func (*chainAggregate) Configure(c dogma.AggregateConfigurer) {
	c.Identity("aggregate", "85980075-1095-4c24-a5aa-8925814a7137")
	c.Routes(
		dogma.HandlesCommand[*command](),
		dogma.RecordsEvent[*triggered](),
	)
}

func (*chainAggregate) RouteCommandToInstance(dogma.Command) string { return "instance" }

func (*chainAggregate) HandleCommand(_ *chainAggregateRoot, s dogma.AggregateCommandScope[*chainAggregateRoot], _ dogma.Command) {
	s.RecordEvent(&triggered{})
}

// chainProcess handles [triggered] and executes [relayed].
type chainProcess struct {
	dogma.StatelessProcessBehavior
	dogma.NoTimeoutMessagesBehavior[dogma.StatelessProcessRoot]
}

func (*chainProcess) Configure(c dogma.ProcessConfigurer) {
	c.Identity("process", "19615953-cf78-4e46-8651-c4fad3cd4042")
	c.Routes(
		dogma.HandlesEvent[*triggered](),
		dogma.ExecutesCommand[*relayed](),
	)
}

func (*chainProcess) RouteEventToInstance(context.Context, dogma.Event) (string, bool, error) {
	return "instance", true, nil
}

func (*chainProcess) HandleEvent(_ context.Context, _ dogma.StatelessProcessRoot, s dogma.ProcessEventScope[dogma.StatelessProcessRoot], _ dogma.Event) error {
	s.ExecuteCommand(&relayed{})
	return nil
}

// chainIntegration handles [relayed] and records [observed].
type chainIntegration struct{}

func (*chainIntegration) Configure(c dogma.IntegrationConfigurer) {
	c.Identity("integration", "2db79e89-5e24-405b-ba5e-9ea1bda9eb43")
	c.Routes(
		dogma.HandlesCommand[*relayed](),
		dogma.RecordsEvent[*observed](),
	)
}

func (*chainIntegration) HandleCommand(_ context.Context, s dogma.IntegrationCommandScope, _ dogma.Command) error {
	s.RecordEvent(&observed{})
	return nil
}

// silentApp routes [command] to [silentHandler], which handles it without
// recording any events. Used to drive the ErrEventObserverNotSatisfied path.
type silentApp struct{}

func (*silentApp) Configure(c dogma.ApplicationConfigurer) {
	c.Identity("app", "c08cf015-17f4-48c0-9020-0b2bc2d174c5")
	c.Routes(dogma.ViaIntegration(&silentHandler{}))
}

type silentHandler struct{}

func (*silentHandler) Configure(c dogma.IntegrationConfigurer) {
	c.Identity("handler", "a12952fe-386c-460a-ab55-547edc08a51b")
	c.Routes(dogma.HandlesCommand[*command]())
}

func (*silentHandler) HandleCommand(context.Context, dogma.IntegrationCommandScope, dogma.Command) error {
	return nil
}
