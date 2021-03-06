package action

import (
	api "github.com/mohitkumar/finch/api/v1"
	"github.com/mohitkumar/finch/container"
	"github.com/mohitkumar/finch/model"
	"github.com/mohitkumar/finch/util"
	"google.golang.org/protobuf/proto"
)

var _ Action = new(UserAction)

type UserAction struct {
	baseAction
	nextAction int
}

func NewUserAction(id int, Type ActionType, name string, inputParams map[string]any, nextAction int, container *container.DIContiner) *UserAction {
	return &UserAction{
		baseAction: *NewBaseAction(id, Type, name, inputParams, container),
		nextAction: nextAction,
	}
}

func (ua *UserAction) Execute(wfName string, flowContext *model.FlowContext) error {
	task := &api.Task{
		WorkflowName: wfName,
		FlowId:       flowContext.Id,
		Data:         util.ConvertToProto(ua.ResolveInputParams(flowContext)),
		ActionId:     int32(flowContext.CurrentAction),
	}
	d, err := proto.Marshal(task)
	if err != nil {
		return err
	}
	flowContext.NextAction = ua.nextAction
	err = ua.container.GetFlowDao().SaveFlowContext(wfName, flowContext.Id, flowContext)
	if err != nil {
		return err
	}
	err = ua.container.GetQueue().Push(ua.GetName(), d)
	if err != nil {
		return err
	}
	return nil
}
