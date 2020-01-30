package configapi

import (
	"errors"
	"fmt"

	"github.com/dogmatiq/configkit"
	"github.com/dogmatiq/configkit/message"
	"github.com/dogmatiq/enginekit/configapi/internal/pb"
)

// marshalApplication marshals an application config to its protobuf
// representation.
//
// See the pb.Application type for details about how the application is
// represented as protocol buffers message.
func marshalApplication(in configkit.Application) (*pb.Application, error) {
	out := &pb.Application{}

	var err error
	out.Identity, err = marshalIdentity(in.Identity())
	if err != nil {
		return nil, err
	}

	out.TypeName = in.TypeName()
	if out.TypeName == "" {
		return nil, errors.New("application type name is empty")
	}

	// indices is mapping of name to index into out.Messages. The index is
	// encoded in the handler to avoid repeating long type names that are likely
	// to be referenced many times within an application.
	indices := map[message.Name]uint32{}

	for _, hIn := range in.Handlers() {
		hOut, err := marshalHandler(out, indices, hIn)
		if err != nil {
			return nil, err
		}

		out.Handlers = append(out.Handlers, hOut)
	}

	return out, nil
}

// marshalHandler marshals a handler config to its protobuf representation.
//
// It populates app.Messages with NameRole pairs for each of the handler's
// consumed/produced messages.
func marshalHandler(
	app *pb.Application,
	indices map[message.Name]uint32,
	in configkit.Handler,
) (*pb.Handler, error) {
	out := &pb.Handler{}

	var err error
	out.Identity, err = marshalIdentity(in.Identity())
	if err != nil {
		return nil, err
	}

	out.TypeName = in.TypeName()
	if out.TypeName == "" {
		return nil, errors.New("handler type name is empty")
	}

	out.Type, err = marshalHandlerType(in.HandlerType())
	if err != nil {
		return nil, err
	}

	names := in.MessageNames()
	out.Produced, err = marshalNameRoles(app, indices, names.Produced)
	if err != nil {
		return nil, err
	}

	out.Consumed, err = marshalNameRoles(app, indices, names.Consumed)
	if err != nil {
		return nil, err
	}

	return out, nil
}

// marshalNameRoles marshals a message.NameRoles collection for a handler into
// protocol buffers application representation.
//
// It populates app.Messages with NameRole pairs for each of the elements in the
// NameRoles collection.
func marshalNameRoles(
	app *pb.Application,
	indices map[message.Name]uint32,
	in message.NameRoles,
) ([]uint32, error) {
	out := make([]uint32, 0, len(in))

	for n, r := range in {
		i, ok := indices[n]

		if !ok {
			nr, err := marshalNameRole(n, r)
			if err != nil {
				return nil, err
			}

			i = uint32(len(app.Messages))
			app.Messages = append(app.Messages, nr)
		}

		out = append(out, i)
	}

	return out, nil
}

// unmarshalApplication unmarshals an application config from its protobuf
// representation.
func unmarshalApplication(in *pb.Application) (configkit.Application, error) {
	out := &application{
		messages: configkit.EntityMessageNames{
			Produced: message.NameRoles{},
			Consumed: message.NameRoles{},
		},
		handlers: configkit.HandlerSet{},
	}

	var err error
	out.identity, err = unmarshalIdentity(in.GetIdentity())
	if err != nil {
		return nil, err
	}

	out.typeName = in.GetTypeName()
	if out.typeName == "" {
		return nil, errors.New("application type name is empty")
	}

	var indices []nameRole
	for _, nrIn := range in.GetMessages() {
		nrOut, err := unmarshalNameRole(nrIn)
		if err != nil {
			return nil, err
		}
		indices = append(indices, nrOut)
	}

	for _, hIn := range in.Handlers {
		hOut, err := unmarshalHandler(indices, hIn)
		if err != nil {
			return nil, err
		}

		out.handlers.Add(hOut)

		for n, r := range hOut.MessageNames().Produced {
			out.messages.Produced[n] = r
		}

		for n, r := range hOut.MessageNames().Consumed {
			out.messages.Consumed[n] = r
		}
	}

	return out, nil
}

