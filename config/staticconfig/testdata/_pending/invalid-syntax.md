# Type Aliased Handlers

This test verifies that static analysis panics when it encounters incorrect
syntax.

```go au:input au:group=matrix
package app

// Even though this file has invalid syntax the import statements are still
// parsed. This import necessary so that the test still considers it a
// possibility that this package has valid Dogma application implementations.
import "github.com/dogmatiq/dogma"

// Below is the deliberate illegal Go syntax to test loading of the packages
// with errors.
<This is the deliberate illegal Go syntax>

```

```au:output au:group=matrix
expected declaration, found '<'
```
