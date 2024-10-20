# Interface as an entity configurer.

This test ensures that the static analyzer does not behaves abnormally when it
encounters an interface that handles configuration. In this case, the static
analysis is not capable of gathering data about what particular entity is
configured behind the interface.

```go au:input au:group=matrix
package app

import (
    . "github.com/dogmatiq/dogma"
)


type Configurer interface {
	ApplyConfiguration(c ApplicationConfigurer)
}

type App struct {
	C Configurer
}

func (a App) Configure(c ApplicationConfigurer) {
	c.Identity("<app>", "7468a57f-20f0-4d11-9aad-48fcd553a908")
	a.C.ApplyConfiguration(c)
}

```

```au:output au:group=matrix
application <app> (7468a57f-20f0-4d11-9aad-48fcd553a908) App
```
