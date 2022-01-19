package kubernetes

import (
	pluginapi "github.com/dapr/dapr/pkg/apis/plugins/v1alpha1"
	"github.com/dapr/dapr/pkg/modes"
	"github.com/dapr/dapr/pkg/plugin"
)

type Launcher struct {
}

func NewKubernetesLauncher(p *pluginapi.Plugin) plugin.Launcher {
	return &Launcher{}
}

func (l *Launcher) CanApply(p *pluginapi.Plugin, mode modes.DaprMode) bool {
	return mode == modes.KubernetesMode && p.Spec.Container != nil
}

func (l *Launcher) Launch(p *pluginapi.Plugin) (plugin.Plugin, error) {
	return nil, nil
}
