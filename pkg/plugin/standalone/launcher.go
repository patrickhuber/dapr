package standalone

import (
	"github.com/dapr/dapr/pkg/modes"
	"github.com/dapr/dapr/pkg/plugin"
	"github.com/dapr/dapr/pkg/sdk"
	"github.com/dapr/dapr/pkg/sdk/state/v1"
	goplugin "github.com/hashicorp/go-plugin"

	pluginapi "github.com/dapr/dapr/pkg/apis/plugins/v1alpha1"
)

type Launcher struct {
}

func NewStandaloneLauncher(p *pluginapi.Plugin) plugin.Launcher {

	return &Launcher{}
}

func (l *Launcher) CanApply(p *pluginapi.Plugin, mode modes.DaprMode) bool {
	return p.Spec.Run != nil && mode == modes.StandaloneMode
}

func (l *Launcher) Launch(p *pluginapi.Plugin) (plugin.Client, error) {

	pluginSet := CreatePluginSet(p)

	runtimeContext := GetRuntimeContextFromString(p.Spec.Run.Runtime)
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
	return &Client{
		clientProtocol: clientProtocol,
	}, nil
}

func CreatePluginSet(p *pluginapi.Plugin) goplugin.PluginSet {
	pluginSet := goplugin.PluginSet{}
	for _, c := range p.Spec.Components {
		switch c.ComponentType {
		case "state":
			pluginSet[c.Name] = state.GRPCStatePlugin{}
		}
	}
	return pluginSet
}
