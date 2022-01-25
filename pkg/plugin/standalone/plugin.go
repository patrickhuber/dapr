package standalone

import (
	"fmt"

	"github.com/dapr/components-contrib/configuration"
	"github.com/dapr/components-contrib/pubsub"
	"github.com/dapr/components-contrib/state"
	"github.com/dapr/dapr/pkg/sdk"
	goplugin "github.com/hashicorp/go-plugin"
)

type Plugin struct {
	clientProtocol goplugin.ClientProtocol
}

func (c *Plugin) Init(m configuration.Metadata) error {
	return nil
}

func (c *Plugin) Store() (state.Store, error) {
	name := string(sdk.ProtocolGRPC)
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
