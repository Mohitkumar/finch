package executor

import (
	"github.com/mohitkumar/finch/action"
	api "github.com/mohitkumar/finch/api/v1"
	"github.com/mohitkumar/finch/flow"
	"github.com/mohitkumar/finch/persistence/factory"
)

type TaskExecutor struct {
	pFactory *factory.PersistenceFactory
}

func NewTaskExecutor(pFactory *factory.PersistenceFactory) *TaskExecutor {
	return &TaskExecutor{
		pFactory: pFactory,
	}
}

func (ex *TaskExecutor) ExecuteAction(wfName string, actionId int, flow flow.Flow, flowContext *api.FlowContext) error {
	if _, ok := flow.Actions[int(actionId)]; !ok {
		flowContext.WorkflowState = api.FlowContext_COMPLETED
		return ex.pFactory.GetFlowDao().SaveFlowContext(wfName, flowContext.Id, flowContext)
	}
	currentAction := flow.Actions[actionId]
	err := currentAction.Execute(wfName, flowContext)
	if err != nil {
		return err
	}
	nextActionId := flowContext.NextAction

	switch currentAction.GetType() {
	case action.ACTION_TYPE_SYSTEM:
		return ex.ExecuteAction(wfName, int(nextActionId), flow, flowContext)
	case action.ACTION_TYPE_USER:

	}
	return nil
}
