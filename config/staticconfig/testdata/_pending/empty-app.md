# Empty application

This test ensures that the static analyzer includes Dogma applications that have
no handlers.

```go au:input
package app

import "github.com/dogmatiq/dogma"

type App struct{}

func (App) Configure(c dogma.ApplicationConfigurer) {
	c.Identity("<app>", "8a6baab1-ee64-402e-a081-e43f4bebc243")
}
```

```au:output
application <app> (8a6baab1-ee64-402e-a081-e43f4bebc243) App
```
