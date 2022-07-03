package action

import (
	api "github.com/mohitkumar/finch/api/v1"
	"github.com/mohitkumar/finch/persistence/factory"
	"google.golang.org/protobuf/proto"
)

var _ Action = new(UserAction)

type UserAction struct {
	baseAction
	nextAction int
}

func NewUserAction(id int, Type ActionType, name string, inputParams map[string]any, nextAction int, pFactory factory.PersistenceFactory) *UserAction {
	return &UserAction{
		baseAction: *NewBaseAction(id, Type, name, inputParams, pFactory),
		nextAction: nextAction,
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
	flowContext.NextAction = int32(ua.nextAction)
	ua.pFactory.GetFlowDao().SaveFlowContext(wfName, flowContext.Id, flowContext)
	ua.pFactory.GetQueue().Push(ua.GetName(), d)
	return nil
}
