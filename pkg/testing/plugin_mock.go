package testing

import (
	"strings"

	"github.com/dapr/components-contrib/configuration"
	"github.com/dapr/components-contrib/pubsub"
	"github.com/dapr/components-contrib/state"
)

type MockPlugin struct {
}

func (p *MockPlugin) Name() string {
	return "name"
}

func (p *MockPlugin) Type() string {
	return "type"
}

func (p *MockPlugin) Version() string {
	return "version"
}

func (p *MockPlugin) FullName() string {
	return strings.ToLower(strings.Join([]string{p.Type(), p.Name()}, "."))
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
