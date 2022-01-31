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
	"strings"

	"github.com/dapr/dapr/pkg/components"
	"github.com/dapr/dapr/pkg/plugin"
	"github.com/pkg/errors"
)

type (
	FactoryMethod func() (plugin.Plugin, error)

	// Plugin is a plugin component definition.
	Plugin struct {
		Name          string
		FactoryMethod FactoryMethod
	}

	// Registry is the interface for callers to get registered plugin components.
	Registry interface {
		Register(components ...Plugin)
		Create(name, version string) (plugin.Plugin, error)
	}

	pluginRegistry struct {
		plugins map[string]FactoryMethod
	}
)

// New creates a Plugin.
func New(name string, factoryMethod FactoryMethod) Plugin {
	return Plugin{
		Name:          name,
		FactoryMethod: factoryMethod,
	}
}

// NewRegistry returns a new pub sub registry.
func NewRegistry() Registry {
	return &pluginRegistry{
		plugins: map[string]FactoryMethod{},
	}
}

// Register registers one or more new message buses.
func (p *pluginRegistry) Register(components ...Plugin) {
	for _, component := range components {
		p.plugins[createFullName(component.Name)] = component.FactoryMethod
	}
}

// Create instantiates a pub/sub based on `name`.
func (p *pluginRegistry) Create(name, version string) (plugin.Plugin, error) {
	if method, ok := p.getplugin(name, version); ok {
		return method()
	}
	return nil, errors.Errorf("couldn't find message bus %s/%s", name, version)
}

func (p *pluginRegistry) getplugin(name, version string) (FactoryMethod, bool) {
	nameLower := strings.ToLower(name)
	versionLower := strings.ToLower(version)
	pluginFn, ok := p.plugins[nameLower+"/"+versionLower]
	if ok {
		return pluginFn, true
	}
	if components.IsInitialVersion(versionLower) {
		pluginFn, ok = p.plugins[nameLower]
	}
	return pluginFn, ok
}

func createFullName(name string) string {
	return strings.ToLower("plugin." + name)
}
