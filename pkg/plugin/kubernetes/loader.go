package kubernetes

import (
	plugins_v1alpha1 "github.com/dapr/dapr/pkg/apis/plugins/v1alpha1"
	"github.com/dapr/dapr/pkg/plugin"
)

type Loader struct {
}

func NewLoader() plugin.Loader {
	return &Loader{}
}

func (l *Loader) Load() ([]*plugins_v1alpha1.Plugin, error) {
	return nil, nil
}
