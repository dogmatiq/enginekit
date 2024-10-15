# Conditional Routes in Dogma Application Handlers

This test verifies that static analysis correctly parses handles that have
conditional routes within their bodies.

```go au:input
package app

import (
	"context"
    "math/rand"
	. "github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/enginetest/stubs"
)

// App implements Application interface.
type App struct{}

// Configure configures the behavior of the engine as it relates to this
// application.
func (App) Configure(c ApplicationConfigurer) {
	c.Identity("<app>", "7e34538e-c407-4af8-8d3c-960e09cde98a")
	c.RegisterIntegration(IntegrationHandler{})
}

// IntegrationHandler is a test implementation of
// IntegrationMessageHandler.
type IntegrationHandler struct{}

// Configure configures the behavior of the engine as it relates to this
// handler.
func (IntegrationHandler) Configure(c IntegrationConfigurer) {
	c.Identity("<integration>", "92cce461-8d30-409b-8d5a-406f656cef2d")

	if rand.Int() == 0 {
		c.Routes(
			HandlesCommand[stubs.CommandStub[stubs.TypeA]](),
			RecordsEvent[stubs.EventStub[stubs.TypeA]](),
		)
	} else {
		c.Routes(
			HandlesCommand[stubs.CommandStub[stubs.TypeB]](),
			RecordsEvent[stubs.EventStub[stubs.TypeB]](),
		)
	}
}

// HandleCommand handles a command message that has been routed to this handler.
func (IntegrationHandler) HandleCommand(
	context.Context,
	IntegrationCommandScope,
	Command,
) error {
	return nil
}

```

```au:output
application <app> (7e34538e-c407-4af8-8d3c-960e09cde98a) App

    - integration <integration> (92cce461-8d30-409b-8d5a-406f656cef2d) IntegrationHandler
        handles CommandStub[TypeA]?
        handles CommandStub[TypeB]?
        records EventStub[TypeA]!
        records EventStub[TypeB]!
```
