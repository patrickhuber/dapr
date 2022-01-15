package plugin

import (
	"net"
	"os/exec"

	pluginapi "github.com/dapr/dapr/pkg/apis/plugins/v1alpha1"
	"github.com/dapr/dapr/pkg/modes"
	"github.com/dapr/dapr/pkg/sdk"
	"github.com/dapr/dapr/pkg/sdk/state/v1"
	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
)

func Launch(m modes.DaprMode, p pluginapi.Plugin) Client {
	pluginSet := CreatePluginSet(p)
	if p.Spec.Run != nil {
		return LaunchStandalone(p.Spec.Run, pluginSet)
	}

	if p.Spec.Container == nil {
		return nil
	}

	if m == modes.StandaloneMode {
		return LaunchContainer(p.Spec.Container, pluginSet)
	}
	return LaunchKubernetes(p.Spec.Container)
}

func LaunchContainer(container *pluginapi.Container, pluginSet plugin.PluginSet) Client {
	return nil
}

func LaunchKubernetes(container *pluginapi.Container) Client {
	return nil
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

func LaunchStandalone(run *pluginapi.Run, pluginSet plugin.PluginSet) Client {

	var cmd *exec.Cmd
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: sdk.Handshake,
		Plugins:         pluginSet,
		Cmd:             cmd,
		AllowedProtocols: []plugin.Protocol{
			plugin.ProtocolNetRPC,
			plugin.ProtocolGRPC,
		},
	})
	return &StandaloneClient{
		client: client,
	}
}

func CreatePluginSet(p pluginapi.Plugin) plugin.PluginSet {
	pluginSet := plugin.PluginSet{}
	for _, c := range p.Spec.Components {
		switch c.ComponentType {
		case "state":
			pluginSet[c.Name] = state.GRPCStatePlugin{}
		}
	}
	return pluginSet
}
