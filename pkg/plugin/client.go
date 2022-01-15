package plugin

import (
	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
)

type Client interface {
}

type GRPCClient struct {
	connection *grpc.ClientConn
}

type StandaloneClient struct {
	client *plugin.Client
}
