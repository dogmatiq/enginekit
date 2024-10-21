# Type Aliased Handlers

This test verifies that static analysis can correctly parse handlers that are
declared as type aliases.

```go au:input au:group=matrix
package app

import (
	"context"
	. "github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/enginetest/stubs"
)

type (
	// IntegrationHandlerAlias is a test type alias of IntegrationHandler.
	IntegrationHandlerAlias = IntegrationHandler
)

// App implements Application interface.
type App struct{}

// Configure configures the behavior of the engine as it relates to this
// application.
func (App) Configure(c ApplicationConfigurer) {
	c.Identity("<handler-as-typealias>", "1b828a1c-eba1-4e4c-88b8-e49f78ad15c7")

	c.RegisterIntegration(IntegrationHandlerAlias{})
}

// IntegrationHandler is a test implementation of
// IntegrationMessageHandler.
type IntegrationHandler struct{}

// Configure configures the behavior of the engine as it relates to this
// handler.
func (IntegrationHandler) Configure(c IntegrationConfigurer) {
	c.Identity("<integration>", "4d8cd3f5-21dc-475b-a8dc-80138adde3f2")

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

```au:output au:group=matrix
application <handler-as-typealias> (1b828a1c-eba1-4e4c-88b8-e49f78ad15c7) App

    - integration <integration> (4d8cd3f5-21dc-475b-a8dc-80138adde3f2) IntegrationHandlerAlias
        handles CommandStub[TypeB]?
```
