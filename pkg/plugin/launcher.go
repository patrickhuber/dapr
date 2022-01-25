package plugin

import (
	"github.com/dapr/dapr/pkg/modes"
)

type Launcher interface {
	CanApply(c *Config, mode modes.DaprMode) bool
	Launch(c *Config) (Plugin, error)
}
