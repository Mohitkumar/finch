package service

import (
	"github.com/mohitkumar/finch/action"
	api "github.com/mohitkumar/finch/api/v1"
	"github.com/mohitkumar/finch/flow"
	"github.com/mohitkumar/finch/persistence"
)

type WorkflowExecutor struct {
	pFactory persistence.PersistenceFactory
}

func (ex *WorkflowExecutor) StartExecution(wfName string, flow flow.Flow, flowContext *api.FlowContext) error {
	err := ex.ExecuteAction(wfName, flow.RootAction, flow, flowContext)
	if err != nil {
		return err
	}
	return nil
}

func (ex *WorkflowExecutor) ExecuteAction(wfName string, actionId int, flow flow.Flow, flowContext *api.FlowContext) error {
	result, err := flow.Actions[actionId].Execute(wfName, flowContext)
	if err != nil {
		return err
	}
	switch result.ActionType {
	case action.ACTION_TYPE_SYSTEM:
		if result.Data != nil {
			flowContext, err = ex.pFactory.GetFlowDao().UpdateFlowContextData(wfName, flow.Id, actionId, result.Data)
			if err != nil {
				return err
			}
		}
		return ex.ExecuteAction(wfName, result.NextAction, flow, flowContext)
	case action.ACTION_TYPE_USER:
		flowContext.NextAction = int32(result.NextAction)
		ex.pFactory.GetFlowDao().SaveFlowContext(wfName, flow.Id, flowContext)
	}

}
