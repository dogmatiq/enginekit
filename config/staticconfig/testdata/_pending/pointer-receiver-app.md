# Pointer-Receiver Dogma App

This test verifies that static analysis correctly parses Dogma applications
declared with pointer-receiver methods.

```go au:input
package app

import (
	. "github.com/dogmatiq/dogma"
)

// App implements Application interface.
type App struct{}

// Configure is implemented using a pointer receiver, such that the *App
// implements Application, and not App itself.
func (a *App) Configure(c ApplicationConfigurer) {
	c.Identity("<app>", "b754902b-47c8-48fc-84d2-d920c9cbdaec")
}
```

```au:output
application <app> (b754902b-47c8-48fc-84d2-d920c9cbdaec) *App
```
