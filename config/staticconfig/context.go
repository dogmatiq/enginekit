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
		Package     *ssa.Package
		Application *types.Interface

		HandlesCommand   *types.Func
		ExecutesCommand  *types.Func
		HandlesEvent     *types.Func
		RecordsEvent     *types.Func
		SchedulesTimeout *types.Func
	}

	Analysis *Analysis
}

// resolveDogmaPackage updates ctx with information about the Dogma package.
//
// It returns false if the Dogma package has not been imported.
func resolveDogmaPackage(ctx *context) bool {
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

		fn := func(n string) *types.Func {
			return pkg.Pkg.
				Scope().
				Lookup(n).(*types.Func)
		}

		ctx.Dogma.Package = pkg
		ctx.Dogma.Application = iface("Application")

		ctx.Dogma.HandlesCommand = fn("HandlesCommand")
		ctx.Dogma.ExecutesCommand = fn("ExecutesCommand")
		ctx.Dogma.HandlesEvent = fn("HandlesEvent")
		ctx.Dogma.RecordsEvent = fn("RecordsEvent")
		ctx.Dogma.SchedulesTimeout = fn("SchedulesTimeout")

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
