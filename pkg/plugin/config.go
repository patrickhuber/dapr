package plugin

import "github.com/dapr/components-contrib/configuration"

const ContainerRepositoryKey = "container.repository"
const ContainerTagKey = "container.tag"

const RunNameKey = "run.name"
const RunVersionKey = "run.version"
const RunRuntimeKey = "run.runtime"
const RunBaseDirectory = "run.baseDirectory"

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
	BaseDirectory string `json:"baseDirectory"`
	Name          string `json:"name"`
	Version       string `json:"version"`
	Runtime       string `json:"runtime"`
}

// Component defines the name and component type of the component. This is used to register the component in the component map
type Component struct {
	Name          string `json:"name"`
	ComponentType string `json:"componentType"`
}

// MapComponentAPIToConfig maps the component API schema to the Configuration schema
func MapComponentAPIToConfig(m configuration.Metadata) *Config {
	return &Config{
		Container: mapPropertiesToContainerConfig(m.Properties),
		Run:       mapPropertiesToRunConfig(m.Properties),
	}
}

func mapPropertiesToContainerConfig(properties map[string]string) *Container {
	if properties == nil {
		return nil
	}
	repository, exists := properties[ContainerRepositoryKey]
	if !exists {
		return nil
	}
	tag, exists := properties[ContainerTagKey]
	if !exists {
		return nil
	}
	return &Container{
		Repository: repository,
		Tag:        tag,
	}
}

func mapPropertiesToRunConfig(properties map[string]string) *Run {
	if properties == nil {
		return nil
	}
	name, exists := properties[RunNameKey]
	if !exists {
		return nil
	}
	version, exists := properties[RunVersionKey]
	if !exists {
		return nil
	}
	runtime, exists := properties[RunRuntimeKey]
	if !exists {
		return nil
	}
	baseDirectory, exists := properties[RunBaseDirectory]
	if !exists {
		return nil
	}
	return &Run{
		Name:          name,
		Version:       version,
		Runtime:       runtime,
		BaseDirectory: baseDirectory,
	}
}
