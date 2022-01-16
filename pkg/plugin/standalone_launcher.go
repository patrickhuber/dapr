package plugin

import (
	"github.com/dapr/dapr/pkg/sdk"
	"github.com/dapr/dapr/pkg/sdk/state/v1"
	"github.com/hashicorp/go-plugin"

	pluginapi "github.com/dapr/dapr/pkg/apis/plugins/v1alpha1"
)

type StandaloneLauncher struct {
	client *plugin.Client
}

func NewStandaloneLauncher(p *pluginapi.Plugin) Launcher {

	pluginSet := CreatePluginSet(p)

	runtimeContext := GetRuntimeContextFromString(p.Spec.Run.Runtime)
	cmd := runtimeContext.Command("")

	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: sdk.Handshake,
		Plugins:         pluginSet,
		Cmd:             cmd,
		AllowedProtocols: []plugin.Protocol{
			plugin.ProtocolNetRPC,
			plugin.ProtocolGRPC,
		},
	})
	return &StandaloneLauncher{
		client: client,
	}
}

func CreatePluginSet(p *pluginapi.Plugin) plugin.PluginSet {
	pluginSet := plugin.PluginSet{}
	for _, c := range p.Spec.Components {
		switch c.ComponentType {
		case "state":
			pluginSet[c.Name] = state.GRPCStatePlugin{}
		}
	}
	return pluginSet
}

func (l *StandaloneLauncher) Launch() (Client, error) {
	clientProtocol, err := l.client.Client()
	if err != nil {
		return nil, err
	}
	return &StandaloneClient{
		clientProtocol: clientProtocol,
	}, nil
}
