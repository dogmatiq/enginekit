package staticconfig

import (
	"cmp"
	"slices"

	"github.com/dogmatiq/enginekit/config"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"
)

// PackagesLoadMode is the minimal [packages.LoadMode] required when loading
// packages for analysis by [FromPackages].
const PackagesLoadMode = packages.NeedFiles |
	packages.NeedCompiledGoFiles |
	packages.NeedImports |
	packages.NeedTypes |
	packages.NeedSyntax |
	packages.NeedTypesInfo |
	packages.NeedDeps

	// FromDir returns the configurations of the [dogma.Application] in the Go
// package at the given directory, and its subdirectories.
//
// The configurations are built by statically analyzing the code; it is never
// executed. As a result, the returned configurations may be invalid or
// incomplete. See [config.Fidelity].
func FromDir(dir string) Analysis {
	pkgs, err := packages.Load(
		&packages.Config{
			Mode: PackagesLoadMode,
			Dir:  dir,
		},
		"./...",
	)
	if err != nil {
		// According to the documentation of [packages.Load], this error relates
		// only to malformed patterns, which should never occur since it's
		// hardcoded above.
		panic(err)
	}

	return FromPackages(pkgs)
}

// FromPackages returns the configurations of the [dogma.Application] in the
// given Go packages.
//
// The configurations are built by statically analyzing the code; it is never
// executed. As a result, the returned configurations may be invalid or
// incomplete. See [config.Fidelity].
//
// The packages must have be loaded from source syntax using the [packages.Load]
// function using [PackagesLoadMode], at a minimum.
func FromPackages(pkgs []*packages.Package) Analysis {
	ctx := &analysisContext{
		Analysis: Analysis{
			Packages: pkgs,
		},
	}

	ctx.SSAProgram, ctx.SSAPackages = ssautil.AllPackages(
		ctx.Packages,
		0,
		// ssa.SanityCheckFunctions, // TODO: document why this is necessary
		// see.InstantiateGenerics // TODO: might this make some generic handling code easier?
	)

	ctx.SSAProgram.Build()

	if !lookupDogmaPackage(ctx) {
		// If the dogma package is not found as an import, none of the packages
		// can possibly have types that implement [dogma.Application] because
		// doing so requires referring to [dogma.ApplicationConfigurer].
		return ctx.Analysis
	}

	// for _, pkg := range ctx.SSAPackages {
	// 	if pkg == nil {
	// 		// Any [packages.Package] that can not be built results in a nil
	// 		// [ssa.Package]. We ignore any such packages so that we can still
	// 		// obtain information about applications from other valid packages.
	// 		continue
	// 	}

	// 	for _, m := range pkg.Members {
	// 		// The sequence of the if-blocks below is important as a type
	// 		// implements an interface only if the methods in the interface's
	// 		// method set have non-pointer receivers. Hence the implementation
	// 		// check for the non-pointer type is made first.
	// 		//
	// 		// A pointer to the type, on the other hand, implements the
	// 		// interface regardless of whether pointer receivers are used or
	// 		// not.
	// 		if types.Implements(m.Type(), dogmaPkg.Application) {
	// 			apps = append(apps, analyzeApplication(prog, dogmaPkg, m.Type()))
	// 			continue
	// 		}

	// 		if p := types.NewPointer(m.Type()); types.Implements(p, dogmaPkg.Application) {
	// 			apps = append(apps, analyzeApplication(prog, dogmaPkg, p))
	// 		}
	// 	}
	// }

	slices.SortFunc(
		ctx.Analysis.Applications,
		func(a, b *config.Application) int {
			return cmp.Compare(
				a.String(),
				b.String(),
			)
		},
	)

	return ctx.Analysis
}

// Analysis encapsulates the results of static analysis.
type Analysis struct {
	Applications []*config.Application
	Packages     []*packages.Package
	SSAProgram   *ssa.Program
	SSAPackages  []*ssa.Package
}

type analysisContext struct {
	Analysis
	Dogma dogmaPkg
}
