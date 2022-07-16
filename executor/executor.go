package executor

import (
	"github.com/mohitkumar/finch/action"
	api "github.com/mohitkumar/finch/api/v1"
	"github.com/mohitkumar/finch/container"
	"github.com/mohitkumar/finch/flow"
)

type TaskExecutor struct {
	container *container.DIContiner
}

func NewTaskExecutor(container *container.DIContiner) *TaskExecutor {
	return &TaskExecutor{
		container: container,
	}
}

func (ex *TaskExecutor) ExecuteAction(wfName string, actionId int, flow flow.Flow, flowContext *api.FlowContext) error {
	if _, ok := flow.Actions[int(actionId)]; !ok {
		return ex.container.GetFlowDao().UpdateFlowStatus(wfName, flowContext.Id, flowContext, api.FlowContext_COMPLETED)
	}
	currentAction := flow.Actions[actionId]
	err := currentAction.Execute(wfName, flowContext)
	if err != nil {
		return err
	}
	nextActionId := flowContext.NextAction

	switch currentAction.GetType() {
	case action.ACTION_TYPE_SYSTEM:
		switch currentAction.GetName() {
		case "switch":
			return ex.ExecuteAction(wfName, int(nextActionId), flow, flowContext)
		case "delay":

		}
	case action.ACTION_TYPE_USER:

	}
	return nil
}
