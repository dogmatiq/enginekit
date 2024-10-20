package config

import (
	"reflect"

	"github.com/dogmatiq/enginekit/collections/maps"
	"github.com/dogmatiq/enginekit/config/internal/renderer"
	"github.com/dogmatiq/enginekit/internal/typename"
	"github.com/dogmatiq/enginekit/message"
	"github.com/dogmatiq/enginekit/optional"
)

// RouteAsConfigured contains the raw unvalidated properties of a [Route].
type RouteAsConfigured struct {
	// RouteType is the type of route, if available.
	RouteType optional.Optional[RouteType]

	// MessageTypeName is the fully-qualified name of the Go type that
	// implements the [dogma.Message] interface, if available.
	MessageTypeName optional.Optional[string]

	// MessageType is the [message.Type], if available.
	MessageType optional.Optional[message.Type]

	// Fidelity describes the configuration's accuracy in comparison to the
	// actual configuration that would be used at runtime.
	Fidelity Fidelity
}

// Route represents a message route to or from a handler.
type Route struct {
	AsConfigured RouteAsConfigured
}

// RouteType returns the type of route, or panics if the route type is not
// available.
func (r *Route) RouteType() RouteType {
	return r.AsConfigured.RouteType.Get()
}

// MessageType returns the [message.Type] associated with the route, or panics
// if the message type is not available.
func (r *Route) MessageType() message.Type {
	return r.AsConfigured.MessageType.Get()
}

// Fidelity returns information about how well the configuration represents
// the actual configuration that would be used at runtime.
func (r *Route) Fidelity() Fidelity {
	return r.AsConfigured.Fidelity
}

func (r *Route) String() string {
	return RenderDescriptor(r)
}

func (r *Route) renderDescriptor(ren *renderer.Renderer) {
	ren.Print("route")

	if rt, ok := r.AsConfigured.RouteType.TryGet(); ok {
		ren.Print(":", rt.String())
	}

	if mt, ok := r.AsConfigured.MessageTypeName.TryGet(); ok {
		ren.Print(":", typename.Unqualified(mt))
	}
}

func (r *Route) renderDetails(ren *renderer.Renderer) {
	f, errs := validate(r)

	renderFidelity(ren, f, errs)

	if rt, ok := r.AsConfigured.RouteType.TryGet(); ok {
		ren.Print(rt.String(), " ")
	}

	ren.Print("route")

	if mt, ok := r.AsConfigured.MessageTypeName.TryGet(); ok {
		ren.Print(" for ", mt)

		if !r.AsConfigured.MessageType.IsPresent() {
			ren.Print(" (runtime type unavailable)")
		}
	}

	ren.Print("\n")
	renderErrors(ren, errs)
}

func (r *Route) clone() Component {
	return &Route{r.AsConfigured}
}

func (r *Route) normalize(ctx *normalizationContext) {
	routeType, routeTypeOK := r.AsConfigured.RouteType.TryGet()
	typeName, typeNameOK := r.AsConfigured.MessageTypeName.TryGet()
	messageType, messageTypeOK := r.AsConfigured.MessageType.TryGet()

	if !routeTypeOK {
		r.AsConfigured.Fidelity |= Incomplete
	}

	if !typeNameOK {
		r.AsConfigured.Fidelity |= Incomplete
	}

	if messageTypeOK {
		if routeTypeOK && routeType.MessageKind() != messageType.Kind() {
			ctx.Fail(MessageKindMismatchError{routeType, messageType})
		}

		actualTypeName := typename.Get(messageType.ReflectType())
		if typeNameOK && typeName != actualTypeName {
			ctx.Fail(TypeNameMismatchError{typeName, actualTypeName})
		}

		r.AsConfigured.MessageTypeName = optional.Some(actualTypeName)
	} else if ctx.Options.RequireValues {
		ctx.Fail(RuntimeValueUnavailableError{reflect.TypeFor[message.Type]()})
	}
}

func reportRouteErrors(
	ctx *normalizationContext,
	h Handler,
	routes []*Route,
) {
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

	for _, r := range routes {
		if rt, ok := r.AsConfigured.RouteType.TryGet(); ok {
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
