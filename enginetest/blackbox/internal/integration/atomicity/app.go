package atomicity

import (
	"context"

	"github.com/dogmatiq/dogma"
)

// app routes [command] to [handler].
type app struct{}

func (*app) Configure(c dogma.ApplicationConfigurer) {
	c.Identity("app", "06d0b3a6-5237-416b-b26b-f11c754c941f")
	c.Routes(dogma.ViaIntegration(&handler{}))
}

// handler handles [command] and records [started] and [completed] in the same
// scope, exercising the atomicity guarantee.
type handler struct{}

func (*handler) Configure(c dogma.IntegrationConfigurer) {
	c.Identity("handler", "17c4ecf3-e588-46cf-8dc1-72ece11722ed")
	c.Routes(
		dogma.HandlesCommand[*command](),
		dogma.RecordsEvent[*started](),
		dogma.RecordsEvent[*completed](),
	)
}

func (*handler) HandleCommand(_ context.Context, s dogma.IntegrationCommandScope, _ dogma.Command) error {
	s.RecordEvent(&started{})
	s.RecordEvent(&completed{})
	return nil
}
