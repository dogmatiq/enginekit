package configbuilder

import (
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/config"
	"github.com/dogmatiq/enginekit/optional"
)

// Process returns a new [config.Process] as configured by fn.
func Process(fn func(*ProcessBuilder)) *config.Process {
	x := &ProcessBuilder{}
	fn(x)
	return x.Done()
}

// ProcessBuilder constructs a [config.Process].
type ProcessBuilder struct {
	target config.Process
}

// SetSourceTypeName sets the source of the configuration.
func (b *ProcessBuilder) SetSourceTypeName(typeName string) {
	setSourceTypeName(&b.target.AsConfigured.Source, typeName)
}

// SetSource sets the source of the configuration.
func (b *ProcessBuilder) SetSource(h dogma.ProcessMessageHandler) {
	setSource(&b.target.AsConfigured.Source, h)
}

// Identity calls fn which configures a [config.Identity] that is added to the
// handler.
func (b *ProcessBuilder) Identity(fn func(*IdentityBuilder)) {
	x := &IdentityBuilder{}
	fn(x)
	b.target.AsConfigured.Identities = append(
		b.target.AsConfigured.Identities,
		x.Done(),
	)
}

// Route calls fn which configures a [config.Route] that is added to the
// handler.
func (b *ProcessBuilder) Route(fn func(*RouteBuilder)) {
	x := &RouteBuilder{}
	fn(x)
	b.target.AsConfigured.Routes = append(
		b.target.AsConfigured.Routes,
		x.Done(),
	)
}

// SetDisabled sets whether the handler is disabled or not.
func (b *ProcessBuilder) SetDisabled(disabled bool) {
	b.target.AsConfigured.IsDisabled = optional.Some(disabled)
}

// Edit calls fn, which can apply arbitrary changes to the handler.
func (b *ProcessBuilder) Edit(fn func(*config.ProcessAsConfigured)) {
	fn(&b.target.AsConfigured)
}

// UpdateFidelity merges f with the current fidelity of the handler.
func (b *ProcessBuilder) UpdateFidelity(f config.Fidelity) {
	b.target.AsConfigured.Fidelity |= f
}

// Done completes the configuration of the handler.
func (b *ProcessBuilder) Done() *config.Process {
	if b.target.AsConfigured.Fidelity&config.Incomplete == 0 {
		if !b.target.AsConfigured.Source.TypeName.IsPresent() {
			panic("handler must have a source or be marked as incomplete")
		}
		if !b.target.AsConfigured.IsDisabled.IsPresent() {
			panic("handler must be known to be enabled or disabled, or be marked as incomplete")
		}
	}

	return &b.target
}
