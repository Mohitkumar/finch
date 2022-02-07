package statemachine

import (
	api "github.com/mohitkumar/finch/api/v1"
)

type (
	StateMachineExecutor interface {
		Execute(request *api.ExecutionRequest) (*api.ExecutionResponse, error)
	}

	stateMachineExecutorImpl struct {
	}
)
