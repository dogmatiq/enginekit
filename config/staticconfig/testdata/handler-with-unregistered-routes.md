# Handler with unregistered routes

This test verifies that static analyzer does not include information about
routes that are constructed but never passed to the configurer's `Routes()`
method.

```au:output au:group=matrix
valid application github.com/dogmatiq/enginekit/config/staticconfig.App (value unavailable)
  - valid identity app/f2c08525-623e-4c76-851c-3172953269e3
  - invalid integration github.com/dogmatiq/enginekit/config/staticconfig.Integration (value unavailable)
      - no handles-command routes
      - valid identity handler/ac391765-da58-4e7c-a478-e4725eb2b0e9
```

## No call to Routes()

```go au:input au:group=matrix
package app

import (
	"context"
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/enginetest/stubs"
)

type Integration struct{}

func (Integration) Configure(c dogma.IntegrationConfigurer) {
	c.Identity("handler", "ac391765-da58-4e7c-a478-e4725eb2b0e9")
	dogma.HandlesCommand[stubs.CommandStub[stubs.TypeX]]()
}

type App struct{}

func (App) Configure(c dogma.ApplicationConfigurer) {
	c.Identity("app", "f2c08525-623e-4c76-851c-3172953269e3")
	c.RegisterIntegration(Integration{})
}

func (Integration) HandleCommand(context.Context, dogma.IntegrationCommandScope, dogma.Command) error { return nil }
```

## Route appended to slice after calling Routes()

```go au:input au:group=matrix
package app

import (
	"context"
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/enginetest/stubs"
)

type Integration struct{}

func (Integration) Configure(c dogma.IntegrationConfigurer) {
	c.Identity("handler", "ac391765-da58-4e7c-a478-e4725eb2b0e9")

	var routes []dogma.IntegrationRoute

	c.Routes(routes...)

	routes = append(
		routes,
		dogma.HandlesCommand[stubs.CommandStub[stubs.TypeX]](),
	)
}

type App struct{}

func (App) Configure(c dogma.ApplicationConfigurer) {
	c.Identity("app", "f2c08525-623e-4c76-851c-3172953269e3")


	c.RegisterIntegration(Integration{})
}

func (Integration) HandleCommand(context.Context, dogma.IntegrationCommandScope, dogma.Command) error { return nil }
```
