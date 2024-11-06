# Indirect configuration

This test verifies that the static analyzer traverses into code called from the
`Configure()` method if that method is given access to the
`ApplicationConfigurer` interface.

```au:output au:group=matrix
valid application github.com/dogmatiq/enginekit/config/staticconfig.App (value unavailable)
  - valid identity app/de142370-93ee-409c-9336-5084d9c5e285
```

## Method call

```go au:input au:group=matrix
package app

import "github.com/dogmatiq/dogma"

type App struct {}

func (a App) Configure(c dogma.ApplicationConfigurer) {
    a.setup(c)
}

func (App) setup(c dogma.ApplicationConfigurer) {
    c.Identity("app", "de142370-93ee-409c-9336-5084d9c5e285")
}
```

## Function call

```go au:input au:group=matrix
package app

import "github.com/dogmatiq/dogma"

type App struct {}

func (App) Configure(c dogma.ApplicationConfigurer) {
    setup(c)
}

func setup(c dogma.ApplicationConfigurer) {
    c.Identity("app", "de142370-93ee-409c-9336-5084d9c5e285")
}
```

## Generic function call

```go au:input au:group=matrix
package app

import "github.com/dogmatiq/dogma"

type App struct {}

func (App) Configure(c dogma.ApplicationConfigurer) {
    setup(c)
}

func setup[T dogma.ApplicationConfigurer](c T) {
    c.Identity("app", "de142370-93ee-409c-9336-5084d9c5e285")
}
```

## Deferred

```go au:input au:group=matrix
package app

import "github.com/dogmatiq/dogma"

type App struct {}

func (App) Configure(c dogma.ApplicationConfigurer) {
    defer c.Identity("app", "de142370-93ee-409c-9336-5084d9c5e285")
}
```

## Separate goroutine

This test guarantees that the identity configured in a separate goroutine is
detected by the static analyzer, but this usage would like introduce a race
condition in any real `ApplicationConfigurer` implementation.

```go au:input au:group=matrix
package app

import "github.com/dogmatiq/dogma"

type App struct {}

func (App) Configure(c dogma.ApplicationConfigurer) {
    go c.Identity("app", "de142370-93ee-409c-9336-5084d9c5e285")
}
```
