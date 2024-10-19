package staticconfig

import (
	"fmt"
	"go/types"

	"golang.org/x/tools/go/ssa"
)

type context struct {
	Program  *ssa.Program
	Packages []*ssa.Package

	Dogma struct {
		Package               *ssa.Package
		Application           *types.Interface
		ApplicationConfigurer *types.Interface
	}

	Analysis *Analysis
}

// findDogma updates ctx with information about the Dogma package.
//
// It returns false if the Dogma package has not been imported.
func findDogma(ctx *context) bool {
	for _, pkg := range ctx.Program.AllPackages() {
		if pkg.Pkg.Path() != "github.com/dogmatiq/dogma" {
			continue
		}

		iface := func(n string) *types.Interface {
			return pkg.Pkg.
				Scope().
				Lookup(n).
				Type().
				Underlying().(*types.Interface)
		}

		ctx.Dogma.Package = pkg
		ctx.Dogma.Application = iface("Application")
		ctx.Dogma.ApplicationConfigurer = iface("ApplicationConfigurer")

		return true
	}

	return false
}

func (c *context) LookupMethod(t types.Type, name string) *ssa.Function {
	fn := c.Program.LookupMethod(t, packageOf(t), name)
	if fn == nil {
		panic(fmt.Sprintf("method not found: %s.%s", t, name))
	}
	return fn
}
