# Handler with speculative routes

This test verifies that the static analyzer correctly identifies routes as
"speculative" under various complex conditions.

Some of these scenarios could be improved, potentially avoiding false positives
for the "speculative" flag. In general, however, it is preferred that the
analyzer errs on the side of caution and marks routes as "speculative" when it
is unsure.

```au:output au:group=matrix
valid application github.com/dogmatiq/enginekit/config/staticconfig.App (value unavailable)
  - valid identity app/e4183527-c234-42d5-8709-3dc8b9d5caa4
  - valid integration github.com/dogmatiq/enginekit/config/staticconfig.Integration (value unavailable)
      - valid identity handler/5752bb84-0b65-4a7f-b2fa-bfb77a53a97f
      - valid speculative handles-command route for github.com/dogmatiq/enginekit/enginetest/stubs.CommandStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA] (type unavailable)
      - valid speculative records-event route for github.com/dogmatiq/enginekit/enginetest/stubs.EventStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA] (type unavailable)
```

## Random index

```go au:input au:group=matrix
package app

import "context"
import "math/rand"
import "github.com/dogmatiq/dogma"
import "github.com/dogmatiq/enginekit/enginetest/stubs"

type Integration struct{}

func (Integration) Configure(c dogma.IntegrationConfigurer) {
	c.Identity("handler", "5752bb84-0b65-4a7f-b2fa-bfb77a53a97f")

	routes := make([]dogma.IntegrationRoute, 0)

	routes[0] = dogma.HandlesCommand[stubs.CommandStub[stubs.TypeA]]()
	routes[rand.Int()] = dogma.RecordsEvent[stubs.EventStub[stubs.TypeA]]()

	c.Routes(routes...)
}

func (Integration) HandleCommand(context.Context, dogma.IntegrationCommandScope, dogma.Command) error { return nil }

type App struct{}

func (App) Configure(c dogma.ApplicationConfigurer) {
	c.Identity("app", "e4183527-c234-42d5-8709-3dc8b9d5caa4")
	c.RegisterIntegration(Integration{})
}
```

## Colliding indices

```go au:input au:group=matrix
package app

import "context"
import "github.com/dogmatiq/dogma"
import "github.com/dogmatiq/enginekit/enginetest/stubs"

type Integration struct{}

func (Integration) Configure(c dogma.IntegrationConfigurer) {
	c.Identity("handler", "5752bb84-0b65-4a7f-b2fa-bfb77a53a97f")

	routes := make([]dogma.IntegrationRoute, 0)

	routes[0] = dogma.HandlesCommand[stubs.CommandStub[stubs.TypeA]]()
	routes[0] = dogma.RecordsEvent[stubs.EventStub[stubs.TypeA]]()

	c.Routes(routes...)
}

func (Integration) HandleCommand(context.Context, dogma.IntegrationCommandScope, dogma.Command) error { return nil }

type App struct{}

func (App) Configure(c dogma.ApplicationConfigurer) {
	c.Identity("app", "e4183527-c234-42d5-8709-3dc8b9d5caa4")
	c.RegisterIntegration(Integration{})
}
```
