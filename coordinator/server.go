package coordinator

import (
	"context"

	api "github.com/mohitkumar/finch/api/v1"
	"google.golang.org/grpc"
)

type (
	ICoordinator interface {
		CreateFlow(flow *api.Flow) (*api.FlowCreateResponse, error)
		GetFlow(req *api.FlowGetRequest) (*api.FlowGetResponse, error)
	}

	GetServerer interface {
		GetServers() ([]*api.Server, error)
	}

	GrpcConfig struct {
		Coordinator ICoordinator
		GetServerer GetServerer
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
		return status, err
	}
	return status, nil
}

func (s *grpcServer) GetFlow(ctx context.Context, req *api.FlowGetRequest) (*api.FlowGetResponse, error) {
	res, err := s.Coordinator.GetFlow(req)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (s *grpcServer) GetServers(
	ctx context.Context, req *api.GetServersRequest,
) (
	*api.GetServersResponse, error) {
	servers, err := s.GetServerer.GetServers()
	if err != nil {
		return nil, err
	}
	return &api.GetServersResponse{Servers: servers}, nil
}
