package standalone_test

import (
	"testing"
	"testing/fstest"

	"github.com/dapr/components-contrib/configuration"
	config "github.com/dapr/dapr/pkg/config/modes"
	"github.com/dapr/dapr/pkg/plugin"
	"github.com/dapr/dapr/pkg/plugin/standalone"
	"github.com/dapr/kit/logger"
	"github.com/stretchr/testify/require"
)

func TestInitCanInitialize(t *testing.T) {
	var mapFS = fstest.MapFS{
		"root/plugins/test/v1/dapr-test-v1": {},
	}
	m := configuration.Metadata{
		Properties: map[string]string{},
	}
	p := standalone.NewPlugin(logger.NewLogger("default"), plugin.Config{
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
	}, mapFS)
	err := p.Init(m)
	require.Nil(t, err)
}
