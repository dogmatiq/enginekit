# Multiple applications in a single package

This test verifies that the static analyzer discovers multiple Dogma application
types defined within the same Go package.

```au:output
valid application github.com/dogmatiq/enginekit/config/staticconfig/testdata/pkg.One (runtime type unavailable)
  - valid identity one/4fec74a1-6ed4-46f4-8417-01e0910be8f1

valid application github.com/dogmatiq/enginekit/config/staticconfig/testdata/pkg.Two (runtime type unavailable)
  - valid identity two/6e97d403-3cb8-4a59-a7ec-74e8e219a7bc
```

```go au:input
package app

import "github.com/dogmatiq/dogma"

type (
	One struct{}
	Two struct{}
)

func (One) Configure(c dogma.ApplicationConfigurer) {
	c.Identity("one", "4fec74a1-6ed4-46f4-8417-01e0910be8f1")
}

func (Two) Configure(c dogma.ApplicationConfigurer) {
	c.Identity("two", "6e97d403-3cb8-4a59-a7ec-74e8e219a7bc")
}
```
