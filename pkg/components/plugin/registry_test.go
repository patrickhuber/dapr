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

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	"github.com/dapr/dapr/pkg/plugin"
	daprt "github.com/dapr/dapr/pkg/testing"
)

func TestCreateFullName(t *testing.T) {
	t.Run("create redis plugin key name", func(t *testing.T) {
		assert.Equal(t, "plugin.redis", createFullName("redis"))
	})

	t.Run("create kafka plugin key name", func(t *testing.T) {
		assert.Equal(t, "plugin.kafka", createFullName("kafka"))
	})

	t.Run("create azure service bus plugin key name", func(t *testing.T) {
		assert.Equal(t, "plugin.azure.servicebus", createFullName("azure.servicebus"))
	})

	t.Run("create rabbitmq plugin key name", func(t *testing.T) {
		assert.Equal(t, "plugin.rabbitmq", createFullName("rabbitmq"))
	})
}

func TestCreatePlugin(t *testing.T) {
	testRegistry := NewRegistry()

	t.Run("plugin messagebus is registered", func(t *testing.T) {
		const (
			pluginName    = "mockPlugin"
			pluginNameV2  = "mockPlugin/v2"
			componentName = "plugin." + pluginName
		)

		// Initiate mock object
		mockPlugin := new(daprt.MockPlugin)
		mockPluginV2 := new(daprt.MockPlugin)

		// act
		testRegistry.Register(New(pluginName, func() (plugin.Plugin, error) {
			return mockPlugin, nil
		}))
		testRegistry.Register(New(pluginNameV2, func() (plugin.Plugin, error) {
			return mockPluginV2, nil
		}))

		// assert v0 and v1
		p, e := testRegistry.Create(componentName, "v0")
		assert.NoError(t, e)
		assert.Same(t, mockPlugin, p)

		p, e = testRegistry.Create(componentName, "v1")
		assert.NoError(t, e)
		assert.Same(t, mockPlugin, p)

		// assert v2
		pV2, e := testRegistry.Create(componentName, "v2")
		assert.NoError(t, e)
		assert.Same(t, mockPluginV2, pV2)

		// check case-insensitivity
		pV2, e = testRegistry.Create(strings.ToUpper(componentName), "V2")
		assert.NoError(t, e)
		assert.Same(t, mockPluginV2, pV2)
	})

	t.Run("plugin messagebus is not registered", func(t *testing.T) {
		const PluginName = "fakePlugin"

		// act
		p, actualError := testRegistry.Create(createFullName(PluginName), "v1")
		expectedError := errors.Errorf("couldn't find message bus %s/v1", createFullName(PluginName))
		// assert
		assert.Nil(t, p)
		assert.Equal(t, expectedError.Error(), actualError.Error())
	})
}
