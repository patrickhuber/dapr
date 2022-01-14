package state

import (
	"net/rpc"

	"github.com/dapr/components-contrib/state"
)

const (
	GetMethod      = "Plugin.Get"
	SetMethod      = "Plugin.Set"
	DeleteMethod   = "Plugin.Delete"
	InitMethod     = "Plugin.Init"
	PingMethod     = "Plugin.Ping"
	FeaturesMethod = "Plugin.Features"
)

type RPCClient struct {
	client   *rpc.Client
	features []state.Feature
}

func (s *RPCClient) Init(metadata state.Metadata) error {
	var resp interface{}
	if err := s.client.Call(InitMethod, &metadata, &resp); err != nil {
		return err
	}

	// we need to call Features here because it could return an error and the features method doesn't support that
	var features []state.Feature
	if err := s.client.Call(FeaturesMethod, new(interface{}), &features); err != nil {
		return err
	}
	s.features = features

	return nil
}

func (s *RPCClient) Features() []state.Feature {
	return s.features
}

func (s *RPCClient) Delete(req *state.DeleteRequest) error {
	var resp interface{}
	return s.client.Call(DeleteMethod, req, &resp)
}

func (s *RPCClient) Set(req *state.SetRequest) error {
	var resp interface{}
	return s.client.Call(SetMethod, map[string]interface{}{
		"req": req,
	}, &resp)
}

func (s *RPCClient) Ping() error {
	var resp interface{}
	err := s.client.Call(PingMethod, new(interface{}), &resp)
	if err != nil {
		return err
	}
	return nil
}

func (s *RPCClient) BulkDelete(req []state.DeleteRequest) error {
	return nil
}

func (s *RPCClient) BulkGet(req []state.GetRequest) (bool, []state.BulkGetResponse, error) {
	return false, nil, nil
}

func (s *RPCClient) BulkSet(req []state.SetRequest) error {
	return nil
}
