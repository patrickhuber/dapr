package kubernetes

import (
	"fmt"
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
	cfg        plugin.Config
	logger     logger.Logger
	discovery  Discovery
}

type Metadata struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
	Address string `yaml:"address"`
	Port    int    `yaml:"port"`
}

func NewPlugin(logger logger.Logger, cfg plugin.Config, discovery Discovery) plugin.Plugin {
	return &Plugin{
		logger:    logger,
		cfg:       cfg,
		discovery: discovery,
	}
}

func (p *Plugin) Init(m configuration.Metadata) error {

	// use the discovery to locate the plugin
	metadata, ok, err := p.discovery.Lookup(p.cfg.Name, p.cfg.Version)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("unable to locate plugin metadata for plugin %s version %s", p.cfg.Name, p.cfg.Version)
	}

	ipAddress := net.ParseIP(metadata.Address)
	if ipAddress == nil {
		return fmt.Errorf("%s is not a valid ip address", metadata.Address)
	}

	addressAndPort := fmt.Sprintf("%s:%d", metadata.Address, metadata.Port)

	conn, err := grpc.Dial(addressAndPort)
	if err != nil {
		return err
	}
	p.connection = conn
	return nil
}

func (p *Plugin) Store() (state.Store, error) {
	client := stateproto.NewStoreClient(p.connection)
	return statesdk.NewGRPCClient(client), nil
}

func (p *Plugin) PubSub() (pubsub.PubSub, error) {
	return nil, nil
}
