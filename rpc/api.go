package rpc

import (
	"context"

	api "github.com/mohitkumar/finch/api/v1"
)

var _ api.TaskServiceServer = (*grpcServer)(nil)

func (srv *grpcServer) Poll(ctx context.Context, req *api.TaskPollRequest) (*api.Task, error) {
	return srv.TaskService.Poll(req.TaskType)
}

func (srv *grpcServer) Push(ctx context.Context, req *api.TaskResult) (*api.TaskResultPushResponse, error) {
	err := srv.TaskService.Push(req)
	if err != nil {
		return &api.TaskResultPushResponse{
			Status: false,
		}, err
	}
	return &api.TaskResultPushResponse{Status: true}, nil
}

func (srv *grpcServer) PollStream(req *api.TaskPollRequest, stream api.TaskService_PollStreamServer) error {
	for {
		select {
		case <-stream.Context().Done():
			return nil
		default:
			res, err := srv.TaskService.Poll(req.TaskType)
			if err != nil {
				return err
			}
			if err = stream.Send(res); err != nil {
				return err
			}
		}
	}
}
