package grpc

import (
	"net"

	"github.com/dapr/components-contrib/configuration"
	"github.com/dapr/components-contrib/pubsub"
	"github.com/dapr/components-contrib/state"
	"github.com/dapr/dapr/pkg/plugin"
	stateproto "github.com/dapr/dapr/pkg/proto/state/v1"
	statesdk "github.com/dapr/dapr/pkg/sdk/state/v1"
	"github.com/dapr/kit/logger"
	"google.golang.org/grpc"
)

type Plugin struct {
	connection *grpc.ClientConn
	logger     logger.Logger
}

func NewPlugin(logger logger.Logger) plugin.Plugin {
	return &Plugin{
		logger: logger,
	}
}

func (c *Plugin) Init(m configuration.Metadata) error {

	cfg := plugin.MapComponentAPIToConfig(m)
	address, err := c.startContainer(cfg.Container)
	if err != nil {
		return err
	}
	conn, err := grpc.Dial(address.String())
	if err != nil {
		return err
	}
	c.connection = conn
	return nil
}

func (c *Plugin) startContainer(container *plugin.Container) (net.Addr, error) {
	var address net.Addr
	return address, nil
}

func (c *Plugin) Store() (state.Store, error) {
	client := stateproto.NewStoreClient(c.connection)
	return statesdk.NewGRPCClient(client), nil
}

func (c *Plugin) PubSub() (pubsub.PubSub, error) {
	return nil, nil
}
