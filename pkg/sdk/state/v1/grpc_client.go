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

	etag := ""
	emptyResponse := &state.GetResponse{
		ETag:     &etag,
		Metadata: map[string]string{},
		Data:     []byte{},
	}

	response, err := c.client.Get(context.TODO(), c.mapGetRequest(req))
	if err != nil {
		return emptyResponse, err
	}
	if response == nil {
		return emptyResponse, fmt.Errorf("response is nil")
	}

	return c.mapGetResponse(response), nil
}

func (c *GRPCClient) mapGetRequest(req *state.GetRequest) *proto.GetRequest {
	consistency, ok := common.StateOptions_StateConsistency_value[req.Key]
	if !ok {
		consistency = int32(common.StateOptions_CONSISTENCY_UNSPECIFIED)
	}
	return &proto.GetRequest{
		Key:         req.Key,
		Metadata:    req.Metadata,
		Consistency: common.StateOptions_StateConsistency(consistency),
	}
}

func (c *GRPCClient) mapGetResponse(resp *proto.GetResponse) *state.GetResponse {
	var etag *string
	if resp.Etag != nil {
		etag = &resp.Etag.Value
	}
	return &state.GetResponse{
		Data:     resp.GetData(),
		ETag:     etag,
		Metadata: resp.GetMetadata(),
	}
}

func (c *GRPCClient) Set(req *state.SetRequest) error {
	protoRequest, err := c.mapSetRequest(req)
	if err != nil {
		return err
	}
	_, err = c.client.Set(context.TODO(), protoRequest)
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
	protoRequests := []*proto.GetRequest{}
	for _, request := range req {
		protoRequest := c.mapGetRequest(&request)
		protoRequests = append(protoRequests, protoRequest)
	}
	bulkGetRequest := &proto.BulkGetRequest{
		Items: protoRequests,
	}
	bulkGetResponse, err := c.client.BulkGet(context.TODO(), bulkGetRequest)
	if err != nil {
		return false, nil, err
	}
	items := []state.BulkGetResponse{}
	for _, resp := range bulkGetResponse.Items {
		bulkGet := state.BulkGetResponse{
			Key:      resp.GetKey(),
			Data:     resp.GetData(),
			ETag:     &resp.GetEtag().Value,
			Metadata: resp.GetMetadata(),
			Error:    resp.Error,
		}
		items = append(items, bulkGet)
	}
	return bulkGetResponse.Got, items, nil
}

func (c *GRPCClient) BulkSet(req []state.SetRequest) error {
	requests := []*proto.SetRequest{}
	for _, r := range req {
		protoRequest, err := c.mapSetRequest(&r)
		if err != nil {
			return err
		}
		requests = append(requests, protoRequest)
	}
	var err error
	_, err = c.client.BulkSet(context.TODO(), &proto.BulkSetRequest{
		Items: requests,
	})
	return err
}

func (c *GRPCClient) mapSetRequest(req *state.SetRequest) (*proto.SetRequest, error) {
	var bytes []byte
	switch t := req.Value.(type) {
	case string:
		bytes = []byte(t)
	case []byte:
		bytes = t
	default:
		if t == nil {
			return nil, fmt.Errorf("set nil value")
		}
		var err error
		if bytes, err = utils.Marshal(t, json.Marshal); err != nil {
			return nil, err
		}
	}
	var etag *common.Etag
	if req.ETag != nil {
		etag = &common.Etag{
			Value: *req.ETag,
		}
	}
	return &proto.SetRequest{
		Key:      req.GetKey(),
		Value:    bytes,
		Etag:     etag,
		Metadata: req.GetMetadata(),
		Options: &common.StateOptions{
			Concurrency: c.getConcurrency(req.Options.Concurrency),
			Consistency: c.getConsistency(req.Options.Consistency),
		},
	}, nil
}
