package service

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/mohitkumar/finch/logger"
	"github.com/mohitkumar/finch/persistence"
	"go.uber.org/zap"
)

type WorkflowExecutionService struct {
	workflowDao    persistence.WorkflowDao
	flowDao        persistence.FlowDao
	actionExecutor ActionExecutor
}

func (s *WorkflowExecutionService) StartFlow(name string, data map[string]any) error {
	wf, err := s.workflowDao.Get(name)
	if err != nil {
		logger.Error("workflow not found", zap.String("name", name))
		return fmt.Errorf("workflow = %s not found", name)
	}
	flow := wf.Convert(uuid.New().String())
	flowCtx, err := s.flowDao.CreateAndSaveFlowContext(name, flow.Id, flow.RootAction, data)
	if err != nil {
		return err
	}
	return s.actionExecutor.Execute(name, flow.Actions[flow.RootAction], flowCtx)
}
