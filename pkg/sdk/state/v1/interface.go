package state

import (
	"context"
	"net/rpc"

	"github.com/dapr/components-contrib/state"
	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"

	proto "github.com/dapr/dapr/pkg/proto/state/v1"
)

const (
	ProtocolRPC  = "state"
	ProtocolGRPC = "state_grpc"
)

var PluginMap = plugin.PluginSet{
	ProtocolGRPC: &GRPCStatePlugin{},
	ProtocolRPC:  &RPCStatePlugin{},
}

func CreatePluginMap(store state.Store) map[string]plugin.Plugin {
	return map[string]plugin.Plugin{
		ProtocolGRPC: &GRPCStatePlugin{
			Impl: store,
		},
		ProtocolRPC: &RPCStatePlugin{
			Impl: store,
		},
	}
}

type GRPCStatePlugin struct {
	plugin.Plugin
	Impl state.Store
}

func (p *GRPCStatePlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	proto.RegisterStoreServer(s, &GRPCServer{Impl: p.Impl})
	return nil
}

func (p *GRPCStatePlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &GRPCClient{client: proto.NewStoreClient(c)}, nil
}

type RPCStatePlugin struct {
	Impl state.Store
}

func (p *RPCStatePlugin) Server(b *plugin.MuxBroker) (interface{}, error) {
	return &RPCServer{
		Impl: p.Impl,
	}, nil
}

func (p *RPCStatePlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &RPCClient{
		client: c,
	}, nil
}
