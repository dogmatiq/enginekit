# Handler Adaptors

This test verifies that static analysis correctly parses handler adaptors.

```go au:input
package app

import (
	. "github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/enginetest/stubs"
)

// App implements Application interface.
type App struct{}

// Configure configures the behavior of the engine as it relates to this
// application.
func (App) Configure(c ApplicationConfigurer) {
	c.Identity("<adaptor-func>", "f610eae4-f5d0-4eea-a9c9-6cbbfa9b2060")

	c.RegisterIntegration(AdaptIntegration(IntegrationHandler{}))
}

// IntegrationHandler is the type that provides the handler logic, but is not
// itself an implementation of IntegrationMessageHandler.
type IntegrationHandler struct{}

// Configure configures the behavior of the engine as it relates to this
// handler.
func (IntegrationHandler) Configure(c IntegrationConfigurer) {
	c.Identity("<integration>", "099b5b8d-9e04-422f-bcc3-bb0d451158c7")

	c.Routes(
		HandlesCommand[stubs.CommandStub[stubs.TypeI]](),
		RecordsEvent[stubs.EventStub[stubs.TypeI]](),
	)
}

// PartialIntegrationMessageHandler is the subset of
// IntegrationMessageHandler that must be implemented for a type to be
// detected as a concrete implementation.
type PartialIntegrationMessageHandler interface {
	Configure(c IntegrationConfigurer)
}

// AdaptIntegration adapts the argument to the IntegrationMessageHandler interface.
func AdaptIntegration(PartialIntegrationMessageHandler) IntegrationMessageHandler {
	panic("the implementation of this function is irrelevant to the analyzer")
}
```

```au:output
application <adaptor-func> (f610eae4-f5d0-4eea-a9c9-6cbbfa9b2060) App

    - integration <integration> (099b5b8d-9e04-422f-bcc3-bb0d451158c7) IntegrationHandler
        handles CommandStub[TypeI]?
        records EventStub[TypeI]!
```
