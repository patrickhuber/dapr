package standalone

import (
	"github.com/dapr/dapr/pkg/modes"
	"github.com/dapr/dapr/pkg/plugin"
	"github.com/dapr/dapr/pkg/sdk"
	"github.com/dapr/dapr/pkg/sdk/state/v1"
	goplugin "github.com/hashicorp/go-plugin"
)

type Launcher struct {
}

func NewStandaloneLauncher() plugin.Launcher {

	return &Launcher{}
}

func (l *Launcher) CanApply(c *plugin.Config, mode modes.DaprMode) bool {
	return c.Run != nil && mode == modes.StandaloneMode
}

func (l *Launcher) Launch(c *plugin.Config) (plugin.Plugin, error) {

	pluginSet := CreatePluginSet(c)

	runtimeContext := GetRuntimeContextFromString(c.Run.Runtime)
	cmd := runtimeContext.Command("")

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
		return nil, err
	}
	return &Plugin{
		clientProtocol: clientProtocol,
	}, nil
}

func CreatePluginSet(c *plugin.Config) goplugin.PluginSet {
	pluginSet := goplugin.PluginSet{}
	for _, c := range c.Components {
		switch c.ComponentType {
		case "state":
			pluginSet[c.Name] = state.GRPCStatePlugin{}
		}
	}
	return pluginSet
}
