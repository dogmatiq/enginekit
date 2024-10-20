# Nil handlers

This test ensures that the static analyzer includes basic information about the
presence of `nil` handlers.

```au:output au:group=matrix
valid application github.com/dogmatiq/enginekit/config/staticconfig/testdata.App (runtime type unavailable)
  - valid identity app/0726ae0d-67e4-4a50-8a19-9f58eae38e51
  - incomplete aggregate
      - no identity is configured
      - no "handles-command" routes are configured
      - no "records-event" routes are configured
  - incomplete process
      - no identity is configured
      - no "handles-event" routes are configured
      - no "executes-command" routes are configured
  - incomplete integration
      - no identity is configured
      - no "handles-command" routes are configured
  - incomplete projection
      - no identity is configured
      - no "handles-event" routes are configured
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
