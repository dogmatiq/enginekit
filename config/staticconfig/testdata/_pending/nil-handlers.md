# Nil-valued handlers

This test ensures that the static analyzer ignores handlers that are `nil`, but
still includes the application itself.

```go au:input
package app

import "github.com/dogmatiq/dogma"

type App struct{}

func (App) Configure(c dogma.ApplicationConfigurer) {
	c.Identity("<app>", "0726ae0d-67e4-4a50-8a19-9f58eae38e51")

	c.RegisterAggregate(nil)
	c.RegisterProcess(nil)
	c.RegisterProjection(nil)
	c.RegisterIntegration(nil)
}
```

```au:output
application <app> (0726ae0d-67e4-4a50-8a19-9f58eae38e51) App
```
