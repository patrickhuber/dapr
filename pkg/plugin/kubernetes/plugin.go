package kubernetes

import (
	"fmt"
	"net"
	"strings"

	"github.com/dapr/components-contrib/configuration"
	"github.com/dapr/components-contrib/pubsub"
	"github.com/dapr/components-contrib/state"
	"github.com/dapr/dapr/pkg/env"
	"github.com/dapr/dapr/pkg/plugin"
	stateproto "github.com/dapr/dapr/pkg/proto/state/v1"
	statesdk "github.com/dapr/dapr/pkg/sdk/state/v1"
	"github.com/dapr/kit/logger"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v2"
)

type Plugin struct {
	connection  *grpc.ClientConn
	cfg         plugin.Config
	logger      logger.Logger
	environment env.Env
}

type Metadata struct {
	Name    string `json:"name" yaml:"name"`
	Version string `json:"version" yaml:"version"`
	Address string `json:"address" yaml:"address"`
	Port    int    `json:"port" yaml:"port"`
}

func NewPlugin(logger logger.Logger, cfg plugin.Config, environment env.Env) plugin.Plugin {
	return &Plugin{
		logger:      logger,
		cfg:         cfg,
		environment: environment,
	}
}

func (p *Plugin) Init(m configuration.Metadata) error {

	// get the injected metadata from the dapr plugin environment variables
	metadata, err := p.getCurrentComponentPluginMetdata()
	if err != nil {
		return err
	}

	ipAddress := net.ParseIP(metadata.Address)
	if ipAddress == nil {
		return fmt.Errorf("%s is not a valid ip address", metadata.Address)
	}

	addressAndPort := fmt.Sprintf("%s:%d", metadata.Address, metadata.Port)

	conn, err := grpc.Dial(addressAndPort)
	if err != nil {
		return err
	}
	p.connection = conn
	return nil
}

func (p *Plugin) getPluginEnvironmentVariables() []string {
	all := p.environment.List()
	prefix := "DAPR_PLUGIN_"
	matched := make([]string, 0)
	for _, v := range all {
		if strings.HasPrefix(v, prefix) {
			matched = append(matched, v)
		}
	}
	return matched
}

func (p *Plugin) getCurrentComponentPluginMetdata() (Metadata, error) {
	metadataList, err := p.getPluginMetadata()
	var zero Metadata
	if err != nil {
		return zero, err
	}
	for _, metadata := range metadataList {
		if metadata.Name != p.cfg.Name {
			continue
		}
		if metadata.Version != p.cfg.Version {
			continue
		}
		return metadata, nil
	}
	return zero, fmt.Errorf("metada for component %s version %s not found in injected environment variables", p.cfg.Name, p.cfg.Type)
}

// getPluginMetadata gets the plugin metadata as an encoded yaml string and turnes it into a structure
//
// input:
// 	DAPR_PLUGIN_GOMEMORY_V1=name: abc|version: v1|ip: 192.168.1.1|port: 1234
//
// to yaml:
// 	name: abc123
// 	version: v1
// 	ip: 192.168.1.1
// 	port: 1234
//
// to metadata
// allMetadata := []Metadata{
//    {
// 			Name: "abc123",
//			Version: "v1",
// 			IP: "192.168.1.1",
// 			Port: "1234",
//    },
// }
func (p *Plugin) getPluginMetadata() ([]Metadata, error) {
	names := p.getPluginEnvironmentVariables()
	allMetadata := []Metadata{}
	for _, name := range names {
		value, ok := p.environment.Lookup(name)
		if !ok {
			continue
		}
		metadata := Metadata{}

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
		if err != nil {
			return nil, err
		}

		// append to the list of all metadata
		allMetadata = append(allMetadata, metadata)
	}
	return allMetadata, nil
}

func (p *Plugin) Store() (state.Store, error) {
	client := stateproto.NewStoreClient(p.connection)
	return statesdk.NewGRPCClient(client), nil
}

func (p *Plugin) PubSub() (pubsub.PubSub, error) {
	return nil, nil
}
