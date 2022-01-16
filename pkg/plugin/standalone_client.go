package plugin

import (
	"fmt"

	"github.com/dapr/components-contrib/pubsub"
	"github.com/dapr/components-contrib/state"
	"github.com/hashicorp/go-plugin"
)

type StandaloneClient struct {
	clientProtocol plugin.ClientProtocol
}

func (c *StandaloneClient) Store(name string) (state.Store, error) {
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

func (s *StandaloneClient) PubSub(name string) (pubsub.PubSub, error) {
	return nil, nil
}
