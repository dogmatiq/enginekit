# No applications

This test verifies that the identity specified with constants is correctly
parsed.

```go au:input
package app

import . "github.com/dogmatiq/dogma"

const (
	// AppName is the application name.
	AppName = "<app>"
	// AppKey is the application key.
	AppKey = "04e12cf2-3c66-4414-9203-e045ddbe02c7"
)

// App implements Application interface.
type App struct{}

// Configure sets the application identity using non-literal constant
// expressions.
func (App) Configure(c ApplicationConfigurer) {
	c.Identity(AppName, AppKey)
}
```

```au:output
application <app> (04e12cf2-3c66-4414-9203-e045ddbe02c7) App
```
