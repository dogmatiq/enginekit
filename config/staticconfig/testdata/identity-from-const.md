# Identity built from constants

This test verifies that the static analyzer can discover the values within an
entity's identity when they are sourced from non-literal constant expressions.

```au:output
valid application github.com/dogmatiq/enginekit/config/staticconfig/testdata/pkg.App (runtime type unavailable)
  - valid identity app/d0de39ba-aaaf-43fd-8e8f-7c4e3be309ec
```

```go au:input
package app

import "github.com/dogmatiq/dogma"

const (
	Name = "app"
	Key = "d0de39ba-aaaf-43fd-8e8f-7c4e3be309ec"
)

type App struct{}

func (App) Configure(c dogma.ApplicationConfigurer) {
	c.Identity(Name, Key)
}
```
