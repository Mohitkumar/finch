package service

import (
	"github.com/mohitkumar/finch/action"
	api "github.com/mohitkumar/finch/api/v1"
	"github.com/mohitkumar/finch/flow"
	"github.com/mohitkumar/finch/persistence/factory"
)

type WorkflowExecutor struct {
	pFactory factory.PersistenceFactory
}

func (ex *WorkflowExecutor) StartExecution(wfName string, flow flow.Flow, flowContext *api.FlowContext) error {
	err := ex.ExecuteAction(wfName, flow.RootAction, flow, flowContext)
	if err != nil {
		return err
	}
	return nil
}

func (ex *WorkflowExecutor) ExecuteAction(wfName string, actionId int, flow flow.Flow, flowContext *api.FlowContext) error {
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
