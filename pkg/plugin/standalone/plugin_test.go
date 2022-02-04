package standalone_test

import (
	"testing"

	"github.com/dapr/components-contrib/configuration"
	"github.com/dapr/dapr/pkg/plugin"
	"github.com/dapr/dapr/pkg/plugin/standalone"
	"github.com/dapr/kit/logger"
	"github.com/stretchr/testify/require"
)

func TestInitCanInitialize(t *testing.T) {
	m := configuration.Metadata{
		Properties: map[string]string{
			plugin.RunNameKey:       "gomemory",
			plugin.RunBaseDirectory: "",
			plugin.RunRuntimeKey:    string(standalone.RuntimeExec),
			plugin.RunVersionKey:    "v1",
		},
	}
	p := standalone.NewPlugin(logger.NewLogger("default"))
	err := p.Init(m)
	require.Nil(t, err)
}
