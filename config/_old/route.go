package config

import (
	"reflect"

	"github.com/dogmatiq/enginekit/collections/maps"
	"github.com/dogmatiq/enginekit/config/internal/renderer"
	"github.com/dogmatiq/enginekit/internal/typename"
	"github.com/dogmatiq/enginekit/message"
	"github.com/dogmatiq/enginekit/optional"
)

// Route represents a message route to or from a handler.
type Route struct {
	ComponentProperties

	// RouteType is the type of route, if available.
	RouteType optional.Optional[RouteType]

	// MessageTypeName is the fully-qualified name of the Go type that
	// implements the [dogma.Message] interface, if available.
	MessageTypeName optional.Optional[string]

	// MessageType is the [message.Type], if available.
	MessageType optional.Optional[message.Type]
}

func (r *Route) String() string {
	return RenderDescriptor(r)
}

func (r *Route) renderDescriptor(ren *renderer.Renderer) {
	ren.Print("route")

	if rt, ok := r.RouteType.TryGet(); ok {
		ren.Print(":", rt.String())
	}

	if mt, ok := r.MessageTypeName.TryGet(); ok {
		ren.Print(":", typename.Unqualified(mt))
	}
}

func (r *Route) renderDetails(ren *renderer.Renderer) {
	f, errs := validate(r)

	renderFidelity(ren, f, errs)

	if rt, ok := r.RouteType.TryGet(); ok {
		ren.Print(rt.String(), " ")
	}

	ren.Print("route")

	if mt, ok := r.MessageTypeName.TryGet(); ok {
		ren.Print(" for ", mt)

		if !r.MessageType.IsPresent() {
			ren.Print(" (runtime type unavailable)")
		}
	}

	ren.Print("\n")
	renderErrors(ren, errs)
}

func (r *Route) normalize(ctx *normalizationContext) {
	routeType, routeTypeOK := r.RouteType.TryGet()
	typeName, typeNameOK := r.MessageTypeName.TryGet()
	messageType, messageTypeOK := r.MessageType.TryGet()

	if !routeTypeOK {
		r.Fidelity |= Incomplete
	}

	if !typeNameOK {
		r.Fidelity |= Incomplete
	}

	if messageTypeOK {
		if routeTypeOK && routeType.MessageKind() != messageType.Kind() {
			ctx.Fail(MessageKindMismatchError{routeType, messageType})
		}

		actualTypeName := typename.Get(messageType.ReflectType())
		if typeNameOK && typeName != actualTypeName {
			ctx.Fail(TypeNameMismatchError{typeName, actualTypeName})
		}

		r.MessageTypeName = optional.Some(actualTypeName)
	} else if ctx.Options.RequireValues {
		ctx.Fail(RuntimeValueUnavailableError{reflect.TypeFor[message.Type]()})
	}
}

func (r *Route) clone() any {
	return &Route{
		clone(r.ComponentProperties),
		r.RouteType,
		r.MessageTypeName,
		r.MessageType,
	}
}

func reportRouteErrors(ctx *normalizationContext, h Handler) {
	var (
		capabilities = h.HandlerType().RouteCapabilities()
		missing      maps.Ordered[RouteType, MissingRequiredRouteError]
		duplicate    maps.OrderedByKey[routeKey, DuplicateRouteError]
	)

	for rt, req := range capabilities.RouteTypes {
		if req == RouteTypeRequired {
			missing.Set(rt, MissingRequiredRouteError{rt})
		}
	}

	for _, r := range h.CommonHandlerProperties().RouteComponents {
		if rt, ok := r.RouteType.TryGet(); ok {
			if capabilities.RouteTypes[rt] == RouteTypeDisallowed {
				ctx.Fail(UnexpectedRouteTypeError{r})
			} else {
				missing.Remove(rt)
			}
		}

		if k, ok := routeKeyOf(r); ok {
			duplicate.Update(
				k,
				func(err *DuplicateRouteError) {
					err.RouteType = k.RouteType
					err.MessageTypeName = k.MessageTypeName
					err.DuplicateRoutes = append(err.DuplicateRoutes, r)
				},
			)
		}
	}

	for err := range missing.Values() {
		ctx.Fail(err)
	}

	for err := range duplicate.Values() {
		if len(err.DuplicateRoutes) > 1 {
			ctx.Fail(err)
		}
	}
}
