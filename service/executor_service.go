package service

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/mohitkumar/finch/container"
	"github.com/mohitkumar/finch/executor"
	"github.com/mohitkumar/finch/flow"
	"github.com/mohitkumar/finch/logger"
	"go.uber.org/zap"
)

type WorkflowExecutionService struct {
	container    *container.DIContiner
	taskExecutor *executor.TaskExecutor
}

func NewWorkflowExecutionService(container *container.DIContiner) *WorkflowExecutionService {
	return &WorkflowExecutionService{
		container:    container,
		taskExecutor: executor.NewTaskExecutor(container),
	}
}
func (s *WorkflowExecutionService) StartFlow(name string, input map[string]any) error {
	wf, err := s.container.GetWorkflowDao().Get(name)
	if err != nil {
		logger.Error("workflow not found", zap.String("name", name))
		return fmt.Errorf("workflow = %s not found", name)
	}
	flow := flow.Convert(wf, uuid.New().String(), s.container)
	flowCtx, err := s.container.GetFlowDao().CreateAndSaveFlowContext(name, flow.Id, flow.RootAction, input)
	if err != nil {
		return err
	}
	logger.Info("starting workflow", zap.String("workflow", name), zap.Int("rootAction", flow.RootAction))
	return s.taskExecutor.ExecuteAction(name, flow.RootAction, flow, flowCtx)
}
