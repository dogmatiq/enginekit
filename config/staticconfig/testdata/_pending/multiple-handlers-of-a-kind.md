# Multiple handlers of a same kind

This test verifies that static analysis can correctly parse multiple handlers of
a same kind.

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
	c.Identity("<multiple-handlers-of-a-kind>", "8961f548-1afc-4996-894c-956835c83199")

	c.RegisterIntegration(FirstIntegrationHandler{})
	c.RegisterIntegration(SecondIntegrationHandler{})
}

// FirstIntegrationHandler is a test implementation of
// IntegrationMessageHandler.
type FirstIntegrationHandler struct{}

// Configure configures the behavior of the engine as it relates to this
// handler.
func (FirstIntegrationHandler) Configure(c IntegrationConfigurer) {
	c.Identity("<first-integration>", "14cf2812-eead-43b3-9c9c-10db5b469e94")

	c.Routes(
		HandlesCommand[stubs.CommandStub[stubs.TypeC]](),
	)
}

// RouteCommandToInstance returns the ID of the integration instance that is
// targetted by m.
func (FirstIntegrationHandler) RouteCommandToInstance(Command) string {
	return "<first-integration>"
}

// HandleCommand handles a command message that has been routed to this handler.
func (FirstIntegrationHandler) HandleCommand(
	context.Context,
	IntegrationCommandScope,
	Command,
) error {
	return nil
}

// SecondIntegrationHandler is a test implementation of
// IntegrationMessageHandler.
type SecondIntegrationHandler struct{}

// Configure configures the behavior of the engine as it relates to this
// handler.
func (SecondIntegrationHandler) Configure(c IntegrationConfigurer) {
	c.Identity("<second-integration>", "6bed3fbc-30e2-44c7-9a5b-e440ffe370d9")

	c.Routes(
		HandlesCommand[stubs.CommandStub[stubs.TypeD]](),
	)
}

// RouteCommandToInstance returns the ID of the integration instance that is
// targetted by m.
func (SecondIntegrationHandler) RouteCommandToInstance(Command) string {
	return "<second-integration>"
}

// HandleCommand handles a command message that has been routed to this handler.
func (SecondIntegrationHandler) HandleCommand(
	context.Context,
	IntegrationCommandScope,
	Command,
) error {
	return nil
}
```

```au:output
application <multiple-handlers-of-a-kind> (8961f548-1afc-4996-894c-956835c83199) App

    - integration <first-integration> (14cf2812-eead-43b3-9c9c-10db5b469e94) FirstIntegrationHandler
        handles CommandStub[TypeC]?

    - integration <second-integration> (6bed3fbc-30e2-44c7-9a5b-e440ffe370d9) SecondIntegrationHandler
        handles CommandStub[TypeD]?
```
