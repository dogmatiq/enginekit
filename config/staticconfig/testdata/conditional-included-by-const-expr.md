# Conditional with constant expression that includes configuration

This test verifies that the static analyzer includes information about an
entity's identity if it appears in a conditional block that is always executed.
Note that the identity is not marked as "speculative".

```au:output au:group=matrix
valid application github.com/dogmatiq/enginekit/config/staticconfig.App (runtime type unavailable)
  - valid identity app/de142370-93ee-409c-9336-5084d9c5e285
```

## After conditional return

```go au:input au:group=matrix
package app

import "github.com/dogmatiq/dogma"

type App struct {}

func (a App) Configure(c dogma.ApplicationConfigurer) {
	if false {
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

func (a App) Configure(c dogma.ApplicationConfigurer) {
	if true {
		c.Identity("app", "de142370-93ee-409c-9336-5084d9c5e285")
	}
}
```

## In defer that is scheduled conditionally

```go au:input au:group=matrix
package app

import "github.com/dogmatiq/dogma"

type App struct {}

func (a App) Configure(c dogma.ApplicationConfigurer) {
	if true {
		defer c.Identity("app", "de142370-93ee-409c-9336-5084d9c5e285")
	}
}
```

## If statement with non-const static condition

```go au:input au:group=matrix
package app

import "github.com/dogmatiq/dogma"

type App struct {}

func (a App) Configure(c dogma.ApplicationConfigurer) {
	if cond() {
		c.Identity("app", "de142370-93ee-409c-9336-5084d9c5e285")
	}
}

func cond() bool {
	return true
}
```