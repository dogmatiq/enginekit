# Invalid syntax

This test verifies that static analyzer does not fail catastrophically when the
analyzed code does not compile.

```au:output au:group=matrix
main.go:10:1: expected declaration, found '<'
```

```go au:input au:group=matrix
package app

// Even though this file has invalid syntax the import statements are still
// parsed. This import necessary so that the test still considers it a
// possibility that this package has valid Dogma application implementations.
import _ "github.com/dogmatiq/dogma"

// Below is the deliberate illegal Go syntax to test loading of the packages
// with errors.
<This is the deliberate illegal Go syntax>

```
