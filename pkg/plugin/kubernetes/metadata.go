package kubernetes

type Metadata struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
	Address string `yaml:"address"`
	Port    int    `yaml:"port"`
}
