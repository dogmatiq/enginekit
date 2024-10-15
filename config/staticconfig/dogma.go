package staticconfig

import (
	"go/types"
)

const (
	// dogmaPkgPath is the full path of dogma package.
	dogmaPkgPath = "github.com/dogmatiq/dogma"
)

// dogma encapsulates information about the dogma package.
type dogmaPkg struct {
	Package *types.Package
}

// lookupDogmaPackage returns information about the dogma package.
//
// It returns false if the Dogma package has not been imported.
func lookupDogmaPackage(ctx *analysisContext) bool {
	pkg := ctx.SSAProgram.ImportedPackage(dogmaPkgPath)
	if pkg == nil {
		return false
	}

	ctx.Dogma = dogmaPkg{
		Package: pkg.Pkg,
	}

	return true
}
