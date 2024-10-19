package staticconfig

import (
	"cmp"
	"fmt"
	"go/types"
	"iter"
	"slices"

	"github.com/dogmatiq/enginekit/config"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"
)

// Analysis encapsulates the results of static analysis.
type Analysis struct {
	Applications []*config.Application
	Artifacts    Artifacts
}

// Artifacts contains the intermediate results of the analysis.
type Artifacts struct {
	Packages    []*packages.Package
	SSAProgram  *ssa.Program
	SSAPackages []*ssa.Package
}

// Errors returns a sequence of errors that occurred during analysis, not
// including errors with the Dogma configuration itself.
func (a Analysis) Errors() iter.Seq[error] {
	return func(yield func(error) bool) {
		for _, pkg := range a.Artifacts.Packages {
			for _, err := range pkg.Errors {
				if !yield(err) {
					return
				}
			}
		}
	}
}

// PackagesLoadMode is the minimal [packages.LoadMode] required when loading
// packages for analysis by [Analyze].
const PackagesLoadMode = packages.NeedFiles |
	packages.NeedCompiledGoFiles |
	packages.NeedImports |
	packages.NeedTypes |
	packages.NeedSyntax |
	packages.NeedTypesInfo |
	packages.NeedDeps

// LoadAndAnalyze returns the configurations of the [dogma.Application]
// implementations in the Go package at the given directory, and its
// subdirectories.
//
// The configurations are built by statically analyzing the code, which is never
// executed. As a result, the returned configurations may be invalid or
// incomplete. See [config.Fidelity].
func LoadAndAnalyze(dir string) Analysis {
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

	return Analyze(pkgs)
}

// Analyze returns the configurations of the [dogma.Application] implementations
// in the given Go packages.
//
// The configurations are built by statically analyzing the code, which is never
// executed. As a result, the returned configurations may be invalid or
// incomplete. See [config.Fidelity].
//
// The packages must have be loaded from source syntax using the [packages.Load]
// function using [PackagesLoadMode], at a minimum.
func Analyze(pkgs []*packages.Package) Analysis {
	prog, ssaPackages := ssautil.AllPackages(
		pkgs,
		ssa.InstantiateGenerics, // | ssa.SanityCheckFunctions, // TODO: document why this is necessary
	)

	prog.Build()

	ctx := &context{
		Program:  prog,
		Packages: ssaPackages,
		Analysis: &Analysis{
			Artifacts: struct {
				Packages    []*packages.Package
				SSAProgram  *ssa.Program
				SSAPackages []*ssa.Package
			}{
				pkgs,
				prog,
				ssaPackages,
			},
		},
	}

	if !findDogma(ctx) {
		// If the dogma package is not found as an import, none of the packages
		// can possibly have types that implement [dogma.Application] because
		// doing so requires referring to [dogma.ApplicationConfigurer].
		return *ctx.Analysis
	}

	for _, pkg := range ctx.Packages {
		if pkg == nil {
			// Any [packages.Package] that can not be built results in a nil
			// [ssa.Package]. We ignore any such packages so that we can still
			// obtain information about applications from other valid packages.
			continue
		}

		for _, m := range pkg.Members {
			if t, ok := m.(*ssa.Type); ok {
				analyzeType(ctx, t)
			}
		}
	}

	// Ensure the applications are in a deterministic order.
	slices.SortFunc(
		ctx.Analysis.Applications,
		func(a, b *config.Application) int {
			return cmp.Compare(
				a.String(),
				b.String(),
			)
		},
	)

	return *ctx.Analysis
}

// packageOf returns the package in which t is declared.
//
// It panics if t is not a named type or a pointer to a named type.
func packageOf(t types.Type) *types.Package {
	switch t := t.(type) {
	case *types.Named:
		return t.Obj().Pkg()
	case *types.Alias:
		return t.Obj().Pkg()
	case *types.Pointer:
		return packageOf(t.Elem())
	default:
		panic(fmt.Sprintf("cannot determine package for anonymous or built-in type %v", t))
	}
}

func analyzeType(ctx *context, m *ssa.Type) {
	t := m.Type()

	if isAbstract(t) {
		// We're only interested in concrete types; otherwise there's nothing to
		// analyze!
		return
	}

	// The sequence of the if-blocks below is important as a type
	// implements an interface only if the methods in the interface's
	// method set have non-pointer receivers. Hence the implementation
	// check for the "raw" (non-pointer) type is made first.
	//
	// A pointer to the type, on the other hand, implements the
	// interface regardless of whether pointer receivers are used or
	// not.

	if types.Implements(t, ctx.Dogma.Application) {
		analyzeApplication(ctx, t)
		return
	}

	p := types.NewPointer(t)
	if types.Implements(p, ctx.Dogma.Application) {
		analyzeApplication(ctx, p)
		return
	}
}
