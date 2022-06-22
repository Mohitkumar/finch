package service

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/mohitkumar/finch/logger"
	"github.com/mohitkumar/finch/persistence"
	"go.uber.org/zap"
)

type WorkflowExecutionService struct {
	workflowDao persistence.WorkflowDao
	flowDao     persistence.FlowDao
	queue       persistence.Queue
}

func (s *WorkflowExecutionService) StartFlow(name string, data map[string]any) error {
	wf, err := s.workflowDao.Get(name)
	if err != nil {
		logger.Error("workflow not found", zap.String("name", name))
		return fmt.Errorf("workflow = %s not found", name)
	}
	flow := wf.Convert(uuid.New().String(), data)
	flowCtx, err := s.flowDao.CreateAndSaveFlowContext(name, flow.RootAction, flow)
	if err != nil {
		return err
	}
	return flow.Actions[flow.RootAction].Execute(flowCtx)
}
