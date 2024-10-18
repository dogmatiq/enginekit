# Conditional identity

This test verifies that the static analyzer excludes information about an
entity's identity even if it appears in an unreachable branch.

```au:output
valid application github.com/dogmatiq/enginekit/config/staticconfig/testdata/pkg.App (runtime type unavailable)
  - no identities configured
```

## If statement

```go au:input
package app

import "github.com/dogmatiq/dogma"

type App struct {}

func (a App) Configure(c dogma.ApplicationConfigurer) {
	if true || false {
        return
	}

	c.Identity("app", "de142370-93ee-409c-9336-5084d9c5e285")
}
```
