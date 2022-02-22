package coordinator

import (
	"context"

	api "github.com/mohitkumar/finch/api/v1"
	"google.golang.org/grpc"
)

type (
	ICoordinator interface {
		CreateFlow(flow *api.Flow) (FlowCreateStatus, error)
	}

	GrpcConfig struct {
		Coordinator ICoordinator
	}

	grpcServer struct {
		api.UnimplementedCoordinatorServer
		*GrpcConfig
	}
)

var _ api.CoordinatorServer = (*grpcServer)(nil)

func newGrpcServer(config *GrpcConfig) (*grpcServer, error) {
	srv := &grpcServer{
		GrpcConfig: config,
	}
	return srv, nil
}

func NewServer(config *GrpcConfig) (*grpc.Server, error) {
	gsrv := grpc.NewServer()

	srv, err := newGrpcServer(config)
	if err != nil {
		return nil, err
	}

	api.RegisterCoordinatorServer(gsrv, srv)

	return gsrv, nil
}

func (s *grpcServer) CreateFlow(ctx context.Context, req *api.FlowCreateRequest) (*api.FlowCreateResponse, error) {
	status, err := s.Coordinator.CreateFlow(req.Flow)
	if err != nil {
		return nil, err
	}
	var st api.FlowCreateResponse_Status
	if status == FlowCreateStatusSuccess {
		st = api.FlowCreateResponse_SUCCESS
	} else {
		st = api.FlowCreateResponse_SUCCESS
	}
	return &api.FlowCreateResponse{Status: st}, nil
}
