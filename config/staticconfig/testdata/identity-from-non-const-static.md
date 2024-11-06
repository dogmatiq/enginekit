# Identity built from non-constant values that can be resolved statically

This test verifies that the static analyzer includes an entity's identity, even
if it cannot determine the values used.

```au:output au:group=matrix
valid application github.com/dogmatiq/enginekit/config/staticconfig.App (value unavailable)
  - valid identity app/a0a0edb7-ce45-4eb4-940c-0f77459ae2a0
```

## Function call

```go au:input au:group=matrix
package app

import "github.com/dogmatiq/dogma"

type App struct {
}

func (App) Configure(c dogma.ApplicationConfigurer) {
	c.Identity(name(), "a0a0edb7-ce45-4eb4-940c-0f77459ae2a0")
}

func name() string {
  return "app"
}
```

## Function call with tuple extraction

```go au:input au:group=matrix
package app

import "github.com/dogmatiq/dogma"

type App struct {
}

func (App) Configure(c dogma.ApplicationConfigurer) {
	c.Identity(ident())
}

func ident() (string,string) {
  return "app", "a0a0edb7-ce45-4eb4-940c-0f77459ae2a0"
}
```

## Function call with multiple branches that return the same values

```go au:input au:group=matrix
package app

import "math/rand"
import "github.com/dogmatiq/dogma"

type App struct {}

func (App) Configure(c dogma.ApplicationConfigurer) {
    c.Identity(ident())
}

func ident() (string, string) {
	if rand.Int() == 0 {
		return "app", "a0a0edb7-ce45-4eb4-940c-0f77459ae2a0"
	}
	return "app", "a0a0edb7-ce45-4eb4-940c-0f77459ae2a0"
}
```
