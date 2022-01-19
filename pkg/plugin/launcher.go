package plugin

import (
	pluginapi "github.com/dapr/dapr/pkg/apis/plugins/v1alpha1"
	"github.com/dapr/dapr/pkg/modes"
)

type Launcher interface {
	CanApply(p *pluginapi.Plugin, mode modes.DaprMode) bool
	Launch(p *pluginapi.Plugin) (Plugin, error)
}
