/*
Copyright 2021 The Dapr Authors
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package plugin

import (
	"github.com/dapr/dapr/pkg/modes"
	"github.com/dapr/dapr/pkg/plugin"
)

type (
	FactoryMethod func(plugin.Config) (plugin.Plugin, error)

	// Plugin is a plugin component definition.
	Plugin struct {
		Mode          modes.DaprMode
		FactoryMethod FactoryMethod
	}

	// Registry is the interface for callers to get registered plugin components.
	Registry interface {
		Register(mode modes.DaprMode, components ...Plugin)
		Create(plugin.Config) (plugin.Plugin, error)
	}

	pluginRegistry struct {
		factory FactoryMethod
	}
)

// New creates a Plugin Factory.
func New(mode modes.DaprMode, factoryMethod FactoryMethod) Plugin {
	return Plugin{
		FactoryMethod: factoryMethod,
		Mode:          mode,
	}
}

// NewRegistry returns a new pub sub registry.
func NewRegistry() Registry {
	return &pluginRegistry{}
}

// Register attempts to match the given runtime mode to the list of components
// See the WithPlugin methods in the runtime to see how plugin factories are bound
func (p *pluginRegistry) Register(mode modes.DaprMode, components ...Plugin) {
	for _, component := range components {
		if component.Mode == mode {
			p.factory = component.FactoryMethod
			return
		}
	}
}

// Creates an instance of the plugin
func (p *pluginRegistry) Create(cfg plugin.Config) (plugin.Plugin, error) {
	return p.factory(cfg)
}
