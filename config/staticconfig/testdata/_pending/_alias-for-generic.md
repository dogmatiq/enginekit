# Uninstatiated generic application

This test ensures that the static analyzer does not fail if it encounters a
generic type that implements the `dogma.Application` interface.

```au:output
(no applications found)
```

```go au:input
package app

import "github.com/dogmatiq/dogma"

type App[T any] struct{}

func (App[T]) Configure(c dogma.ApplicationConfigurer) {
    c.Identity("app", "8a6baab1-ee64-402e-a081-e43f4bebc243")
}

type Alias = App[int]
```
