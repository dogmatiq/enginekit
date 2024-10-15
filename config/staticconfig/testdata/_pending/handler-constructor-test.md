# Handler constructors

This test verifies that static analysis correctly parses handler constructors.

```go au:input
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
	c.Identity("<handler-constructor>", "3bc3849b-abe0-4c4e-9db4-e48dc28c9a26")

	c.RegisterIntegration(NewIntegrationHandler())
}

// IntegrationHandler is a test implementation of
// IntegrationMessageHandler.
type IntegrationHandler struct{}

// NewIntegrationHandler returns a new IntegrationHandler.
func NewIntegrationHandler() IntegrationHandler {
	panic("the implementation of this function is irrelevant to the analyzer")
}

// Configure configures the behavior of the engine as it relates to this
// handler.
func (IntegrationHandler) Configure(c IntegrationConfigurer) {
	c.Identity("<integration>", "099b5b8d-9e04-422f-bcc3-bb0d451158c7")

	c.Routes(
		HandlesCommand[stubs.CommandStub[stubs.TypeB]](),
	)
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
application <handler-constructor> (3bc3849b-abe0-4c4e-9db4-e48dc28c9a26) App

    - integration <integration> (099b5b8d-9e04-422f-bcc3-bb0d451158c7) IntegrationHandler
        handles CommandStub[TypeB]?
```
