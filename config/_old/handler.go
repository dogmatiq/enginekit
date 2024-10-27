package config

import "github.com/dogmatiq/enginekit/optional"

// A Handler is a specialization of [QEntity] that represents configuration of a
// Dogma message handler.
type Handler interface {
	QEntity

	// HandlerType returns [HandlerType] of the handler.
	HandlerType() HandlerType

	// IsDisabled returns true if the handler is disabled.
	//
	// It panics if the disabled state cannot be determined.
	IsDisabled() bool

	// CommonHandlerProperties returns the (possibly invalid or incomplete)
	// properties of the handler.
	CommonHandlerProperties() *HandlerProperties
}

// HandlerProperties contains the properties common to all [Handler] implementations.
type HandlerProperties struct {
	EntityProperties

	// RouteComponents is the list of routes configured on the handler.
	RouteComponents []*Route

	// DisabledFlag is a [Flag] that indicates whether the handler is disabled.
	DisabledFlag Flag[Disabled]
}

// CommonHandlerProperties returns the (possibly invalid or incomplete)
// properties of the handler.
func (p *HandlerProperties) CommonHandlerProperties() *HandlerProperties {
	return p
}

func (p HandlerProperties) clone() any {
	cloneInPlace(&p.EntityProperties)
	cloneSliceInPlace(&p.RouteComponents)
	clonePointeeInPlace(&p.DisabledFlag)
	return p
}

// Disabled is the label for a [FlagModification] that indicates that a [Handler] has been
// disabled.
type Disabled struct{ flagSymbol }

func resolveIsDisabled(Handler) bool {
	panic("not implemented")
	// 	ctx := strictContext(h)

	// 	h.handlerProperties().DisabledFlags.resolve(h.Fidelity())
}

func normalizeHandler[T any](
	ctx *normalizationContext,
	h Handler,
	v optional.Optional[T],
) {
	normalizeEntity(ctx, h, v)

	p := h.CommonHandlerProperties()
	normalizeChildren(ctx, p.RouteComponents...)
	normalizeChildren(ctx, &p.DisabledFlag)

	reportRouteErrors(ctx, h)
}
