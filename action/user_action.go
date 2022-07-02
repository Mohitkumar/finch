package action

import (
	api "github.com/mohitkumar/finch/api/v1"
	"github.com/mohitkumar/finch/persistence"
	"google.golang.org/protobuf/proto"
)

var _ Action = new(UserAction)

type UserAction struct {
	baseAction
	nextAction int
}

func NewUserAction(id int, Type string, name string, inputParams map[string]any, nextAction int, pFactory persistence.PersistenceFactory) *UserAction {
	return &UserAction{
		baseAction: *NewBaseAction(id, Type, name, inputParams, pFactory),
		nextAction: nextAction,
	}
}

func (ua *UserAction) Execute(wfName string, flowContext *api.FlowContext) (*ActionResult, error) {
	task := &api.Task{
		WorkflowName: wfName,
		FlowId:       flowContext.Id,
		Data:         flowContext.Data,
	}
	d, err := proto.Marshal(task)
	if err != nil {
		return nil, err
	}
	ua.pFactory.GetQueue().Push(ua.GetName(), d)
	result := &ActionResult{
		NextAction: ua.nextAction,
	}
	return result, nil
}
