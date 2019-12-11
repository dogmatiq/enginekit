package configapi

import (
	"errors"

	"github.com/dogmatiq/configkit/message"
	"github.com/dogmatiq/enginekit/configapi/internal/pb"
)

func marshalRole(in message.Role) (pb.MessageRole, error) {
	if err := in.Validate(); err != nil {
		return pb.MessageRole_UNKNOWN_MESSAGE_ROLE, err
	}

	switch in {
	case message.CommandRole:
		return pb.MessageRole_COMMAND, nil
	case message.EventRole:
		return pb.MessageRole_EVENT, nil
	default: // message.TimeoutRole
		return pb.MessageRole_TIMEOUT, nil
	}
}

func unmarshalRole(in pb.MessageRole) (message.Role, error) {
	switch in {
	case pb.MessageRole_COMMAND:
		return message.CommandRole, nil
	case pb.MessageRole_EVENT:
		return message.EventRole, nil
	case pb.MessageRole_TIMEOUT:
		return message.TimeoutRole, nil
	default:
		return "", errors.New("TODO")
	}
}

func marshalRoles(
	in message.NameRoles,
) (
	[]*pb.MessageRolePair,
	map[message.Name]uint32,
	error,
) {
	var (
		out     []*pb.MessageRolePair
		indices = map[message.Name]uint32{}
		index   uint32
	)

	for n, r := range in {
		p := &pb.MessageRolePair{}

		var err error
		p.Name, err = n.MarshalBinary()
		if err != nil {
			return nil, nil, err
		}

		p.Role, err = marshalRole(r)
		if err != nil {
			return nil, nil, err
		}

		out = append(out, p)
		indices[n] = index
		index++
	}

	return out, indices, nil
}

func unmarshalRoles(in []*pb.MessageRolePair) (
	message.NameRoles,
	[]rolePair,
	error,
) {
	var (
		out     = message.NameRoles{}
		indices []rolePair
	)

	for _, p := range in {
		var n message.Name

		if err := n.UnmarshalBinary(p.Name); err != nil {
			return nil, nil, err
		}

		r, err := unmarshalRole(p.Role)
		if err != nil {
			return nil, nil, err
		}

		out[n] = r
		indices = append(indices, rolePair{n, r})
	}

	return out, indices, nil
}

type rolePair struct {
	Name message.Name
	Role message.Role
}
