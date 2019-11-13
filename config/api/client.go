package api

import (
	"context"

	"github.com/dogmatiq/enginekit/config"
	"github.com/dogmatiq/enginekit/config/api/internal/pb"
	"github.com/dogmatiq/enginekit/identity"
	"github.com/dogmatiq/enginekit/marshaling"
	"google.golang.org/grpc"
)

// Client is used to query a server about it's application configurations.
type Client struct {
	Connection *grpc.ClientConn
	Marshaler  *marshaling.Marshaler
}

// ListApplicationIdentities returns the identities of applications hosted by
// the server.
func (c *Client) ListApplicationIdentities(
	ctx context.Context,
) (_ []identity.Identity, err error) {
	req := &pb.ListApplicationIdentitiesRequest{}
	res, err := pb.NewConfigClient(c.Connection).ListApplicationIdentities(ctx, req)
	if err != nil {
		return nil, err
	}

	defer catch(&err)

	var idents []identity.Identity
	for _, i := range res.Identities {
		idents = append(idents, unmarshalIdentity(i))
	}

	return idents, nil
}

// ListApplications returns the configurations of the applications hosted by
// the server. The handler objects in the returned configuration are nil.
func (c *Client) ListApplications(
	ctx context.Context,
) ([]*config.ApplicationConfig, error) {
	req := &pb.ListApplicationsRequest{}
	res, err := pb.NewConfigClient(c.Connection).ListApplications(ctx, req)
	if err != nil {
		return nil, err
	}

	defer catch(&err)

	var configs []*config.ApplicationConfig
	for _, cfg := range res.Applications {
		configs = append(configs, unmarshalApplication(c.Marshaler, cfg))
	}

	return configs, nil
}

func catch(err *error) {
	r := recover()
	if e, ok := r.(error); ok {
		*err = e
	} else if r != nil {
		panic(r)
	}
}
