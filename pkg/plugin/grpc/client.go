package plugin

import (
	"net"

	"github.com/dapr/components-contrib/pubsub"
	"github.com/dapr/components-contrib/state"
	"github.com/dapr/dapr/pkg/plugin"
	stateproto "github.com/dapr/dapr/pkg/proto/state/v1"
	statesdk "github.com/dapr/dapr/pkg/sdk/state/v1"
	"google.golang.org/grpc"
)

type Client struct {
	connection *grpc.ClientConn
}

func CreateClient(address net.Addr) (plugin.Client, error) {
	serverAddress := address.String()
	conn, err := grpc.Dial(serverAddress)
	if err != nil {
		return nil, err
	}
	return &Client{
		connection: conn,
	}, nil
}

func (c *Client) Store(name string) (state.Store, error) {
	client := stateproto.NewStoreClient(c.connection)
	return statesdk.NewGRPCClient(client), nil
}

func (c *Client) PubSub(name string) (pubsub.PubSub, error) {
	return nil, nil
}
