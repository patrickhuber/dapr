package state

import (
	"context"
	"fmt"

	"github.com/dapr/components-contrib/state"
	"github.com/dapr/dapr/pkg/proto/common/v1"
	statev1pb "github.com/dapr/dapr/pkg/proto/state/v1"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

type GRPCServer struct {
	// this is the real implementation
	Impl state.Store
}

func (s *GRPCServer) Features(ctx context.Context, req *emptypb.Empty) (*statev1pb.FeaturesResponse, error) {
	features := s.Impl.Features()
	featureList := []string{}
	for _, f := range features {
		featureList = append(featureList, string(f))
	}
	return &statev1pb.FeaturesResponse{
		Feature: featureList,
	}, nil
}

func (s *GRPCServer) Init(ctx context.Context, req *statev1pb.MetadataRequest) (*emptypb.Empty, error) {
	metadata := state.Metadata{
		Properties: req.GetProperties(),
	}
	return &emptypb.Empty{}, s.Impl.Init(metadata)
}

func (s *GRPCServer) Get(ctx context.Context, req *statev1pb.GetRequest) (*statev1pb.GetResponse, error) {
	request := &state.GetRequest{
		Key:      req.GetKey(),
		Metadata: req.GetMetadata(),
		Options: state.GetStateOption{
			Consistency: req.Consistency.String(),
		},
	}
	response, err := s.Impl.Get(request)
	if err != nil {
		return nil, err
	}
	if response == nil {
		return nil, fmt.Errorf("response is nil")
	}
	etag := ""
	if response.ETag != nil {
		etag = *response.ETag
	}

	return &statev1pb.GetResponse{
		Data:     response.Data,
		Etag:     &common.Etag{Value: etag},
		Metadata: response.Metadata,
	}, nil
}

func (s *GRPCServer) Set(ctx context.Context, req *statev1pb.SetRequest) (*emptypb.Empty, error) {
	setRequest := &state.SetRequest{
		Key:   req.Key,
		ETag:  &req.Etag.Value,
		Value: req.Value,
		Options: state.SetStateOption{
			Concurrency: req.Options.GetConcurrency().String(),
			Consistency: req.Options.GetConsistency().String(),
		},
	}
	err := s.Impl.Set(setRequest)
	return &emptypb.Empty{}, err
}

func (s *GRPCServer) Delete(ctx context.Context, req *statev1pb.DeleteRequest) (*emptypb.Empty, error) {
	deleteRequest := &state.DeleteRequest{
		Key:      req.Key,
		ETag:     &req.Etag.Value,
		Metadata: req.Metadata,
		Options: state.DeleteStateOption{
			Concurrency: req.Options.GetConcurrency().String(),
			Consistency: req.Options.GetConsistency().String(),
		},
	}
	return &emptypb.Empty{}, s.Impl.Delete(deleteRequest)
}

func (s *GRPCServer) Ping(ctx context.Context, req *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, s.Impl.Ping()
}

func (s *GRPCServer) BulkDelete(ctx context.Context, req *statev1pb.BulkDeleteRequest) (*emptypb.Empty, error) {
	return nil, nil
}

func (s *GRPCServer) BulkGet(ctx context.Context, req *statev1pb.BulkGetRequest) (*statev1pb.BulkGetResponse, error) {
	return nil, nil
}

func (s *GRPCServer) BulkSet(ctx context.Context, req *statev1pb.BulkSetRequest) (*emptypb.Empty, error) {
	return nil, nil
}
