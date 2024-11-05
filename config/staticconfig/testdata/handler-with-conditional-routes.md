# Handler with conditional routes

This test verifies that the static analyzer correctly identifies when routes are
added to a handler conditionally.

```au:output au:group=matrix
valid application github.com/dogmatiq/enginekit/config/staticconfig.App (value unavailable)
  - valid identity app/e4183527-c234-42d5-8709-3dc8b9d5caa4
  - valid integration github.com/dogmatiq/enginekit/config/staticconfig.Integration (value unavailable)
      - valid identity handler/5752bb84-0b65-4a7f-b2fa-bfb77a53a97f
      - valid speculative handles-command route for github.com/dogmatiq/enginekit/enginetest/stubs.CommandStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA] (type unavailable)
      - valid speculative records-event route for github.com/dogmatiq/enginekit/enginetest/stubs.EventStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA] (type unavailable)
      - valid speculative handles-command route for github.com/dogmatiq/enginekit/enginetest/stubs.CommandStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeB] (type unavailable)
      - valid speculative records-event route for github.com/dogmatiq/enginekit/enginetest/stubs.EventStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeB] (type unavailable)
```

## Routes() call within conditional block

```go au:input au:group=matrix
package app

import "context"
import "math/rand"
import "github.com/dogmatiq/dogma"
import "github.com/dogmatiq/enginekit/enginetest/stubs"

type Integration struct{}

func (Integration) Configure(c dogma.IntegrationConfigurer) {
	c.Identity("handler", "5752bb84-0b65-4a7f-b2fa-bfb77a53a97f")

	if rand.Int() == 0 {
		c.Routes(
			dogma.HandlesCommand[stubs.CommandStub[stubs.TypeA]](),
			dogma.RecordsEvent[stubs.EventStub[stubs.TypeA]](),
		)
	} else {
		c.Routes(
			dogma.HandlesCommand[stubs.CommandStub[stubs.TypeB]](),
			dogma.RecordsEvent[stubs.EventStub[stubs.TypeB]](),
		)
	}
}

func (Integration) HandleCommand(context.Context, dogma.IntegrationCommandScope, dogma.Command) error { return nil }

type App struct{}

func (App) Configure(c dogma.ApplicationConfigurer) {
	c.Identity("app", "e4183527-c234-42d5-8709-3dc8b9d5caa4")
	c.RegisterIntegration(Integration{})
}
```

## Slice built within conditional block

```go au:input au:group=matrix
package app

import "context"
import "math/rand"
import "github.com/dogmatiq/dogma"
import "github.com/dogmatiq/enginekit/enginetest/stubs"

type Integration struct{}

func (Integration) Configure(c dogma.IntegrationConfigurer) {
	c.Identity("handler", "5752bb84-0b65-4a7f-b2fa-bfb77a53a97f")

	var routes []dogma.IntegrationRoute

	if rand.Int() == 0 {
		routes = append(
			routes,
			dogma.HandlesCommand[stubs.CommandStub[stubs.TypeA]](),
			dogma.RecordsEvent[stubs.EventStub[stubs.TypeA]](),
		)
	} else {
		routes = append(
			routes,
			dogma.HandlesCommand[stubs.CommandStub[stubs.TypeB]](),
			dogma.RecordsEvent[stubs.EventStub[stubs.TypeB]](),
		)
	}

	c.Routes(routes...)
}

func (Integration) HandleCommand(context.Context, dogma.IntegrationCommandScope, dogma.Command) error { return nil }

type App struct{}

func (App) Configure(c dogma.ApplicationConfigurer) {
	c.Identity("app", "e4183527-c234-42d5-8709-3dc8b9d5caa4")
	c.RegisterIntegration(Integration{})
}
```

## Slice built within multiple conditional blocks

```go au:input au:group=matrix
package app

import "context"
import "math/rand"
import "github.com/dogmatiq/dogma"
import "github.com/dogmatiq/enginekit/enginetest/stubs"

type Integration struct{}

func (Integration) Configure(c dogma.IntegrationConfigurer) {
	c.Identity("handler", "5752bb84-0b65-4a7f-b2fa-bfb77a53a97f")

	var routes []dogma.IntegrationRoute

	if rand.Int() == 0 {
		routes = append(
			routes,
			dogma.HandlesCommand[stubs.CommandStub[stubs.TypeA]](),
			dogma.RecordsEvent[stubs.EventStub[stubs.TypeA]](),
		)
	}

	if rand.Int() == 0 {
		routes = append(
			routes,
			dogma.HandlesCommand[stubs.CommandStub[stubs.TypeB]](),
			dogma.RecordsEvent[stubs.EventStub[stubs.TypeB]](),
		)
	}

	c.Routes(routes...)
}

func (Integration) HandleCommand(context.Context, dogma.IntegrationCommandScope, dogma.Command) error { return nil }

type App struct{}

func (App) Configure(c dogma.ApplicationConfigurer) {
	c.Identity("app", "e4183527-c234-42d5-8709-3dc8b9d5caa4")
	c.RegisterIntegration(Integration{})
}
```
