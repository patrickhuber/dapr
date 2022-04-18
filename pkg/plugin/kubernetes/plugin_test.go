package kubernetes_test

import (
	"context"
	"fmt"
	"log"
	"net"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"

	"github.com/dapr/components-contrib/configuration"
	"github.com/dapr/components-contrib/state"
	"github.com/dapr/dapr/pkg/env"
	"github.com/dapr/dapr/pkg/plugin"
	"github.com/dapr/dapr/pkg/plugin/kubernetes"
	stateproto "github.com/dapr/dapr/pkg/proto/state/v1"
	sdk_state "github.com/dapr/dapr/pkg/sdk/state/v1"
	"github.com/dapr/kit/logger"
	"github.com/stretchr/testify/require"
)

// creates the dialer function for initializing a grpc connection with a dial context
// see: http://www.inanzzz.com/index.php/post/w9qr/unit-testing-golang-grpc-client-and-server-application-with-bufconn-package
func dialer() func(ctx context.Context, s string) (net.Conn, error) {
	listener := bufconn.Listen(1024 * 1024)
	server := grpc.NewServer()
	store := &sdk_state.GRPCServer{
		Impl: plugin.NewMemoryStore(),
	}
	stateproto.RegisterStoreServer(server, store)
	go func() {
		if err := server.Serve(listener); err != nil {
			log.Fatal(err)
		}
	}()

	return func(ctx context.Context, s string) (net.Conn, error) {
		return listener.Dial()
	}
}

func MockConnectionFactory(metadata *kubernetes.Metadata) (*grpc.ClientConn, error) {
	ctx := context.Background()
	return grpc.DialContext(ctx, "", grpc.WithInsecure(), grpc.WithContextDialer(dialer()))
}

func TestStorePlugin(t *testing.T) {
	const ComponentName = "test"
	const ComponentVersion = "v1"
	logger := logger.NewLogger("test")
	cfg := plugin.Config{
		Name:    ComponentName,
		Version: ComponentVersion,
	}
	environment := env.NewMemory()
	environment.Set("DAPR_PLUGIN_TEST", fmt.Sprintf("name: %s|version: %s|ip: 192.168.1.1|port: 9999", ComponentName, ComponentVersion))
	discovery := kubernetes.NewDiscovery(environment)
	p := kubernetes.NewPlugin(logger, cfg, discovery, MockConnectionFactory)
	err := p.Init(configuration.Metadata{})
	require.Nil(t, err)

	t.Run("get empty returns empty", func(t *testing.T) {
		store, err := p.Store()
		require.Nil(t, err)
		require.NotNil(t, store)
		response, err := store.Get(&state.GetRequest{})
		require.Nil(t, err)
		require.Equal(t, &state.GetResponse{}, response)
	})
	t.Run("set allows roundtrip", func(t *testing.T) {
		const TestKey = "test"
		const TestData = "data"

		store, err := p.Store()
		require.Nil(t, err)
		require.NotNil(t, store)
		err = store.Set(&state.SetRequest{
			Key:   TestKey,
			Value: TestData,
		})
		require.Nil(t, err)
		response, err := store.Get(&state.GetRequest{
			Key: TestKey,
			Options: state.GetStateOption{
				Consistency: state.Strong,
			},
		})
		require.Nil(t, err)
		require.NotNil(t, response)
		require.Equal(t, TestData, string(response.Data))
	})
}
