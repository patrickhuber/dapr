package plugin_test

import (
	"testing"

	"github.com/dapr/components-contrib/configuration"
	"github.com/dapr/dapr/pkg/plugin"
	"github.com/stretchr/testify/require"
)

func TestConfigCanMapContainer(t *testing.T) {
	cfg := plugin.MapComponentAPIToConfig(configuration.Metadata{
		Properties: map[string]string{
			plugin.ContainerRepositoryKey: "repository",
			plugin.ContainerTagKey:        "tag",
		},
	})
	require.NotNil(t, cfg)
	require.NotNil(t, cfg.Container)
	require.NotEmpty(t, cfg.Container.Repository)
	require.NotEmpty(t, cfg.Container.Tag)
}

func TestConfigCanMapRun(t *testing.T) {
	cfg := plugin.MapComponentAPIToConfig(configuration.Metadata{
		Properties: map[string]string{
			plugin.RunNameKey:       "name",
			plugin.RunVersionKey:    "version",
			plugin.RunRuntimeKey:    "runtime",
			plugin.RunBaseDirectory: "baseDirectory",
		},
	})
	require.NotNil(t, cfg)
	require.NotNil(t, cfg.Run)
	require.NotEmpty(t, cfg.Run.BaseDirectory)
	require.NotEmpty(t, cfg.Run.Name)
	require.NotEmpty(t, cfg.Run.Runtime)
	require.NotEmpty(t, cfg.Run.Version)
}
