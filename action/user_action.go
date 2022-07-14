package action

import (
	api "github.com/mohitkumar/finch/api/v1"
	"github.com/mohitkumar/finch/persistence/factory"
	"github.com/mohitkumar/finch/util"
	"google.golang.org/protobuf/proto"
)

var _ Action = new(UserAction)

type UserAction struct {
	baseAction
	nextAction int
}

func NewUserAction(id int, Type ActionType, name string, inputParams map[string]any, nextAction int, pFactory *factory.PersistenceFactory) *UserAction {
	return &UserAction{
		baseAction: *NewBaseAction(id, Type, name, inputParams, pFactory),
		nextAction: nextAction,
	}
}

func (ua *UserAction) Execute(wfName string, flowContext *api.FlowContext) error {
	task := &api.Task{
		WorkflowName: wfName,
		FlowId:       flowContext.Id,
		Data:         util.ConvertToProto(ua.ResolveInputParams(flowContext)),
		ActionId:     flowContext.CurrentAction,
	}
	d, err := proto.Marshal(task)
	if err != nil {
		return err
	}
	err = ua.pFactory.GetFlowDao().UpdateFlowContextNextAction(wfName, flowContext.Id, flowContext, ua.nextAction)
	if err != nil {
		return err
	}
	err = ua.pFactory.GetQueue().Push(ua.GetName(), d)
	if err != nil {
		return err
	}
	return nil
}
