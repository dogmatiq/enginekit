# Non-pointer Handlers Registered in a Dogma Application as Pointers.

This test verifies that static analysis can correctly parse non-pointer handlers
registered in a dogma application as pointers using 'address-of' operator.

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
	c.Identity(
        "<non-pointer-handler-registered-as-pointer>",
        "282653ad-9343-44f1-889e-a8b2b095b54b",
    )

	c.RegisterIntegration(&IntegrationHandler{})
}

// IntegrationHandler is a test implementation of
// IntegrationMessageHandler.
type IntegrationHandler struct{}

// Configure configures the behavior of the engine as it relates to this
// handler.
func (IntegrationHandler) Configure(c IntegrationConfigurer) {
	c.Identity("<integration>", "1425ca64-0448-4bfd-b18d-9fe63a95995f")

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
application <non-pointer-handler-registered-as-pointer> (282653ad-9343-44f1-889e-a8b2b095b54b) App

    - integration <integration> (1425ca64-0448-4bfd-b18d-9fe63a95995f) *IntegrationHandler
        handles CommandStub[TypeB]?
```
