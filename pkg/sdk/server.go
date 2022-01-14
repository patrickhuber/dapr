package sdk

import (
	"github.com/hashicorp/go-plugin"
)

type Server interface {
	// Server should return the RPC server compatible struct to serve
	// the methods that the Client calls over net/rpc.
	Server(*plugin.MuxBroker) (interface{}, error)
}

func Serve(pluginSet plugin.PluginSet) {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: Handshake,
		Plugins:         pluginSet,
		GRPCServer:      plugin.DefaultGRPCServer,
	})
}
