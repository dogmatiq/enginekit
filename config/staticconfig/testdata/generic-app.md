# Generic application

This test ensures that the static analyzer finds an instantiated generic type
that implements the `dogma.Application` interface.

```au:output au:group=matrix
valid application github.com/dogmatiq/enginekit/config/staticconfig.App (value unavailable)
  - valid identity app/8a6baab1-ee64-402e-a081-e43f4bebc243
```

## Instatiated using a type alias

```go au:input au:group=matrix
package app

import "github.com/dogmatiq/dogma"

type App = AppImpl[int]

type AppImpl[T any] struct{}

func (AppImpl[T]) Configure(c dogma.ApplicationConfigurer) {
    c.Identity("app", "8a6baab1-ee64-402e-a081-e43f4bebc243")
}
```

## Instatiated by embedding in a named struct

```go au:input au:group=matrix
package app

import "github.com/dogmatiq/dogma"

type App struct {
    AppImpl[int]
}

type AppImpl[T any] struct{}

func (AppImpl[T]) Configure(c dogma.ApplicationConfigurer) {
    c.Identity("app", "8a6baab1-ee64-402e-a081-e43f4bebc243")
}
```
