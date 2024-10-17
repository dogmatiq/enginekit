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
