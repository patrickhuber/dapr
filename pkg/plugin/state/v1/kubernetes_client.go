package state

import (
	"net"

	"github.com/dapr/components-contrib/state"
	proto "github.com/dapr/dapr/pkg/proto/state/v1"
	sdk "github.com/dapr/dapr/pkg/sdk/state/v1"
	"google.golang.org/grpc"
)

type KubernetesClient struct {
	Client
	conn *grpc.ClientConn
}

func NewKubernetesClient(address net.Addr) (state.Store, error) {
	serverAddress := address.String()
	conn, err := grpc.Dial(serverAddress)
	if err != nil {
		return nil, err
	}
	storeClient := proto.NewStoreClient(conn)
	grpcClient := sdk.NewGRPCClient(storeClient)

	// multiplexing?
	return &KubernetesClient{
		conn: conn,
		Client: Client{
			internal: grpcClient,
		},
	}, nil
}

func (c *KubernetesClient) Close() error {
	return c.conn.Close()
}
