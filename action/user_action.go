package action

import (
	api "github.com/mohitkumar/finch/api/v1"
	"github.com/mohitkumar/finch/persistence"
	"google.golang.org/protobuf/proto"
)

var _ Action = new(UserAction)

type UserAction struct {
	baseAction
	queue persistence.Queue
}

func NewUserAction(id int, Type string, name string, inputParams map[string]any, queue persistence.Queue) *UserAction {
	return &UserAction{
		baseAction: *NewBaseAction(id, Type, name, inputParams),
		queue:      queue,
	}
}

func (ua *UserAction) Execute(wfName string, flowContext *api.FlowContext) error {
	task := &api.Task{
		WorkflowName: wfName,
		FlowId:       flowContext.Id,
		Data:         flowContext.Data,
	}
	d, err := proto.Marshal(task)
	if err != nil {
		return err
	}
	ua.queue.Push(ua.GetName(), d)
	return nil
}
