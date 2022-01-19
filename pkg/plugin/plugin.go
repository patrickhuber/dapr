package plugin

import (
	"github.com/dapr/components-contrib/pubsub"
	"github.com/dapr/components-contrib/state"
)

type Plugin interface {
	Store() (state.Store, error)
	PubSub() (pubsub.PubSub, error)
}
