# Integration message handler

This test ensures that the static analyzer supports all aspects of configuring
an integration handler.

```au:output au:group=matrix
valid application github.com/dogmatiq/enginekit/config/staticconfig.App (value unavailable)
  - valid identity app/0726ae0d-67e4-4a50-8a19-9f58eae38e51
  - valid integration github.com/dogmatiq/enginekit/config/staticconfig.Integration (value unavailable)
      - valid identity integration/b92431e6-3a7d-4235-a76f-541622c487ee
      - valid handles-command route for github.com/dogmatiq/enginekit/enginetest/stubs.CommandStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA] (type unavailable)
      - valid records-event route for github.com/dogmatiq/enginekit/enginetest/stubs.EventStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA] (type unavailable)
      - valid disabled flag, set to true
```

```go au:input au:group=matrix
package app

import "context"
import "github.com/dogmatiq/dogma"
import "github.com/dogmatiq/enginekit/enginetest/stubs"

type Integration struct {}

func (Integration) Configure(c dogma.IntegrationConfigurer) {
    c.Identity("integration", "b92431e6-3a7d-4235-a76f-541622c487ee")
    c.Routes(
        dogma.HandlesCommand[stubs.CommandStub[stubs.TypeA]](),
        dogma.RecordsEvent[stubs.EventStub[stubs.TypeA]](),
    )
    c.Disable()
}

type App struct{}

func (App) Configure(c dogma.ApplicationConfigurer) {
	c.Identity("app", "0726ae0d-67e4-4a50-8a19-9f58eae38e51")
	c.RegisterIntegration(Integration{})
}

func (Integration) HandleCommand(context.Context, dogma.IntegrationCommandScope, dogma.Command) error { return nil }
```
