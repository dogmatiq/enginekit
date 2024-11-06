# Generic handler

This test ensures that the static analyzer can analyze handlers that are
implemented using generic types.

```au:output au:group=matrix
valid application github.com/dogmatiq/enginekit/config/staticconfig.App (value unavailable)
  - valid identity app/3109677f-5ed5-4a30-86a1-9975273c5a38
  - valid integration github.com/dogmatiq/enginekit/config/staticconfig.Integration[github.com/dogmatiq/enginekit/enginetest/stubs.CommandStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA]] (value unavailable)
      - valid identity handler/40393d25-f95a-46ea-8702-068643c20ed6
      - valid handles-command route for github.com/dogmatiq/enginekit/enginetest/stubs.CommandStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA] (type unavailable)
```

```go au:input au:group=matrix
package app

import "context"
import "github.com/dogmatiq/dogma"
import "github.com/dogmatiq/enginekit/enginetest/stubs"

type Integration[T dogma.Command] struct{}

func (Integration[T]) Configure(c dogma.IntegrationConfigurer) {
	c.Identity("handler", "40393d25-f95a-46ea-8702-068643c20ed6")
	c.Routes(dogma.HandlesCommand[T]())
}

func (Integration[T]) HandleCommand(context.Context, dogma.IntegrationCommandScope, dogma.Command) error { return nil }

type App struct{}

func (App) Configure(c dogma.ApplicationConfigurer) {
    c.Identity("app", "3109677f-5ed5-4a30-86a1-9975273c5a38")
    c.RegisterIntegration(Integration[stubs.CommandStub[stubs.TypeA]]{})
}
```
