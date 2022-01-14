package state

import (
	"fmt"

	"github.com/dapr/components-contrib/state"
	sdk "github.com/dapr/dapr/pkg/sdk/state/v1"
	goplugin "github.com/hashicorp/go-plugin"
)

type StandaloneClient struct {
	Client
	client goplugin.Client
}

func NewStandaloneClient(client goplugin.Client) (state.Store, error) {
	rpcClient, err := client.Client()
	if err != nil {
		return nil, err
	}

	raw, err := rpcClient.Dispense(sdk.ProtocolGRPC)
	if err != nil {
		return nil, err
	}

	store, ok := raw.(state.Store)
	if !ok {
		return nil, fmt.Errorf("the plugin supplied is not a state store")
	}

	return &StandaloneClient{
		client: client,
		Client: Client{
			internal: store,
		},
	}, nil
}

func (h *StandaloneClient) Close() error {
	// this kills everything running under the plugin
	h.client.Kill()
	return nil
}
