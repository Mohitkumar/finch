package service

import (
	api "github.com/mohitkumar/finch/api/v1"
	"github.com/mohitkumar/finch/model"
	"github.com/mohitkumar/finch/persistence"
	"google.golang.org/protobuf/proto"
)

type ActionExecutor interface {
	Execute(wfName string, action model.Action, flowContext *api.FlowContext) error
}

var _ ActionExecutor = new(userActionExecutor)

type userActionExecutor struct {
	queue persistence.Queue
}

func (e *userActionExecutor) Execute(wfName string, action model.Action, flowContext *api.FlowContext) error {
	task := &api.Task{
		WorkflowName: wfName,
		FlowId:       flowContext.Id,
		Data:         flowContext.Data,
	}
	d, err := proto.Marshal(task)
	if err != nil {
		return err
	}
	e.queue.Push(action.GetName(), d)
	return nil
}
