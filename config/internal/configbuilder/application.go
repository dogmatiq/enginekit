package configbuilder

import (
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/config"
)

// Application returns a new [config.Application] as configured by fn.
func Application(fn func(*ApplicationBuilder)) *config.Application {
	x := &ApplicationBuilder{}
	fn(x)
	return x.Done()
}

// ApplicationBuilder constructs a [config.Application].
type ApplicationBuilder struct {
	target config.Application
}

// SetSourceTypeName sets the source of the configuration.
func (b *ApplicationBuilder) SetSourceTypeName(typeName string) {
	setSourceTypeName(&b.target.AsConfigured.Source, typeName)
}

// SetSource sets the source of the configuration.
func (b *ApplicationBuilder) SetSource(app dogma.Application) {
	setSource(&b.target.AsConfigured.Source, app)
}

// Identity calls fn which configures a [config.Identity] that is added to the
// application.
func (b *ApplicationBuilder) Identity(fn func(*IdentityBuilder)) {
	x := &IdentityBuilder{}
	fn(x)
	b.target.AsConfigured.Identities = append(
		b.target.AsConfigured.Identities,
		x.Done(),
	)
}

// Aggregate calls fn which configures a [config.Aggregate] that is added
// to the application.
func (b *ApplicationBuilder) Aggregate(fn func(*AggregateBuilder)) {
	x := &AggregateBuilder{}
	fn(x)
	b.target.AsConfigured.Handlers = append(
		b.target.AsConfigured.Handlers,
		x.Done(),
	)
}

// Process calls fn which configures a [config.Process] that is added to the
// application.
func (b *ApplicationBuilder) Process(fn func(*ProcessBuilder)) {
	x := &ProcessBuilder{}
	fn(x)
	b.target.AsConfigured.Handlers = append(
		b.target.AsConfigured.Handlers,
		x.Done(),
	)
}

// Integration calls fn which configures a [config.Integration] that is added to
// the application.
func (b *ApplicationBuilder) Integration(fn func(*IntegrationBuilder)) {
	x := &IntegrationBuilder{}
	fn(x)
	b.target.AsConfigured.Handlers = append(
		b.target.AsConfigured.Handlers,
		x.Done(),
	)
}

// Projection calls fn which configures a [config.Projection] that is added to
// the application.
func (b *ApplicationBuilder) Projection(fn func(*ProjectionBuilder)) {
	x := &ProjectionBuilder{}
	fn(x)
	b.target.AsConfigured.Handlers = append(
		b.target.AsConfigured.Handlers,
		x.Done(),
	)
}

// Edit calls fn, which can apply arbitrary changes to the application.
func (b *ApplicationBuilder) Edit(fn func(*config.ApplicationAsConfigured)) {
	fn(&b.target.AsConfigured)
}

// Fidelity returns the fidelity of the configuration.
func (b *ApplicationBuilder) Fidelity() config.Fidelity {
	return b.target.AsConfigured.Fidelity
}

// UpdateFidelity merges f with the current fidelity of the configuration.
func (b *ApplicationBuilder) UpdateFidelity(f config.Fidelity) {
	b.target.AsConfigured.Fidelity |= f
}

// Done completes the configuration of the application.
func (b *ApplicationBuilder) Done() *config.Application {
	if b.target.AsConfigured.Fidelity&config.Incomplete == 0 {
		if !b.target.AsConfigured.Source.TypeName.IsPresent() {
			panic("application must have a source or be marked as incomplete")
		}
	}

	return &b.target
}
