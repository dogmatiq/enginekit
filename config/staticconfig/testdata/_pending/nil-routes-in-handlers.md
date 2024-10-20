# Nil Routes in Dogma Application Handlers

This test verifies that static analysis correctly parses `nil` routes inside
Dogma Application handlers.

```go au:input au:group=matrix
package app

import (
	"context"
	. "github.com/dogmatiq/dogma"
)

// App implements Application interface.
type App struct{}

// Configure configures the behavior of the engine as it relates to this
// application.
func (App) Configure(c ApplicationConfigurer) {
	c.Identity("<app>", "c100edcc-6dcc-42ed-ac75-69eecb3d0ec4")
	c.RegisterIntegration(IntegrationHandler{})
}


// IntegrationHandler is a test implementation of
// IntegrationMessageHandler.
type IntegrationHandler struct{}

// Configure configures the behavior of the engine as it relates to this
// handler.
func (IntegrationHandler) Configure(c IntegrationConfigurer) {
	c.Identity("<integration>", "363039e5-2938-4b2c-9bec-dcb29dee2da1")
	c.Routes(nil)
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
application <app> (c100edcc-6dcc-42ed-ac75-69eecb3d0ec4) App

    - integration <integration> (363039e5-2938-4b2c-9bec-dcb29dee2da1) IntegrationHandler
```
