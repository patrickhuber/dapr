package plugin

import (
	pluginapi "github.com/dapr/dapr/pkg/apis/plugins/v1alpha1"
	"github.com/dapr/dapr/pkg/modes"
	"github.com/dapr/dapr/pkg/plugin"
)

type Launcher struct {
}

func NewContainerLauncher(p *pluginapi.Plugin) plugin.Launcher {

	return &Launcher{}
}

func (l *Launcher) CanApply(c *plugin.Config, mode modes.DaprMode) bool {
	return mode == modes.KubernetesMode && c.Container != nil
}

func (l *Launcher) Launch(c *plugin.Config) (plugin.Plugin, error) {
	return nil, nil
}
