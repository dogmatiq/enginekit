package configapi

import (
	"context"
	"errors"

	"github.com/dogmatiq/configkit"
	"github.com/dogmatiq/configkit/message"
	"github.com/dogmatiq/enginekit/configapi/internal/pb"
)

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

func marshalHandlerType(in configkit.HandlerType) (pb.HandlerType, error) {
	if err := in.Validate(); err != nil {
		return pb.HandlerType_UNKNOWN_HANDLER_TYPE, err
	}

	switch in {
	case configkit.AggregateHandlerType:
		return pb.HandlerType_AGGREGATE, nil
	case configkit.ProcessHandlerType:
		return pb.HandlerType_PROCESS, nil
	case configkit.IntegrationHandlerType:
		return pb.HandlerType_INTEGRATION, nil
	default: // configkit.ProjectionHandlerType
		return pb.HandlerType_PROJECTION, nil
	}
}

func unmarshalHandlerType(in pb.HandlerType) (configkit.HandlerType, error) {
	switch in {
	case pb.HandlerType_AGGREGATE:
		return configkit.AggregateHandlerType, nil
	case pb.HandlerType_PROCESS:
		return configkit.ProcessHandlerType, nil
	case pb.HandlerType_INTEGRATION:
		return configkit.IntegrationHandlerType, nil
	case pb.HandlerType_PROJECTION:
		return configkit.ProjectionHandlerType, nil
	default:
		return "", errors.New("TODO")
	}
}

func marshalHandler(
	indices map[message.Name]uint32,
	in configkit.Handler,
) (*pb.Handler, error) {
	out := pb.Handler{
		Identity: &pb.Identity{},
		TypeName: in.TypeName(),
	}

	var err error
	out.Identity, err = marshalIdentity(in.Identity())
	if err != nil {
		return nil, err
	}

	out.Type, err = marshalHandlerType(in.HandlerType())
	if err != nil {
		return nil, err
	}

	names := in.MessageNames()

	out.Produced, err = marshalMessageIndices(indices, names.Produced)
	if err != nil {
		return nil, err
	}

	out.Consumed, err = marshalMessageIndices(indices, names.Consumed)
	if err != nil {
		return nil, err
	}

	return &out, nil
}

func unmarshalHandler(
	indices []rolePair,
	in *pb.Handler,
) (configkit.Handler, error) {
	out := handler{
		typeName: in.GetTypeName(),
		messages: configkit.EntityMessageNames{
			Roles:    message.NameRoles{},
			Produced: message.NameRoles{},
			Consumed: message.NameRoles{},
		},
	}

	var err error
	out.identity, err = unmarshalIdentity(in.GetIdentity())
	if err != nil {
		return nil, err
	}

	out.handlerType, err = unmarshalHandlerType(in.GetType())
	if err != nil {
		return nil, err
	}

	for _, i32 := range in.Produced {
		i := int(i32)
		if i >= len(indices) {
			return nil, errors.New("TODO")
		}

		p := indices[i]
		out.messages.Roles[p.Name] = p.Role
		out.messages.Produced[p.Name] = p.Role
	}

	for _, i32 := range in.Consumed {
		i := int(i32)
		if i >= len(indices) {
			return nil, errors.New("TODO")
		}

		p := indices[i]
		out.messages.Roles[p.Name] = p.Role
		out.messages.Consumed[p.Name] = p.Role
	}

	return &out, nil
}

func marshalMessageIndices(
	indices map[message.Name]uint32,
	names message.NameRoles,
) ([]uint32, error) {
	var out []uint32

	for n := range names {
		i, ok := indices[n]
		if !ok {
			return nil, errors.New("TODO")
		}

		out = append(out, i)
	}

	return out, nil
}
