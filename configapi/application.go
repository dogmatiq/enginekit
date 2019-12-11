package configapi

import (
	"context"

	"github.com/dogmatiq/configkit"
	"github.com/dogmatiq/configkit/message"
	"github.com/dogmatiq/enginekit/configapi/internal/pb"
)

// application is an implementation of config.Application.
type application struct {
	identity configkit.Identity
	typeName string
	messages configkit.EntityMessageNames
	foreign  configkit.EntityMessageNames
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

func (a *application) ForeignMessageNames() configkit.EntityMessageNames {
	return a.foreign
}

func (a *application) Handlers() configkit.HandlerSet {
	return a.handlers
}

func (a *application) AcceptVisitor(ctx context.Context, v configkit.Visitor) error {
	return v.VisitApplication(ctx, a)
}

// marshalApplication marshals an application config to its protobuf
// representation.
func marshalApplication(in configkit.Application) (*pb.Application, error) {
	out := pb.Application{
		TypeName: in.TypeName(),
	}

	var err error
	out.Identity, err = marshalIdentity(in.Identity())
	if err != nil {
		return nil, err
	}

	var indices map[message.Name]uint32
	out.Messages, indices, err = marshalRoles(in.MessageNames().Roles)
	if err != nil {
		return nil, err
	}

	for _, ih := range in.Handlers() {
		oh, err := marshalHandler(indices, ih)
		if err != nil {
			return nil, err
		}

		out.Handlers = append(out.Handlers, oh)
	}

	return &out, nil
}

// unmarshalApplication unmarshals an application config from its protobuf
// representation.
func unmarshalApplication(in *pb.Application) (configkit.Application, error) {
	out := application{
		typeName: in.GetTypeName(),
		messages: configkit.EntityMessageNames{
			Produced: message.NameRoles{},
			Consumed: message.NameRoles{},
		},
		foreign: configkit.EntityMessageNames{
			Roles:    message.NameRoles{},
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

	var indices []rolePair
	out.messages.Roles, indices, err = unmarshalRoles(in.GetMessages())
	if err != nil {
		return nil, err
	}

	for _, ih := range in.Handlers {
		ih, err := unmarshalHandler(indices, ih)
		if err != nil {
			return nil, err
		}

		out.handlers.Add(ih)

		for n, r := range ih.MessageNames().Produced {
			out.messages.Produced[n] = r
		}

		for n, r := range ih.MessageNames().Consumed {
			out.messages.Consumed[n] = r
		}
	}

	for n, r := range out.messages.Consumed {
		if _, ok := out.messages.Produced[n]; ok {
			continue
		}

		out.foreign.Roles[n] = r
		out.foreign.Consumed[n] = r
	}

	for n, r := range out.messages.Produced {
		if _, ok := out.messages.Consumed[n]; ok {
			continue
		}

		if r == message.CommandRole {
			out.foreign.Roles[n] = r
			out.foreign.Produced[n] = r
		}
	}

	return &out, nil
}