// unmarshalHandler unmarshals a handler configuration from its protocol buffers
// representation.
func unmarshalHandler(
	indices []nameRole,
	in *pb.Handler,
) (configkit.Handler, error) {
	out := &handler{
		messages: configkit.EntityMessageNames{
			Produced: message.NameRoles{},
			Consumed: message.NameRoles{},
		},
	}

	var err error
	out.identity, err = unmarshalIdentity(in.GetIdentity())
	if err != nil {
		return nil, err
	}

	out.typeName = in.GetTypeName()
	if out.typeName == "" {
		return nil, errors.New("handler type name is empty")
	}

	out.handlerType, err = unmarshalHandlerType(in.GetType())
	if err != nil {
		return nil, err
	}

	out.messages.Produced, err = unmarshalNameRoles(indices, in.Produced)
	if err != nil {
		return nil, err
	}

	out.messages.Consumed, err = unmarshalNameRoles(indices, in.Consumed)
	if err != nil {
		return nil, err
	}

	return out, nil
}

// unmarshalNameRoles unmarshals a sequence of name/role indices into a
// NameRoles collection.
func unmarshalNameRoles(
	indices []nameRole,
	in []uint32,
) (message.NameRoles, error) {
	out := message.NameRoles{}
	for _, i := range in {
		if i >= uint32(len(indices)) {
			return nil, errors.New("name/role index out of range")
		}

		nr := indices[i]
		out[nr.Name] = nr.Role
	}

	return out, nil
}

// marshalIdentity marshals a configkit.Identity to its protocol buffers
// representation.
func marshalIdentity(in configkit.Identity) (*pb.Identity, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}

	return &pb.Identity{
		Name: in.Name,
		Key:  in.Key,
	}, nil
}

// unmarshalIdentity unmarshals a configkit.Identity from its protocol buffers
// representation.
func unmarshalIdentity(in *pb.Identity) (configkit.Identity, error) {
	return configkit.NewIdentity(
		in.Name,
		in.Key,
	)
}

// marshalHandlerType marshals a configkit.HandlerType to its protocol buffers
// representation.
func marshalHandlerType(t configkit.HandlerType) (pb.HandlerType, error) {
	if err := t.Validate(); err != nil {
		return pb.HandlerType_UNKNOWN_HANDLER_TYPE, err
	}

	switch t {
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

// unmarshalHandlerType unmarshals a configkit.HandlerType from its protocol
// buffers representation.
func unmarshalHandlerType(t pb.HandlerType) (configkit.HandlerType, error) {
	switch t {
	case pb.HandlerType_AGGREGATE:
		return configkit.AggregateHandlerType, nil
	case pb.HandlerType_PROCESS:
		return configkit.ProcessHandlerType, nil
	case pb.HandlerType_INTEGRATION:
		return configkit.IntegrationHandlerType, nil
	case pb.HandlerType_PROJECTION:
		return configkit.ProjectionHandlerType, nil
	default:
		return "", fmt.Errorf("unknown handler type: %#v", t)
	}
}

// marshalMessageRole marshals a message.Role to its protocol buffers
// representation.
func marshalMessageRole(r message.Role) (pb.MessageRole, error) {
	if err := r.Validate(); err != nil {
		return pb.MessageRole_UNKNOWN_MESSAGE_ROLE, err
	}

	switch r {
	case message.CommandRole:
		return pb.MessageRole_COMMAND, nil
	case message.EventRole:
		return pb.MessageRole_EVENT, nil
	default: // message.TimeoutRole
		return pb.MessageRole_TIMEOUT, nil
	}
}

// unmarshalMessageRole unmarshals a message.Role from its protocol buffers
// representation.
func unmarshalMessageRole(r pb.MessageRole) (message.Role, error) {
	switch r {
	case pb.MessageRole_COMMAND:
		return message.CommandRole, nil
	case pb.MessageRole_EVENT:
		return message.EventRole, nil
	case pb.MessageRole_TIMEOUT:
		return message.TimeoutRole, nil
	default:
		return "", fmt.Errorf("unknown message role: %#v", r)
	}
}

// marshalNameRole marshals a message name and role into a protocol buffers
// NameRole message.
func marshalNameRole(n message.Name, r message.Role) (*pb.NameRole, error) {
	nr := &pb.NameRole{}

	var err error
	nr.Name, err = n.MarshalBinary()
	if err != nil {
		return nil, err
	}

	nr.Role, err = marshalMessageRole(r)
	if err != nil {
		return nil, err
	}

	return nr, nil
}

// nameRole is an in-memory representation of an unmarshaled *pb.NameRole
// message.
type nameRole struct {
	Name message.Name
	Role message.Role
}

// unmarshalNameRole unmarshals a *pb.NameRole to a nameRole.
func unmarshalNameRole(in *pb.NameRole) (nameRole, error) {
	var out nameRole

	err := out.Name.UnmarshalBinary(in.Name)
	if err != nil {
		return out, err
	}

	out.Role, err = unmarshalMessageRole(in.Role)
	return out, err
}
