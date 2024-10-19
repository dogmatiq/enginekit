# No applications

This test ensures that the static analyzer does not fail when the analyzed code
does not contain any Dogma applications.

```au:output
(no applications found)
```

## Empty package

```go au:input
package app
```

## Concrete type with similar structure to a dogma.Application

```go au:input
package app

import _ "github.com/dogmatiq/dogma"

// App looks a lot like a [dogma.Application], but does not actually
// implement the Dogma interface because the local [ApplicationConfigurer]
// interface is not the same type as [dogma.ApplicationConfigurer], even though
// it's compatible.
type App struct{}

func (a App) Configure(c ApplicationConfigurer) {
    c.Identity("name", "ee6ca834-34a3-4e59-8c36-7aeb796401d7")
}

type ApplicationConfigurer interface {
    Identity(name, key string)
}
```

## Interface that is compatible with dogma.Application

```go au:input
package app

import "github.com/dogmatiq/dogma"

// App does implement [dogma.Application], but it is not a concrete type so
// there's nothing to analyze.
type App interface {
    Configure(c dogma.ApplicationConfigurer)
}
```

## Uninstantiated generic application

We can't analyze this code because the application is generic and not
instantiated, meaning that we have no concrete type for `T`. We _could_ chose a
compatible type for `T` and analyze the result of instantiating the generic
type, but the assumption is that the reason the type is generic is because the
application is intended to be used with multiple types.

```go au:input
package app

import "github.com/dogmatiq/dogma"

type App[T any] struct{}

func (App[T]) Configure(c dogma.ApplicationConfigurer) {
    c.Identity("app", "8a6baab1-ee64-402e-a081-e43f4bebc243")
}
```
