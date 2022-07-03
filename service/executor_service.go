package service

import (
	"fmt"

	"github.com/google/uuid"
	api "github.com/mohitkumar/finch/api/v1"
	"github.com/mohitkumar/finch/flow"
	"github.com/mohitkumar/finch/logger"
	"github.com/mohitkumar/finch/persistence/factory"
	"github.com/mohitkumar/finch/util"
	"go.uber.org/zap"
)

type WorkflowExecutionService struct {
	pFactory factory.PersistenceFactory
	executor *WorkflowExecutor
}

func (s *WorkflowExecutionService) StartFlow(name string, data map[string]any) error {
	wf, err := s.pFactory.GetWorkflowDao().Get(name)
	if err != nil {
		logger.Error("workflow not found", zap.String("name", name))
		return fmt.Errorf("workflow = %s not found", name)
	}
	flow := flow.Convert(wf, uuid.New().String(), s.pFactory)
	flowCtx, err := s.pFactory.GetFlowDao().CreateAndSaveFlowContext(name, flow.Id, flow.RootAction, data)
	if err != nil {
		return err
	}
	return s.executor.StartExecution(name, flow, flowCtx)
}

func (s *WorkflowExecutionService) HandleTaskResult(taskResult *api.TaskResult) error {
	wfName := taskResult.WorkflowName
	wfId := taskResult.FlowId
	data := util.ConvertFromProto(taskResult.Data)
	switch taskResult.Status {
	case api.TaskResult_SUCCESS:
		wf, err := s.pFactory.GetWorkflowDao().Get(wfName)
		if err != nil {
			logger.Error("workflow not found", zap.String("name", wfName))
			return fmt.Errorf("workflow = %s not found", wfName)
		}
		flow := flow.Convert(wf, wfId, s.pFactory)
		flowCtx, err := s.pFactory.GetFlowDao().UpdateFlowContextData(wfName, wfId, int(taskResult.ActionId), data)
		if err != nil {
			return err
		}
		s.executor.ExecuteAction(wfName, int(flowCtx.NextAction), flow, flowCtx)
	case api.TaskResult_FAIL:
		//retry logic
	}
	return nil
}
