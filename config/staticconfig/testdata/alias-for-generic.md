# Iinstatiated generic application

This test ensures that the static analyzer finds an instantiated generic type
that implements the `dogma.Application` interface.

```au:output au:group=matrix
valid application github.com/dogmatiq/enginekit/config/staticconfig/testdata.Alias (runtime type unavailable)
  - valid identity app/8a6baab1-ee64-402e-a081-e43f4bebc243
```

```go au:input au:group=matrix
package app

import "github.com/dogmatiq/dogma"

type App[T any] struct{}

func (App[T]) Configure(c dogma.ApplicationConfigurer) {
    c.Identity("app", "8a6baab1-ee64-402e-a081-e43f4bebc243")
}

type Alias = App[int]
```
