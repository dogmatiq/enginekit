package configapi

import (
	"context"

	"github.com/dogmatiq/configkit"
)

// application is an implementation of config.Application.
type application struct {
	identity configkit.Identity
	typeName string
	messages configkit.EntityMessageNames
	handlers configkit.HandlerSet
}

func (a *application) Identity() configkit.Identity {
	return a.identity
}

func (a *application) TypeName() string {
	return a.typeName
}

func (a *application) MessageNames() configkit.EntityMessageNames {
	return a.messages
}

func (a *application) Handlers() configkit.HandlerSet {
	return a.handlers
}

func (a *application) AcceptVisitor(ctx context.Context, v configkit.Visitor) error {
	return v.VisitApplication(ctx, a)
}

// handler is an implementation of config.Handler.
type handler struct {
	identity    configkit.Identity
	typeName    string
	messages    configkit.EntityMessageNames
	handlerType configkit.HandlerType
}

func (h *handler) Identity() configkit.Identity {
	return h.identity
}

func (h *handler) TypeName() string {
	return h.typeName
}

func (h *handler) MessageNames() configkit.EntityMessageNames {
	return h.messages
}

func (h *handler) HandlerType() configkit.HandlerType {
	return h.handlerType
}

func (h *handler) AcceptVisitor(ctx context.Context, v configkit.Visitor) error {
	h.handlerType.MustValidate()

	switch h.handlerType {
	case configkit.AggregateHandlerType:
		return v.VisitAggregate(ctx, h)
	case configkit.ProcessHandlerType:
		return v.VisitProcess(ctx, h)
	case configkit.IntegrationHandlerType:
		return v.VisitIntegration(ctx, h)
	default: // configkit.ProjectionHandlerType
		return v.VisitProjection(ctx, h)
	}
}
