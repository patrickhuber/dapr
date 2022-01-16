package plugin

import (
	pluginapi "github.com/dapr/dapr/pkg/apis/plugins/v1alpha1"
)

type ContainerLauncher struct {
	container *pluginapi.Container
}

func NewContainerLauncher(p *pluginapi.Plugin) Launcher {

	return &ContainerLauncher{container: p.Spec.Container}
}

func (l *ContainerLauncher) Launch() (Client, error) {
	return nil, nil
}
