# Identity built from variables

This test verifies that the static analyzer includes an entity's identity, even
if it cannot determine the values used.

```au:output
valid application github.com/dogmatiq/enginekit/config/staticconfig/testdata.App (runtime type unavailable)
  - incomplete identity ?/?
```

```go au:input
package app

import "github.com/dogmatiq/dogma"

type App struct {
	Name string
	Key  string
}

func (a App) Configure(c dogma.ApplicationConfigurer) {
	c.Identity(a.Name, a.Name)
}
```
