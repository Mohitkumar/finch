package service

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/mohitkumar/finch/executor"
	"github.com/mohitkumar/finch/flow"
	"github.com/mohitkumar/finch/logger"
	"github.com/mohitkumar/finch/persistence/factory"
	"go.uber.org/zap"
)

type WorkflowExecutionService struct {
	pFactory     *factory.PersistenceFactory
	taskExecutor *executor.TaskExecutor
}

func NewWorkflowExecutionService(pFactory *factory.PersistenceFactory) *WorkflowExecutionService {
	return &WorkflowExecutionService{
		pFactory:     pFactory,
		taskExecutor: executor.NewTaskExecutor(pFactory),
	}
}
func (s *WorkflowExecutionService) StartFlow(name string, input map[string]any) error {
	wf, err := s.pFactory.GetWorkflowDao().Get(name)
	if err != nil {
		logger.Error("workflow not found", zap.String("name", name))
		return fmt.Errorf("workflow = %s not found", name)
	}
	flow := flow.Convert(wf, uuid.New().String(), s.pFactory)
	flowCtx, err := s.pFactory.GetFlowDao().CreateAndSaveFlowContext(name, flow.Id, flow.RootAction, input)
	if err != nil {
		return err
	}
	logger.Info("starting workflow", zap.String("workflow", name), zap.Int("rootAction", flow.RootAction))
	return s.taskExecutor.ExecuteAction(name, flow.RootAction, flow, flowCtx)
}
