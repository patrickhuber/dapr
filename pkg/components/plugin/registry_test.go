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
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dapr/dapr/pkg/modes"
	"github.com/dapr/dapr/pkg/plugin"
	daprt "github.com/dapr/dapr/pkg/testing"
)

func TestCreatePlugin(t *testing.T) {
	testRegistry := NewRegistry()

	t.Run("plugin is registered", func(t *testing.T) {
		const (
			pluginName    = "mockPlugin"
			pluginNameV2  = "mockPlugin/v2"
			componentType = "store"
			componentName = componentType + "." + pluginName
		)

		// Initiate mock object
		mockPlugin := new(daprt.MockPlugin)
		mockPluginV2 := new(daprt.MockPlugin)

		// act
		testRegistry.Register(modes.KubernetesMode, New(modes.KubernetesMode, func(cfg plugin.Config) (plugin.Plugin, error) {
			return mockPlugin, nil
		}))
		testRegistry.Register(modes.KubernetesMode, New(modes.KubernetesMode, func(cfg plugin.Config) (plugin.Plugin, error) {
			return mockPluginV2, nil
		}))

		// assert v0 and v1
		p, e := testRegistry.Create(
			plugin.Config{
				Name:    componentName,
				Version: pluginName,
				Type:    "v0",
			})
		assert.NoError(t, e)
		assert.Same(t, mockPlugin, p)

		p, e = testRegistry.Create(
			plugin.Config{
				Name:    componentName,
				Version: pluginName,
				Type:    "v1",
			})
		assert.NoError(t, e)
		assert.Same(t, mockPlugin, p)

		// assert v2
		pV2, e := testRegistry.Create(plugin.Config{
			Name:    "string",
			Version: "string",
			Type:    "string",
		})
		assert.NoError(t, e)
		assert.Same(t, mockPluginV2, pV2)

		// check case-insensitivity
		pV2, e = testRegistry.Create(
			plugin.Config{
				Name:    strings.ToUpper(componentName),
				Version: pluginName,
				Type:    "V2",
			})
		assert.NoError(t, e)
		assert.Same(t, mockPluginV2, pV2)
	})
}
