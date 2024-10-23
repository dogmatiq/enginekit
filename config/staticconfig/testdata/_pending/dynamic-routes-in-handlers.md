# Dynamic routes inside Dogma Application Handlers

This test verifies that static analysis correctly parses routes in handles that
are dynamically populated.

```go au:input au:group=matrix
package app

import (
	"context"
	. "github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/enginetest/stubs"
)

// App implements Application interface.
type App struct{}

// Configure configures the behavior of the engine as it relates to this
// application.
func (App) Configure(c ApplicationConfigurer) {
	c.Identity("<app>", "3bc3849b-abe0-4c4e-9db4-e48dc28c9a26")
	c.RegisterIntegration(IntegrationHandler{})
}

// IntegrationHandler is a test implementation of
// IntegrationMessageHandler.
type IntegrationHandler struct{}

// Configure configures the behavior of the engine as it relates to this
// handler.
func (IntegrationHandler) Configure(c IntegrationConfigurer) {
	c.Identity("<integration>", "3a06b7da-1079-4e4b-a6a6-064c62241918")

	routes := []IntegrationRoute{
		HandlesCommand[stubs.CommandStub[stubs.TypeA]](),
		RecordsEvent[stubs.EventStub[stubs.TypeA]](),
	}

	c.Routes(routes...)
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

```au:output au:group=matrix
application <app> (3bc3849b-abe0-4c4e-9db4-e48dc28c9a26) App

    - integration <integration> (3a06b7da-1079-4e4b-a6a6-064c62241918) IntegrationHandler
        handles CommandStub[TypeA]?
        records EventStub[TypeA]!
```
