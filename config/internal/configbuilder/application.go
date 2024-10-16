package configbuilder

import (
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/config"
)

// Application returns an [ApplicationBuilder] that builds a new
// [config.Application].
func Application() *ApplicationBuilder {
	return &ApplicationBuilder{}
}

// ApplicationBuilder constructs a [config.Application].
type ApplicationBuilder struct {
	target config.Application
}

// SetSourceTypeName sets the source of the configuration.
func (b *ApplicationBuilder) SetSourceTypeName(typeName string) *ApplicationBuilder {
	setSourceTypeName(&b.target.AsConfigured.Source, typeName)
	return b
}

// SetSource sets the source of the configuration.
func (b *ApplicationBuilder) SetSource(app dogma.Application) *ApplicationBuilder {
	setSource(&b.target.AsConfigured.Source, app)
	return b
}

// AddIdentity returns an [IdentityBuilder] that adds a [config.Identity] to the
// application.
func (b *ApplicationBuilder) AddIdentity() *IdentityBuilder {
	return &IdentityBuilder{appendTo: &b.target.AsConfigured.Identities}
}

// BuildIdentity calls fn which configures a [config.Identity] that is added to
// the application.
func (b *ApplicationBuilder) BuildIdentity(
	fn func(*IdentityBuilder) config.Fidelity,
) *ApplicationBuilder {
	x := b.AddIdentity()
	x.Done(fn(x))
	return b
}

// AddAggregate returns an [AggregateBuilder] that adds a [config.Aggregate] to
// the application.
func (b *ApplicationBuilder) AddAggregate() *AggregateBuilder {
	return &AggregateBuilder{appendTo: &b.target.AsConfigured.Handlers}
}

// BuildAggregate calls fn which configures a [config.Aggregate] that is added
// to the application.
func (b *ApplicationBuilder) BuildAggregate(
	fn func(*AggregateBuilder) config.Fidelity,
) *ApplicationBuilder {
	x := b.AddAggregate()
	x.Done(fn(x))
	return b
}

// AddProcess returns an [ProcessBuilder] that adds a [config.Process] to the
// application.
func (b *ApplicationBuilder) AddProcess() *ProcessBuilder {
	return &ProcessBuilder{appendTo: &b.target.AsConfigured.Handlers}
}

// BuildProcess calls fn which configures a [config.Process] that is added to
// the application.
func (b *ApplicationBuilder) BuildProcess(
	fn func(*ProcessBuilder) config.Fidelity,
) *ApplicationBuilder {
	x := b.AddProcess()
	x.Done(fn(x))
	return b
}

// AddIntegration returns an [IntegrationBuilder] that adds a
// [config.Integration] to the application.
func (b *ApplicationBuilder) AddIntegration() *IntegrationBuilder {
	return &IntegrationBuilder{appendTo: &b.target.AsConfigured.Handlers}
}

// BuildIntegration calls fn which configures a [config.Integration] that is
// added to the application.
func (b *ApplicationBuilder) BuildIntegration(
	fn func(*IntegrationBuilder) config.Fidelity,
) *ApplicationBuilder {
	x := b.AddIntegration()
	x.Done(fn(x))
	return b
}

// AddProjection returns an [ProjectionBuilder] that adds a [config.Projection]
// to the application.
func (b *ApplicationBuilder) AddProjection() *ProjectionBuilder {
	return &ProjectionBuilder{appendTo: &b.target.AsConfigured.Handlers}
}

// BuildProjection calls fn which configures a [config.Projection] that is added
// to the application.
func (b *ApplicationBuilder) BuildProjection(
	fn func(*ProjectionBuilder) config.Fidelity,
) *ApplicationBuilder {
	x := b.AddProjection()
	x.Done(fn(x))
	return b
}

// Edit calls fn, which can apply arbitrary changes to the application.
func (b *ApplicationBuilder) Edit(fn func(*config.ApplicationAsConfigured)) *ApplicationBuilder {
	fn(&b.target.AsConfigured)
	return b
}

// Done completes the configuration of the application.
func (b *ApplicationBuilder) Done(f config.Fidelity) *config.Application {
	if f&config.Incomplete == 0 {
		if !b.target.AsConfigured.Source.TypeName.IsPresent() {
			panic("aggregate must have a source or be marked as incomplete")
		}
	}

	b.target.AsConfigured.Fidelity = f

	return &b.target
}
