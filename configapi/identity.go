package configapi

import (
	"github.com/dogmatiq/configkit"
	"github.com/dogmatiq/enginekit/configapi/internal/pb"
)

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
