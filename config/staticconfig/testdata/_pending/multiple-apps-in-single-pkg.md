# Multiple Dogma Apps in a Single Package

This test verifies that static analysis correctly parses multiple Dogma
applications in a single Go package.

```go au:input
package app

import (
	. "github.com/dogmatiq/dogma"
)

// AppFirst implements Application interface.
type AppFirst struct{}

// Configure configures the behavior of the engine as it relates to this
// application.
func (AppFirst) Configure(c ApplicationConfigurer) {
	c.Identity("<app-first>", "4fec74a1-6ed4-46f4-8417-01e0910be8f1")
}

// AppSecond implements Application interface.
type AppSecond struct{}

// Configure configures the behavior of the engine as it relates to this
// application.
func (a AppSecond) Configure(c ApplicationConfigurer) {
	c.Identity("<app-second>", "6e97d403-3cb8-4a59-a7ec-74e8e219a7bc")
}
```

```au:output
application <app-first> (4fec74a1-6ed4-46f4-8417-01e0910be8f1) AppFirst

application <app-second> (6e97d403-3cb8-4a59-a7ec-74e8e219a7bc) AppSecond
```
