package plugin

import (
	"encoding/json"
	"fmt"

	"github.com/dapr/components-contrib/state"
	"github.com/dapr/components-contrib/state/utils"
)

// MemoryStore is a store used for testing
type MemoryStore struct {
	data map[string][]byte
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		data: map[string][]byte{},
	}
}

func (s *MemoryStore) Init(metadata state.Metadata) error {

	for k := range s.data {
		delete(s.data, k)
	}
	return nil
}

func (s *MemoryStore) Features() []state.Feature {

	return []state.Feature{}
}

func (s *MemoryStore) Delete(req *state.DeleteRequest) error {

	delete(s.data, req.Key)
	return nil
}

func (s *MemoryStore) Get(req *state.GetRequest) (*state.GetResponse, error) {
	emptyResponse := &state.GetResponse{
		Data:     nil,
		ETag:     nil,
		Metadata: nil,
	}
	if req == nil {

		return emptyResponse, nil
	}

	value, ok := s.data[req.Key]
	if !ok {

		return emptyResponse, nil
	}

	metadata := map[string]string{}
	for k, v := range req.Metadata {
		metadata[k] = v
	}

	var etag *string

	return &state.GetResponse{
		Data:     value,
		ETag:     etag,
		Metadata: map[string]string{},
	}, nil
}

func (s *MemoryStore) Set(req *state.SetRequest) error {

	var bytes []byte

	switch t := req.Value.(type) {
	case string:
		bytes = []byte(t)
	case []byte:
		bytes = t
	default:
		if t == nil {
			return fmt.Errorf("set: request body is nil")
		}
		var err error
		if bytes, err = utils.Marshal(t, json.Marshal); err != nil {
			return err
		}
	}

	s.data[req.Key] = bytes

	return nil
}

func (s *MemoryStore) Ping() error {

	return nil
}

func (s *MemoryStore) BulkDelete(req []state.DeleteRequest) error {
	return nil
}

func (s *MemoryStore) BulkGet(req []state.GetRequest) (bool, []state.BulkGetResponse, error) {
	return false, nil, nil
}

func (s *MemoryStore) BulkSet(req []state.SetRequest) error {
	for _, r := range req {
		err := s.Set(&r)
		if err != nil {
			return err
		}
	}
	return nil
}
