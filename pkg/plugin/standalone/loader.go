package standalone

import (
	"io/ioutil"
	"strings"

	plugins_v1alpha1 "github.com/dapr/dapr/pkg/apis/plugins/v1alpha1"
	config "github.com/dapr/dapr/pkg/config/modes"
	"github.com/dapr/dapr/pkg/plugin"
	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

type Loader struct {
	config *config.StandaloneConfig
}

func NewLoader(cfg *config.StandaloneConfig) plugin.Loader {
	return &Loader{
		config: cfg,
	}
}

func (l *Loader) Load() ([]*plugins_v1alpha1.Plugin, error) {
	filePath := l.config.PluginsPath
	input, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	splitInput := strings.Split(string(input), "---")

	scheme := runtime.NewScheme()
	err = apps.AddToScheme(scheme)
	if err != nil {
		return nil, err
	}
	err = core.AddToScheme(scheme)
	if err != nil {
		return nil, err
	}

	factory := serializer.NewCodecFactory(scheme)
	decoder := factory.UniversalDeserializer()
	objs := make([]interface{}, 0)

	for _, input := range splitInput {
		obj, _, err := decoder.Decode([]byte(input), nil, nil)
		if err != nil {
			return nil, err
		}
		objs = append(objs, obj)
	}

	plugins := []*plugins_v1alpha1.Plugin{}
	for _, o := range objs {
		switch t := o.(type) {
		case (*plugins_v1alpha1.Plugin):
			plugins = append(plugins, t)
		}
	}
	return plugins, nil
}
