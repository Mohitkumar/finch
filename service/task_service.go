package service

import (
	"fmt"

	api "github.com/mohitkumar/finch/api/v1"
	"github.com/mohitkumar/finch/executor"
	"github.com/mohitkumar/finch/flow"
	"github.com/mohitkumar/finch/logger"
	"github.com/mohitkumar/finch/persistence/factory"
	"github.com/mohitkumar/finch/util"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

type TaskExecutionService struct {
	pFactory     *factory.PersistenceFactory
	taskExecutor *executor.TaskExecutor
}

func NewTaskExecutionService(pFactory *factory.PersistenceFactory) *TaskExecutionService {
	return &TaskExecutionService{
		pFactory:     pFactory,
		taskExecutor: executor.NewTaskExecutor(pFactory),
	}
}
func (ts *TaskExecutionService) Poll(taskName string) (*api.Task, error) {
	data, err := ts.pFactory.GetQueue().Pop(taskName)
	if err != nil {
		return nil, err
	}
	task := &api.Task{}
	err = proto.Unmarshal(data, task)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (ts *TaskExecutionService) Push(res *api.TaskResult) error {
	return ts.HandleTaskResult(res)
}

func (s *TaskExecutionService) HandleTaskResult(taskResult *api.TaskResult) error {
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
		flowCtx, err := s.pFactory.GetFlowDao().AddActionOutputToFlowContext(wfName, wfId, int(taskResult.ActionId), data)
		if err != nil {
			return err
		}
		s.taskExecutor.ExecuteAction(wfName, int(flowCtx.NextAction), flow, flowCtx)
	case api.TaskResult_FAIL:
		//retry logic
	}
	return nil
}
