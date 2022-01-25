package plugin

import (
	"fmt"

	"github.com/dapr/dapr/pkg/modes"
)

type LauncherFactory interface {
	Create(*Config, modes.DaprMode) (Launcher, error)
}

type launcherFactory struct {
	launchers []Launcher
}

func NewLauncherFactory(launchers ...Launcher) LauncherFactory {
	return &launcherFactory{
		launchers: launchers,
	}
}

func (f *launcherFactory) Create(c *Config, m modes.DaprMode) (Launcher, error) {
	for _, l := range f.launchers {
		if l.CanApply(c, m) {
			return l, nil
		}
	}
	return nil, fmt.Errorf("unable to find a launcher for the given config and mode")
}
