package plugin

import (
	"github.com/dapr/components-contrib/pubsub"
	"github.com/dapr/components-contrib/state"
)

type Client interface {
	Store(name string) (state.Store, error)
	PubSub(name string) (pubsub.PubSub, error)
}
