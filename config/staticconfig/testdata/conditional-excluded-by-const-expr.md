# Conditional with constant expression that excludes configuration

This test verifies that the static analyzer excludes information about an
entity's identity if it appears in an unreachable branch.

```au:output au:group=matrix
invalid application github.com/dogmatiq/enginekit/config/staticconfig.App (value unavailable)
  - no identity
```

## After conditional return

```go au:input au:group=matrix
package app

import "github.com/dogmatiq/dogma"

type App struct {}

func (App) Configure(c dogma.ApplicationConfigurer) {
	if true {
        return
	}

	c.Identity("app", "de142370-93ee-409c-9336-5084d9c5e285")
}
```

## Within conditional block

```go au:input au:group=matrix
package app

import "github.com/dogmatiq/dogma"

type App struct {}

func (App) Configure(c dogma.ApplicationConfigurer) {
	if false {
		c.Identity("app", "de142370-93ee-409c-9336-5084d9c5e285")
	}
}
```

## In defer that is never scheduled

```go au:input au:group=matrix
package app

import "github.com/dogmatiq/dogma"

type App struct {}

func (App) Configure(c dogma.ApplicationConfigurer) {
	panic("prevent defer")
	defer c.Identity("app", "de142370-93ee-409c-9336-5084d9c5e285")
}
```
