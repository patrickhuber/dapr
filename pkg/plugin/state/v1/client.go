package state

import (
	"strings"

	"github.com/dapr/components-contrib/state"
	"github.com/dapr/dapr/pkg/plugin"
)

type Client struct {
	internal state.Store
}

func (h *Client) Features() []state.Feature {
	return h.internal.Features()
}

func (h *Client) Init(metadata state.Metadata) error {

	// call the plugin, filter out any host metadata that isn't needed by the downstream
	return h.internal.Init(h.filterMetadata(metadata))
}

// filterMetadata removes any "plugin." prefixed parameters
func (h *Client) filterMetadata(metadata state.Metadata) state.Metadata {
	properties := map[string]string{}
	for k, v := range metadata.Properties {
		if strings.HasPrefix(k, plugin.MetadataPrefix) {
			continue
		}
		properties[k] = v
	}
	return state.Metadata{
		Properties: properties,
	}
}

func (h *Client) Get(req *state.GetRequest) (*state.GetResponse, error) {
	return h.internal.Get(req)
}

func (h *Client) Set(req *state.SetRequest) error {
	return h.internal.Set(req)
}

func (h *Client) Ping() error {
	return h.internal.Ping()
}

func (h *Client) Delete(req *state.DeleteRequest) error {
	return h.internal.Delete(req)
}

func (h *Client) BulkDelete(req []state.DeleteRequest) error {
	return h.internal.BulkDelete(req)
}

func (h *Client) BulkGet(req []state.GetRequest) (bool, []state.BulkGetResponse, error) {
	return h.internal.BulkGet(req)
}

func (h *Client) BulkSet(req []state.SetRequest) error {
	return h.internal.BulkSet(req)
}
