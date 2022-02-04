package standalone

import (
	"fmt"
	"path/filepath"

	"github.com/dapr/components-contrib/configuration"
	"github.com/dapr/components-contrib/pubsub"
	"github.com/dapr/components-contrib/state"
	"github.com/dapr/dapr/pkg/plugin"
	"github.com/dapr/dapr/pkg/sdk"
	"github.com/dapr/kit/logger"
	goplugin "github.com/hashicorp/go-plugin"

	state_sdk "github.com/dapr/dapr/pkg/sdk/state/v1"
)

type Plugin struct {
	clientProtocol goplugin.ClientProtocol
	logger         logger.Logger
}

func NewPlugin(logger logger.Logger) plugin.Plugin {
	return &Plugin{
		logger: logger,
	}
}

func (p *Plugin) Init(m configuration.Metadata) error {
	cfg := plugin.MapComponentAPIToConfig(m)
	pluginSet := CreatePluginSet(cfg)

	runtimeContext := GetRuntimeContextFromString(cfg.Run.Runtime)
	path := CreateComponentPath(cfg.Run)
	cmd := runtimeContext.Command(path)

	p.logger.Debugf("loading runtime '%s' plugin %s", cfg.Run.Runtime, cmd)
	client := goplugin.NewClient(&goplugin.ClientConfig{
		HandshakeConfig: sdk.Handshake,
		Plugins:         pluginSet,
		Cmd:             cmd,
		AllowedProtocols: []goplugin.Protocol{
			goplugin.ProtocolNetRPC,
			goplugin.ProtocolGRPC,
		},
	})
	clientProtocol, err := client.Client()
	if err != nil {
		return err
	}
	p.clientProtocol = clientProtocol
	return nil
}

func (c *Plugin) Store() (state.Store, error) {
	name := string(sdk.ProtocolGRPC)
	value, err := c.clientProtocol.Dispense(name)
	if err != nil {
		return nil, err
	}
	store, ok := value.(state.Store)
	if !ok {
		return nil, fmt.Errorf("expected %s to be state.Store", name)
	}
	return store, nil
}

func (c *Plugin) PubSub() (pubsub.PubSub, error) {
	name := string(sdk.ProtocolGRPC)
	value, err := c.clientProtocol.Dispense(name)
	if err != nil {
		return nil, err
	}
	store, ok := value.(pubsub.PubSub)
	if !ok {
		return nil, fmt.Errorf("expected %s to be pubsub.PubSub", name)
	}
	return store, nil
}

func (c *Plugin) Close() error {
	return c.clientProtocol.Close()
}

func CreatePluginSet(c *plugin.Config) goplugin.PluginSet {
	pluginSet := goplugin.PluginSet{}
	for _, c := range c.Components {
		switch c.ComponentType {
		case "state":
			pluginSet[c.Name] = state_sdk.GRPCStatePlugin{}
		}
	}
	return pluginSet
}

func CreateComponentPath(c *plugin.Run) string {
	fileName := fmt.Sprintf("dapr-%s-%s", c.Name, c.Version)
	return filepath.Join(c.BaseDirectory, c.Name, c.Version, fileName)
}
