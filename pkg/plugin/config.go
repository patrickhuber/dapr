package plugin

import config "github.com/dapr/dapr/pkg/config/modes"

// Config defines the configuration for a plugin
type Config struct {
	Name       string
	Version    string
	Type       string
	Standalone config.StandaloneConfig
	Kubernetes config.KubernetesConfig
}
