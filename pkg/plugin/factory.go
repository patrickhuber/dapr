package plugin

import (
	"fmt"

	pluginapi "github.com/dapr/dapr/pkg/apis/plugins/v1alpha1"
	"github.com/dapr/dapr/pkg/modes"
)

type LauncherFactory interface {
	Create(*pluginapi.Plugin, modes.DaprMode) (Launcher, error)
}

type launcherFactory struct {
	launchers []Launcher
}

func NewLauncherFactory(launchers ...Launcher) LauncherFactory {
	return &launcherFactory{
		launchers: launchers,
	}
}

func (f *launcherFactory) Create(p *pluginapi.Plugin, m modes.DaprMode) (Launcher, error) {
	for _, l := range f.launchers {
		if l.CanApply(p, m) {
			return l, nil
		}
	}
	return nil, fmt.Errorf("unable to find a launcher for the given config and mode")
}
