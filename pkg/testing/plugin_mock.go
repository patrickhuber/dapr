package testing

import (
	"github.com/dapr/components-contrib/configuration"
	"github.com/dapr/components-contrib/pubsub"
	"github.com/dapr/components-contrib/state"
)

type MockPlugin struct {
}

func (p *MockPlugin) Init(metadata configuration.Metadata) error {
	return nil
}

func (p *MockPlugin) Store() (state.Store, error) {
	return nil, nil
}

func (p *MockPlugin) PubSub() (pubsub.PubSub, error) {
	return nil, nil
}
