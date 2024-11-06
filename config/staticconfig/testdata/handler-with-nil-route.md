# Handler with nil route

This test verifies that the static analyzer correctly detects a `nil` route.

```au:output au:group=matrix
valid application github.com/dogmatiq/enginekit/config/staticconfig.App (value unavailable)
  - valid identity app/e4183527-c234-42d5-8709-3dc8b9d5caa4
  - invalid integration github.com/dogmatiq/enginekit/config/staticconfig.Integration (value unavailable)
      - no handles-command routes
      - valid identity handler/5752bb84-0b65-4a7f-b2fa-bfb77a53a97f
      - incomplete route
          - route type is unavailable
          - message type name is unavailable
```

```go au:input au:group=matrix
package app

import "context"
import "github.com/dogmatiq/dogma"

type Integration struct{}

func (Integration) Configure(c dogma.IntegrationConfigurer) {
	c.Identity("handler", "5752bb84-0b65-4a7f-b2fa-bfb77a53a97f")
	c.Routes(nil)
}

func (Integration) HandleCommand(context.Context, dogma.IntegrationCommandScope, dogma.Command) error { return nil }

type App struct{}

func (App) Configure(c dogma.ApplicationConfigurer) {
	c.Identity("app", "e4183527-c234-42d5-8709-3dc8b9d5caa4")
	c.RegisterIntegration(Integration{})
}
```
