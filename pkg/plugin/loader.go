package plugin

import plugins_v1alpha1 "github.com/dapr/dapr/pkg/apis/plugins/v1alpha1"

// Loader loads the plugin definition
type Loader interface {
	Load() ([]*plugins_v1alpha1.Plugin, error)
}
