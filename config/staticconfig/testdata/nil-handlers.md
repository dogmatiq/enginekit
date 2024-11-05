# Nil handlers

This test ensures that the static analyzer includes basic information about the
presence of `nil` handlers.

```au:output au:group=matrix
valid application github.com/dogmatiq/enginekit/config/staticconfig.App (value unavailable)
  - valid identity app/0726ae0d-67e4-4a50-8a19-9f58eae38e51
  - incomplete aggregate
      - could not evaluate entire configuration: the handler's type is unknown
      - no identity
      - no handles-command routes
      - no records-event routes
  - incomplete process
      - could not evaluate entire configuration: the handler's type is unknown
      - no identity
      - no handles-event routes
      - no executes-command routes
  - incomplete integration
      - could not evaluate entire configuration: the handler's type is unknown
      - no identity
      - no handles-command routes
  - incomplete projection
      - could not evaluate entire configuration: the handler's type is unknown
      - no identity
      - no handles-event routes
```

```go au:input au:group=matrix
package app

import "github.com/dogmatiq/dogma"

type App struct{}

func (App) Configure(c dogma.ApplicationConfigurer) {
	c.Identity("app", "0726ae0d-67e4-4a50-8a19-9f58eae38e51")

	c.RegisterAggregate(nil)
	c.RegisterProcess(nil)
	c.RegisterIntegration(nil)
	c.RegisterProjection(nil)
}
```
