package service

import (
	api "github.com/mohitkumar/finch/api/v1"
	"github.com/mohitkumar/finch/flow"
)

type WorkflowExecutor struct {
}

func (ex *WorkflowExecutor) StartExecution(wfName string, flow flow.Flow, flowContext *api.FlowContext) error {
	result, err := flow.Actions[flow.RootAction].Execute(wfName, flowContext)
	if err != nil {
		return err
	}

}
