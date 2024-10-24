package configbuilder

import (
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/config"
)

// Aggregate returns a new [config.AggregateX] as configured by fn.
func Aggregate(fn func(*AggregateBuilder)) *config.Aggregate {
	x := &AggregateBuilder{}
	fn(x)
	return x.Done()
}

// AggregateBuilder constructs a [config.Aggregate].
type AggregateBuilder struct {
	target config.Aggregate
}

// SetSourceTypeName sets the source of the configuration.
func (b *AggregateBuilder) SetSourceTypeName(typeName string) {
	setSourceTypeName(&b.target.Properties.Source, typeName)
}

// SetSource sets the source of the configuration.
func (b *AggregateBuilder) SetSource(h dogma.AggregateMessageHandler) {
	setSource(&b.target.Properties.Source, h)
}

// Identity calls fn which configures a [config.Identity] that is added to the
// handler.
func (b *AggregateBuilder) Identity(fn func(*IdentityBuilder)) {
	x := &IdentityBuilder{}
	fn(x)
	b.target.Properties.Identities = append(
		b.target.Properties.Identities,
		x.Done(),
	)
}

// Route calls fn which configures a [config.Route] that is added to the
// handler.
func (b *AggregateBuilder) Route(fn func(*RouteBuilder)) {
	x := &RouteBuilder{}
	fn(x)
	b.target.Properties.Routes = append(
		b.target.Properties.Routes,
		x.Done(),
	)
}

// Disable calls fn which configures a [config.Flag] that indicates whether the
// handler is disabled.
func (b *AggregateBuilder) Disable(fn func(*FlagBuilder[config.Disabled])) {
	x := &FlagBuilder[config.Disabled]{}
	fn(x)
	b.target.Properties.DisabledFlags = append(
		b.target.Properties.DisabledFlags,
		x.Done(),
	)
}

// Edit calls fn, which can apply arbitrary changes to the handler.
func (b *AggregateBuilder) Edit(fn func(*config.AggregateProperties)) {
	fn(&b.target.Properties)
}

// Fidelity returns the fidelity of the configuration.
func (b *AggregateBuilder) Fidelity() config.Fidelity {
	return b.target.Properties.Fidelity
}

// UpdateFidelity merges f with the current fidelity of the configuration.
func (b *AggregateBuilder) UpdateFidelity(f config.Fidelity) {
	b.target.Properties.Fidelity |= f
}

// Done completes the configuration of the handler.
func (b *AggregateBuilder) Done() *config.Aggregate {
	if b.target.Properties.Fidelity&config.Incomplete == 0 {
		if !b.target.Properties.Source.TypeName.IsPresent() {
			panic("handler must have a source or be marked as incomplete")
		}
	}

	return &b.target
}
