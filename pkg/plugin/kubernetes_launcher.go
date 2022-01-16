package plugin

import (
	pluginapi "github.com/dapr/dapr/pkg/apis/plugins/v1alpha1"
)

type KubernetesLauncher struct {
	container *pluginapi.Container
}

func NewKubernetesLauncher(p *pluginapi.Plugin) Launcher {
	return &KubernetesLauncher{
		container: p.Spec.Container,
	}
}

func (l *KubernetesLauncher) Launch() (Client, error) {
	return nil, nil
}
