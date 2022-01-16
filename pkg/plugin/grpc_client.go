package plugin

import (
	"net"

	"github.com/dapr/components-contrib/pubsub"
	"github.com/dapr/components-contrib/state"
	stateproto "github.com/dapr/dapr/pkg/proto/state/v1"
	statesdk "github.com/dapr/dapr/pkg/sdk/state/v1"
	"google.golang.org/grpc"
)

type GRPCClient struct {
	connection *grpc.ClientConn
}

func CreateGRPCCLient(address net.Addr) (Client, error) {
	serverAddress := address.String()
	conn, err := grpc.Dial(serverAddress)
	if err != nil {
		return nil, err
	}
	return &GRPCClient{
		connection: conn,
	}, nil
}

func (c *GRPCClient) Store(name string) (state.Store, error) {
	client := stateproto.NewStoreClient(c.connection)
	return statesdk.NewGRPCClient(client), nil
}

func (c *GRPCClient) PubSub(name string) (pubsub.PubSub, error) {
	return nil, nil
}
