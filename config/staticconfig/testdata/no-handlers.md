# Applications with no handlers

This test ensures that the static analyzer includes Dogma applications that have
no handlers.

```au:output au:group=matrix
valid application github.com/dogmatiq/enginekit/config/staticconfig.App (value unavailable)
  - valid identity app/8a6baab1-ee64-402e-a081-e43f4bebc243
```

```go au:input au:group=matrix
package app

import "github.com/dogmatiq/dogma"

type App struct{}

func (App) Configure(c dogma.ApplicationConfigurer) {
	c.Identity("app", "8a6baab1-ee64-402e-a081-e43f4bebc243")
}
```
