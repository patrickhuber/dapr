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

type Plugin struct {
	connection *grpc.ClientConn
}

func CreatePlugin(address net.Addr) (plugin.Plugin, error) {
	serverAddress := address.String()
	conn, err := grpc.Dial(serverAddress)
	if err != nil {
		return nil, err
	}
	return &Plugin{
		connection: conn,
	}, nil
}

func (c *Plugin) Init(m *plugin.Metadata) error {
	return nil
}

func (c *Plugin) Store() (state.Store, error) {
	client := stateproto.NewStoreClient(c.connection)
	return statesdk.NewGRPCClient(client), nil
}

func (c *Plugin) PubSub() (pubsub.PubSub, error) {
	return nil, nil
}
