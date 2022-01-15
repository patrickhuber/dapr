package plugin

import (
	"net"

	"google.golang.org/grpc"
)


func NewGRPCClient(address net.Addr) (Client, error) {
	serverAddress := address.String()
	conn, err := grpc.Dial(serverAddress)
	if err != nil {
		return nil, err
	}
	return &GRPCClient{
		connection: conn,
	}, nil
}
