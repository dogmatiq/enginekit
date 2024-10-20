# Unregistered Routes in Dogma Application Handlers

This test verifies that static analysis ignores unregistered routes in Dogma
application handlers.

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
	c.Identity("<app>", "f2c08525-623e-4c76-851c-3172953269e3")
	c.RegisterIntegration(IntegrationHandler{})
}

// IntegrationHandler is a test implementation of
// IntegrationMessageHandler.
type IntegrationHandler struct{}

// Configure configures the behavior of the engine as it relates to this
// handler.
func (IntegrationHandler) Configure(c IntegrationConfigurer) {
	c.Identity("<integration>", "ac391765-da58-4e7c-a478-e4725eb2b0e9")

	// Create a route that is never passed to c.Routes().
	HandlesCommand[stubs.CommandStub[stubs.TypeX]]()

	// Ensure there is still _some_ call to Routes().
	c.Routes(
		HandlesCommand[stubs.CommandStub[stubs.TypeA]](),
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
application <app> (f2c08525-623e-4c76-851c-3172953269e3) App

    - integration <integration> (ac391765-da58-4e7c-a478-e4725eb2b0e9) IntegrationHandler
        handles CommandStub[TypeA]?
```
