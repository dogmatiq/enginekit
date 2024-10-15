# No applications

This test verifies that the identity specified with variables is correctly parsed.

```go au:input
package app

import . "github.com/dogmatiq/dogma"

// App implements Application interface.
type App struct {
	Name string
	Key  string
}

// Configure sets the application identity using non-constant expressions.
func (a App) Configure(c ApplicationConfigurer) {
	c.Identity(a.Name, a.Name)
}
```

```au:output
application  () App
```
