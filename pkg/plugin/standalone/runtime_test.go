package standalone_test

import (
	"testing"

	"github.com/dapr/dapr/pkg/plugin/standalone"
	"github.com/stretchr/testify/require"
)

func TestGetRuntimeContextCanGetDefault(t *testing.T) {
	runtimeContext := standalone.GetRuntimeContext(standalone.RuntimeDefault)
	require.Equal(t, standalone.RuntimeDefault, runtimeContext.Name())
}

func TestGetRuntimeContextCanGetDotnet(t *testing.T) {
	runtimeContext := standalone.GetRuntimeContext(standalone.RuntimeDotnet)
	require.Equal(t, standalone.RuntimeDotnet, runtimeContext.Name())
}

func TestGetRuntimeContextCanGetJava(t *testing.T) {
	runtimeContext := standalone.GetRuntimeContext(standalone.RuntimeJava)
	require.Equal(t, standalone.RuntimeJava, runtimeContext.Name())
}
