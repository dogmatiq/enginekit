# No applications

This test verifies that the identity specified with string literals is correctly
parsed.

```go au:input
package app

import . "github.com/dogmatiq/dogma"

// App implements Application interface.
type App struct{}

// Configure sets the application identity using literal string values.
func (App) Configure(c ApplicationConfigurer) {
	c.Identity("<app>", "9d0af85d-f506-4742-b676-ce87730bb1a0")
}
```

```au:output
application <app> (9d0af85d-f506-4742-b676-ce87730bb1a0) App
```
