package configapi

import (
	"context"

	"github.com/dogmatiq/configkit"
	"github.com/dogmatiq/enginekit/configapi/internal/pb"
	"google.golang.org/grpc"
)

// Client is used to query a server about its application configurations.
type Client struct {
	Connection *grpc.ClientConn
}

// ListApplicationIdentities returns the identities of applications hosted by
// the server.
func (c *Client) ListApplicationIdentities(
	ctx context.Context,
) (_ []configkit.Identity, err error) {
	req := &pb.ListApplicationIdentitiesRequest{}
	res, err := pb.NewConfigClient(c.Connection).ListApplicationIdentities(ctx, req)
	if err != nil {
		return nil, err
	}

	var idents []configkit.Identity
	for _, in := range res.Identities {
		out, err := unmarshalIdentity(in)
		if err != nil {
			return nil, err
		}

		idents = append(idents, out)
	}

	return idents, nil
}

// ListApplications returns the configurations of the applications hosted by
// the server. The handler objects in the returned configuration are nil.
func (c *Client) ListApplications(
	ctx context.Context,
) ([]configkit.Application, error) {
	req := &pb.ListApplicationsRequest{}
	res, err := pb.NewConfigClient(c.Connection).ListApplications(ctx, req)
	if err != nil {
		return nil, err
	}

	var configs []configkit.Application
	for _, in := range res.Applications {
		out, err := unmarshalApplication(in)
		if err != nil {
			return nil, err
		}

		configs = append(configs, out)
	}

	return configs, nil
}
