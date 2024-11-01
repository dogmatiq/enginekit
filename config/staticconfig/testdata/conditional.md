# Conditional identity

This test verifies that the static analyzer includes information about an
entity's identity when it is defined within a conditional statement.

```au:output au:group=matrix
valid application github.com/dogmatiq/enginekit/config/staticconfig.App (value unavailable)
  - valid speculative identity app/de142370-93ee-409c-9336-5084d9c5e285
```

## If statement

```go au:input au:group=matrix
package app

import "math/rand"
import "github.com/dogmatiq/dogma"

type App struct {}

func (a App) Configure(c dogma.ApplicationConfigurer) {
	if rand.Int() == 0 {
		c.Identity("app", "de142370-93ee-409c-9336-5084d9c5e285")
	}
}
```

## Else statement

```go au:input au:group=matrix
package app

import "math/rand"
import "github.com/dogmatiq/dogma"

type App struct {}

func (a App) Configure(c dogma.ApplicationConfigurer) {
	if rand.Int() == 0 {
	} else {
		c.Identity("app", "de142370-93ee-409c-9336-5084d9c5e285")
	}
}
```

## After conditional return

```go au:input au:group=matrix
package app

import "math/rand"
import "github.com/dogmatiq/dogma"

type App struct {}

func (a App) Configure(c dogma.ApplicationConfigurer) {
	if rand.Int() == 0 {
		return
	}

	c.Identity("app", "de142370-93ee-409c-9336-5084d9c5e285")
}
```

## After conditonal panic

```go au:input au:group=matrix
package app

import "math/rand"
import "github.com/dogmatiq/dogma"

type App struct {}

func (a App) Configure(c dogma.ApplicationConfigurer) {
	if rand.Int() == 0 {
		panic("oh no")
	}

	c.Identity("app", "de142370-93ee-409c-9336-5084d9c5e285")
}
```

## Switch statement

```go au:input au:group=matrix
package app

import "math/rand"
import "github.com/dogmatiq/dogma"

type App struct {}

func (a App) Configure(c dogma.ApplicationConfigurer) {
	switch rand.Int() {
	case 0:
		c.Identity("app", "de142370-93ee-409c-9336-5084d9c5e285")
	}
}
```

## For statement

```go au:input au:group=matrix
package app

import "math/rand"
import "github.com/dogmatiq/dogma"

type App struct {}

func (a App) Configure(c dogma.ApplicationConfigurer) {
	for range rand.Int() {
		c.Identity("app", "de142370-93ee-409c-9336-5084d9c5e285")
	}
}
```

## Select statement

```go au:input au:group=matrix
package app

import "math/rand"
import "time"
import "github.com/dogmatiq/dogma"

type App struct {}

func (a App) Configure(c dogma.ApplicationConfigurer) {
	select {
	case <-time.After(time.Duration(rand.Int())):
		c.Identity("app", "de142370-93ee-409c-9336-5084d9c5e285")
	default:
	}
}
```

## Deferred

```go au:input au:group=matrix
package app

import "math/rand"
import "github.com/dogmatiq/dogma"

type App struct {}

func (a App) Configure(c dogma.ApplicationConfigurer) {
	if rand.Int() == 0 {
		defer c.Identity("app", "de142370-93ee-409c-9336-5084d9c5e285")
	}
}
```
