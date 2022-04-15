package kubernetes

import (
	"fmt"
	"strings"

	"github.com/dapr/dapr/pkg/env"
	"gopkg.in/yaml.v2"
)

type Discovery interface {
	Lookup(name, version string) (*Metadata, bool, error)
}

const DaprPluginPrefix = "DAPR_PLUGIN_"

type discovery struct {
	environment env.Env
}

func NewDiscovery(environment env.Env) Discovery {
	return &discovery{
		environment: environment,
	}
}

func (d *discovery) Lookup(name, version string) (*Metadata, bool, error) {
	filtered := d.filter(d.environment.List(), DaprPluginPrefix)
	all, err := d.transformAll(filtered)
	if err != nil {
		return nil, false, err
	}
	return d.find(all, name, version)
}

func (d *discovery) find(all []*Metadata, name, version string) (*Metadata, bool, error) {
	for _, item := range all {
		if name == item.Name && version == item.Version {
			return item, true, nil
		}
	}
	return nil, false, nil
}

func (d *discovery) filter(list map[string]string, prefix string) map[string]string {
	matched := map[string]string{}
	for k, v := range list {
		if !strings.HasPrefix(k, prefix) {
			continue
		}
		matched[k] = v
	}
	return matched
}

func (d *discovery) transformAll(variables map[string]string) ([]*Metadata, error) {
	var all []*Metadata
	for k, v := range variables {
		metadata, err := d.transform(k, v)
		if err != nil {
			return nil, err
		}
		all = append(all, metadata)
	}
	return all, nil
}

func (d *discovery) transform(key, value string) (*Metadata, error) {
	metadata := &Metadata{}
	// first split the string by pipe and reset the value variable
	splits := strings.Split(value, "|")
	value = ""
	// create a yaml string by joining the segments into lines
	for i, s := range splits {
		if i > 0 {
			value += fmt.Sprintln()
		}
		value = fmt.Sprint(s)
	}
	// unmarshal the yaml string to the metadata variable
	err := yaml.Unmarshal([]byte(value), &metadata)

	return metadata, err
}
