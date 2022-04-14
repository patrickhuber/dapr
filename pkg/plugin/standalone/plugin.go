package standalone

import (
	"fmt"
	"io/fs"
	"path/filepath"

	"github.com/dapr/components-contrib/configuration"
	"github.com/dapr/components-contrib/pubsub"
	"github.com/dapr/components-contrib/state"
	"github.com/dapr/dapr/pkg/plugin"
	"github.com/dapr/dapr/pkg/sdk"
	"github.com/dapr/kit/logger"
	goplugin "github.com/hashicorp/go-plugin"

	state_sdk "github.com/dapr/dapr/pkg/sdk/state/v1"
)

type Plugin struct {
	clientProtocol goplugin.ClientProtocol
	cfg            plugin.Config
	logger         logger.Logger
	filesystem     fs.FS
}

const BaseDirectoryKey = "standalone.BaseDirectory"

func NewPlugin(logger logger.Logger, cfg plugin.Config, filesystem fs.FS) plugin.Plugin {
	p := &Plugin{
		logger:     logger,
		cfg:        cfg,
		filesystem: filesystem,
	}
	return p
}

func (p *Plugin) Init(m configuration.Metadata) error {
	pluginPath, err := p.GetPluginPath()
	if err != nil {
		return err
	}

	runtimeContext, err := p.MatchRuntimeContext(pluginPath)
	if err != nil {
		return err
	}

	// enumerate the files in the plugin directory
	cmd := runtimeContext.Command(pluginPath)

	// additional maps would be joined to this
	pluginMap := goplugin.PluginSet{}
	for k, v := range state_sdk.PluginMap {
		pluginMap[k] = v
	}

	p.logger.Debugf("loading runtime '%s' plugin %s", runtimeContext.Name(), cmd)
	client := goplugin.NewClient(&goplugin.ClientConfig{
		HandshakeConfig: sdk.Handshake,
		Plugins:         pluginMap,
		Cmd:             cmd,
		AllowedProtocols: []goplugin.Protocol{
			goplugin.ProtocolGRPC,
		},
	})
	clientProtocol, err := client.Client()
	if err != nil {
		return err
	}
	p.clientProtocol = clientProtocol
	return nil
}

func (c *Plugin) Store() (state.Store, error) {
	name := string(state_sdk.ProtocolGRPC)
	value, err := c.clientProtocol.Dispense(name)
	if err != nil {
		return nil, err
	}
	store, ok := value.(state.Store)
	if !ok {
		return nil, fmt.Errorf("expected %s to be state.Store", name)
	}
	return store, nil
}

func (c *Plugin) PubSub() (pubsub.PubSub, error) {
	name := string(sdk.ProtocolGRPC)
	value, err := c.clientProtocol.Dispense(name)
	if err != nil {
		return nil, err
	}
	store, ok := value.(pubsub.PubSub)
	if !ok {
		return nil, fmt.Errorf("expected %s to be pubsub.PubSub", name)
	}
	return store, nil
}

func (c *Plugin) Close() error {
	return c.clientProtocol.Close()
}

func (p *Plugin) CreatePluginWildcardPath() string {
	fileName := fmt.Sprintf("dapr-%s-%s*", p.cfg.Name, p.cfg.Version)
	return filepath.Join(p.cfg.Standalone.PluginsPath, p.cfg.Name, p.cfg.Version, fileName)
}

func (p *Plugin) GetPluginPath() (string, error) {
	// create the plugin path
	pluginWildcardPath := p.CreatePluginWildcardPath()

	// look in the path for anything that matches the plugin file spec
	entries, err := fs.Glob(p.filesystem, pluginWildcardPath)
	if err != nil {
		return "", err
	}

	if len(entries) != 1 {
		return "", fmt.Errorf("found (%d) entries that match the path spec %s. expected one", len(entries), pluginWildcardPath)
	}
	return entries[0], nil
}

func (p *Plugin) MatchRuntimeContext(pluginPath string) (RuntimeContext, error) {
	runtimeContexts := MatchRuntimeContext(pluginPath)
	if len(runtimeContexts) != 1 {
		return nil, fmt.Errorf("found (%d) runtime contexts that match the file path %s", len(runtimeContexts), pluginPath)
	}
	return runtimeContexts[0], nil
}
