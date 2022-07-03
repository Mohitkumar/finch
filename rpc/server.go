package rpc

import "github.com/mohitkumar/finch/service"

type grpcServer struct {
	executorService *service.WorkflowExecutionService
}
