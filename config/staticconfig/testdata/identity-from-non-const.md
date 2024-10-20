# Identity built from non-constant values

This test verifies that the static analyzer includes an entity's identity, even
if it cannot determine the values used.

```au:output au:group=matrix
valid application github.com/dogmatiq/enginekit/config/staticconfig/testdata.App (runtime type unavailable)
  - incomplete identity ?/?
```

## Variables

```go au:input au:group=matrix
package app

import "github.com/dogmatiq/dogma"

type App struct {
	Name string
	Key  string
}

func (a App) Configure(c dogma.ApplicationConfigurer) {
	c.Identity(a.Name, a.Key)
}
```

## Function call

```go au:input au:group=matrix
package app

import "github.com/dogmatiq/dogma"

type App struct {
    name func() string
	key func() string
}

func (a App) Configure(c dogma.ApplicationConfigurer) {
    c.Identity(a.name(), a.key())
}
```

## Function call with a non-deterministic return value

```go au:input au:group=matrix
package app

import "math/rand"
import "github.com/dogmatiq/dogma"

type App struct {}

func (a App) Configure(c dogma.ApplicationConfigurer) {
    c.Identity(ident())
}

func ident() (string, string) {
	if rand.Int() == 0 {
		return "app1", "a0a0edb7-ce45-4eb4-940c-0f77459ae2a0"
	}
	return "app2", "08905dce-9059-4601-a48f-f449c6fba70b"
}
```
