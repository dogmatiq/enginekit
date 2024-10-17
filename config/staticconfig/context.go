package staticconfig

import (
	"go/types"

	"github.com/dogmatiq/enginekit/config"
	"golang.org/x/tools/go/ssa"
)

type context struct {
	Program  *ssa.Program
	Packages []*ssa.Package

	Dogma struct {
		Package     *ssa.Package
		Application *types.Interface
	}

	Applications []*config.Application
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

		return true
	}

	return false
}

func (ctx *context) LookupMethod(t types.Type, name string) *ssa.Function {
	fn := ctx.Program.LookupMethod(t, packageOf(t), name)
	if fn == nil {
		panic("method not found")
	}
	return fn
}
