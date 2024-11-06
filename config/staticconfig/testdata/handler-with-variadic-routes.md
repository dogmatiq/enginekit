# Handler with variadic routes

This test verifies that the static analyzer correctly identifies routes that are
configured via a slice which is used as the variadic parameter to the `Routes()`
method.

```au:output au:group=matrix
valid application github.com/dogmatiq/enginekit/config/staticconfig.App (value unavailable)
  - valid identity app/e4183527-c234-42d5-8709-3dc8b9d5caa4
  - valid integration github.com/dogmatiq/enginekit/config/staticconfig.Integration (value unavailable)
      - valid identity handler/5752bb84-0b65-4a7f-b2fa-bfb77a53a97f
      - valid handles-command route for github.com/dogmatiq/enginekit/enginetest/stubs.CommandStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA] (type unavailable)
      - valid records-event route for github.com/dogmatiq/enginekit/enginetest/stubs.EventStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA] (type unavailable)
```

## Appended

```go au:input au:group=matrix
package app

import "context"
import "github.com/dogmatiq/dogma"
import "github.com/dogmatiq/enginekit/enginetest/stubs"

type Integration struct{}

func (Integration) Configure(c dogma.IntegrationConfigurer) {
	c.Identity("handler", "5752bb84-0b65-4a7f-b2fa-bfb77a53a97f")

	var routes []dogma.IntegrationRoute

	routes = append(
		routes,
		dogma.HandlesCommand[stubs.CommandStub[stubs.TypeA]](),
		dogma.RecordsEvent[stubs.EventStub[stubs.TypeA]](),
	)

	c.Routes(routes...)
}

func (Integration) HandleCommand(context.Context, dogma.IntegrationCommandScope, dogma.Command) error { return nil }

type App struct{}

func (App) Configure(c dogma.ApplicationConfigurer) {
	c.Identity("app", "e4183527-c234-42d5-8709-3dc8b9d5caa4")
	c.RegisterIntegration(Integration{})
}
```

## Assigned to index

```go au:input au:group=matrix
package app

import "context"
import "github.com/dogmatiq/dogma"
import "github.com/dogmatiq/enginekit/enginetest/stubs"

type Integration struct{}

func (Integration) Configure(c dogma.IntegrationConfigurer) {
	c.Identity("handler", "5752bb84-0b65-4a7f-b2fa-bfb77a53a97f")

	routes := make([]dogma.IntegrationRoute, 1)
	routes[0] = dogma.HandlesCommand[stubs.CommandStub[stubs.TypeA]]()
	routes[1] = dogma.RecordsEvent[stubs.EventStub[stubs.TypeA]]()

	c.Routes(routes...)
}

func (Integration) HandleCommand(context.Context, dogma.IntegrationCommandScope, dogma.Command) error { return nil }

type App struct{}

func (App) Configure(c dogma.ApplicationConfigurer) {
	c.Identity("app", "e4183527-c234-42d5-8709-3dc8b9d5caa4")
	c.RegisterIntegration(Integration{})
}
```

## Assigned to index of sub-slice

```go au:input au:group=matrix
package app

import "context"
import "github.com/dogmatiq/dogma"
import "github.com/dogmatiq/enginekit/enginetest/stubs"

type Integration struct{}

func (Integration) Configure(c dogma.IntegrationConfigurer) {
	c.Identity("handler", "5752bb84-0b65-4a7f-b2fa-bfb77a53a97f")

	routes := make([]dogma.IntegrationRoute, 1)
	routes[:1][0] = dogma.HandlesCommand[stubs.CommandStub[stubs.TypeA]]()
	routes[1:][0] = dogma.RecordsEvent[stubs.EventStub[stubs.TypeA]]()

	c.Routes(routes...)
}

func (Integration) HandleCommand(context.Context, dogma.IntegrationCommandScope, dogma.Command) error { return nil }

type App struct{}

func (App) Configure(c dogma.ApplicationConfigurer) {
	c.Identity("app", "e4183527-c234-42d5-8709-3dc8b9d5caa4")
	c.RegisterIntegration(Integration{})
}
```
