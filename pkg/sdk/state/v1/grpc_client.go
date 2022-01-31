package state

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/dapr/components-contrib/state"
	"github.com/dapr/components-contrib/state/utils"
	"github.com/dapr/dapr/pkg/proto/common/v1"
	proto "github.com/dapr/dapr/pkg/proto/state/v1"

	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// see https://developers.google.com/protocol-buffers/docs/reference/go-generated
// see https://github.com/src-d/proteus

// GRPCClient provides a grpc client for the state store
type GRPCClient struct {
	client   proto.StoreClient
	features []state.Feature
}

func NewGRPCClient(client proto.StoreClient) state.Store {
	return &GRPCClient{
		client: client,
	}
}

func (c *GRPCClient) Features() []state.Feature {
	return c.features
}

func (c *GRPCClient) Init(req state.Metadata) error {
	metadata := &proto.MetadataRequest{
		Properties: map[string]string{},
	}
	for k, v := range req.Properties {
		metadata.Properties[k] = v
	}

	// we need to call the method here because features could return an error and the features interface doesn't support errors
	featureResponse, err := c.client.Features(context.TODO(), &emptypb.Empty{})
	if err != nil {
		return err
	}

	c.features = []state.Feature{}
	for _, f := range featureResponse.Feature {
		feature := state.Feature(f)
		c.features = append(c.features, feature)
	}

	_, err = c.client.Init(context.TODO(), metadata)
	return err
}

func (c *GRPCClient) getConsistency(value string) common.StateOptions_StateConsistency {
	consistency, ok := common.StateOptions_StateConsistency_value[value]
	if !ok {
		return common.StateOptions_CONSISTENCY_UNSPECIFIED
	}
	return common.StateOptions_StateConsistency(consistency)
}

func (c *GRPCClient) getConcurrency(value string) common.StateOptions_StateConcurrency {
	concurrency, ok := common.StateOptions_StateConcurrency_value[value]
	if !ok {
		return common.StateOptions_CONCURRENCY_UNSPECIFIED
	}
	return common.StateOptions_StateConcurrency(concurrency)
}

func (c *GRPCClient) Get(req *state.GetRequest) (*state.GetResponse, error) {

	consistency, ok := common.StateOptions_StateConsistency_value[req.Key]
	if !ok {
		consistency = int32(common.StateOptions_CONSISTENCY_UNSPECIFIED)
	}
	request := &proto.GetRequest{
		Key:         req.Key,
		Metadata:    req.Metadata,
		Consistency: common.StateOptions_StateConsistency(consistency),
	}

	etag := ""
	emptyResponse := &state.GetResponse{
		ETag:     &etag,
		Metadata: map[string]string{},
		Data:     []byte{},
	}

	response, err := c.client.Get(context.TODO(), request)
	if err != nil {
		return emptyResponse, err
	}
	if response == nil {
		return emptyResponse, fmt.Errorf("response is nil")
	}

	etag = response.GetEtag().Value
	return &state.GetResponse{
		Data:     response.GetData(),
		ETag:     &etag,
		Metadata: response.GetMetadata(),
	}, nil
}

func (c *GRPCClient) Set(req *state.SetRequest) error {
	var bytes []byte
	switch t := req.Value.(type) {
	case string:
		bytes = []byte(t)
	case []byte:
		bytes = t
	default:
		var err error
		if bytes, err = utils.Marshal(t, json.Marshal); err != nil {
			return err
		}
	}
	request := &proto.SetRequest{
		Key:   req.GetKey(),
		Value: bytes,
		Etag: &common.Etag{
			Value: *req.ETag,
		},
		Metadata: req.GetMetadata(),
		Options: &common.StateOptions{
			Concurrency: c.getConcurrency(req.Options.Concurrency),
			Consistency: c.getConsistency(req.Options.Consistency),
		},
	}
	_, err := c.client.Set(context.TODO(), request)
	return err
}

func (c *GRPCClient) Ping() error {
	empty := &emptypb.Empty{}
	_, err := c.client.Ping(context.TODO(), empty)
	return err
}

func (c *GRPCClient) Delete(req *state.DeleteRequest) error {
	request := &proto.DeleteRequest{
		Key: req.GetKey(),
		Etag: &common.Etag{
			Value: *req.ETag,
		},
		Metadata: req.GetMetadata(),
		Options: &common.StateOptions{
			Concurrency: c.getConcurrency(req.Options.Concurrency),
			Consistency: c.getConsistency(req.Options.Consistency),
		},
	}
	_, err := c.client.Delete(context.TODO(), request)
	return err
}

func (c *GRPCClient) BulkDelete(req []state.DeleteRequest) error {
	return nil
}

func (c *GRPCClient) BulkGet(req []state.GetRequest) (bool, []state.BulkGetResponse, error) {
	return false, nil, nil
}

func (c *GRPCClient) BulkSet(req []state.SetRequest) error {
	return nil
}
