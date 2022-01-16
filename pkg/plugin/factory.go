package plugin

import (
	"fmt"

	pluginapi "github.com/dapr/dapr/pkg/apis/plugins/v1alpha1"
	"github.com/dapr/dapr/pkg/modes"
)

func CreateLauncher(m modes.DaprMode, p *pluginapi.Plugin) (Launcher, error) {

	if p.Spec.Run != nil {
		return NewStandaloneLauncher(p), nil
	}

	if p.Spec.Container == nil {
		return nil, fmt.Errorf("plugin api must specify one of plugin.spec.run or plugin.spec.container")
	}

	if m == modes.StandaloneMode {
		return NewContainerLauncher(p), nil
	}
	return NewKubernetesLauncher(p), nil
}
