package testing

import (
	"github.com/dapr/components-contrib/pubsub"
	"github.com/dapr/components-contrib/state"
)

type MockPlugin struct {
}

func (p *MockPlugin) Store(name string) (state.Store, error) {
	return nil, nil
}

func (p *MockPlugin) PubSub(name string) (pubsub.PubSub, error) {
	return nil, nil
}
