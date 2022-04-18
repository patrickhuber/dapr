package standalone_test

import (
	"fmt"
	"testing"
	"testing/fstest"

	"github.com/dapr/components-contrib/configuration"
	"github.com/dapr/components-contrib/state"
	config "github.com/dapr/dapr/pkg/config/modes"
	"github.com/dapr/dapr/pkg/plugin"
	"github.com/dapr/dapr/pkg/plugin/standalone"
	state_sdk "github.com/dapr/dapr/pkg/sdk/state/v1"
	"github.com/dapr/kit/logger"
	goplugin "github.com/hashicorp/go-plugin"
	"github.com/stretchr/testify/require"
)

type mockClientProtocol struct {
}

func (p *mockClientProtocol) Dispense(name string) (interface{}, error) {
	switch name {
	case state_sdk.ProtocolGRPC:
		return plugin.NewMemoryStore(), nil
	}
	return nil, fmt.Errorf("unrecognized service %s", name)
}

func (p *mockClientProtocol) Ping() error {
	return nil
}

func (p *mockClientProtocol) Close() error {
	return nil
}

func MockClientProtocolFactory(client *goplugin.Client) (goplugin.ClientProtocol, error) {
	return &mockClientProtocol{}, nil
}

func TestStorePlugin(t *testing.T) {
	var mapFS = fstest.MapFS{
		"root/plugins/test/v1/dapr-test-v1": {},
	}
	m := configuration.Metadata{
		Properties: map[string]string{},
	}
	p := standalone.NewPlugin(
		logger.NewLogger("default"),
		plugin.Config{
			Name:    "test",
			Version: "v1",
			Type:    "state",
			Standalone: config.StandaloneConfig{
				ComponentsPath: "string",
				PluginsPath:    "root/plugins",
			},
			Kubernetes: config.KubernetesConfig{
				ControlPlaneAddress: "string",
			},
		},
		mapFS,
		MockClientProtocolFactory)

	err := p.Init(m)
	require.Nil(t, err)

	store, err := p.Store()
	require.Nil(t, err)
	require.NotNil(t, store)

	t.Run("get returns empty response", func(t *testing.T) {
		const TestKey = "TEST"

		response, err := store.Get(&state.GetRequest{
			Key: TestKey,
		})
		require.Nil(t, err)
		require.NotNil(t, response)
		require.Equal(t, &state.GetResponse{}, response)
	})
	t.Run("get roundtrips set", func(t *testing.T) {
		const TestKey = "TEST"
		const TestValue = "data"

		err := store.Set(&state.SetRequest{
			Key:   TestKey,
			Value: TestValue,
		})
		require.Nil(t, err)

		response, err := store.Get(&state.GetRequest{
			Key: TestKey,
		})
		require.Nil(t, err)
		require.NotNil(t, response)
		require.Equal(t, TestValue, string(response.Data))
	})
}
