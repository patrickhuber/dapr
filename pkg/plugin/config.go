package plugin

import "github.com/dapr/components-contrib/configuration"

// Config defines the configuration for a plugin
type Config struct {
	Container  *Container  `json:"container"`
	Run        *Run        `json:"run"`
	Components []Component `json:"components"`
}

// Container defines the desired container for the plugin
type Container struct {
	Repository string `json:"repository"`
	Tag        string `json:"tag"`
}

// Run defines the desired run command for the plugin
type Run struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Runtime string `json:"runtime"`
}

// Component defines the name and component type of the component. This is used to register the component in the component map
type Component struct {
	Name          string `json:"name"`
	ComponentType string `json:"componentType"`
}

// MapComponentAPIToConfig maps the component API schema to the Configuration schema
func MapComponentAPIToConfig(m configuration.Metadata) *Config {
	return &Config{}
}
