package plugin

import (
	"fmt"

	"github.com/dapr/components-contrib/configuration"
	"github.com/dapr/components-contrib/pubsub"
	"github.com/dapr/components-contrib/state"
)

// ErrComponentNotImplemented defines a not found error
var ErrComponentNotImplemented = fmt.Errorf("the plugin component service is not implemented")

type Plugin interface {
	// Init is called after the plugin is initialized with the metadata loaded from the component CRD
	Init(metadata configuration.Metadata) error
	// Store returns the state store served by this plugin. If the component is not implemented, ErrComponentNotImplemented is returned
	Store() (state.Store, error)
	// PubSub returns the pubsub service served by this plugin. If the component is not implemented, ErrComponentNotImplemented is returned
	PubSub() (pubsub.PubSub, error)
}
