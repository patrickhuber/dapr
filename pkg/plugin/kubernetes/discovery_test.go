package kubernetes_test

import (
	"testing"

	"github.com/dapr/dapr/pkg/env"
	"github.com/dapr/dapr/pkg/plugin/kubernetes"
	"github.com/stretchr/testify/require"
)

func TestLookup(t *testing.T) {
	environment := env.NewMemory()
	environment.Set("DAPR_PLUGIN_TEST", "name: test|version: v1|ip: 192.168.1.1|port: 9999")
	discovery := kubernetes.NewDiscovery(environment)
	metadata, ok, err := discovery.Lookup("test", "v1")
	require.Nil(t, err)
	require.True(t, ok)
	require.NotNil(t, metadata)
}
