# Incomplete entity configuration

This test verifies that the static analyzer marks configuration of an entity as
incomplete if the `Configure()` method calls into code that is unable to be
analyzed.

```au:output au:group=matrix
incomplete application github.com/dogmatiq/enginekit/config/staticconfig.App (runtime type unavailable)
  - valid identity app/de142370-93ee-409c-9336-5084d9c5e285
```

## Call to closure

```go au:input au:group=matrix
package app

import "github.com/dogmatiq/dogma"

type App struct {
    setup func(dogma.ApplicationConfigurer)
}

func (a App) Configure(c dogma.ApplicationConfigurer) {
    c.Identity("app", "de142370-93ee-409c-9336-5084d9c5e285")
    a.setup(c)
}
```

## Call to method on interface

```go au:input au:group=matrix
package app

import "github.com/dogmatiq/dogma"

type App struct {
    iface interface {
      setup(dogma.ApplicationConfigurer)
    }
}

func (a App) Configure(c dogma.ApplicationConfigurer) {
    c.Identity("app", "de142370-93ee-409c-9336-5084d9c5e285")
    a.iface.setup(c)
}
```
