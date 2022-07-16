package service

import (
	"fmt"

	api "github.com/mohitkumar/finch/api/v1"
	"github.com/mohitkumar/finch/container"
	"github.com/mohitkumar/finch/executor"
	"github.com/mohitkumar/finch/flow"
	"github.com/mohitkumar/finch/logger"
	"github.com/mohitkumar/finch/util"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

type TaskExecutionService struct {
	container    *container.DIContiner
	taskExecutor *executor.TaskExecutor
}

func NewTaskExecutionService(container *container.DIContiner) *TaskExecutionService {
	return &TaskExecutionService{
		container:    container,
		taskExecutor: executor.NewTaskExecutor(container),
	}
}
func (ts *TaskExecutionService) Poll(taskName string) (*api.Task, error) {
	data, err := ts.container.GetQueue().Pop(taskName)
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
		wf, err := s.container.GetWorkflowDao().Get(wfName)
		if err != nil {
			logger.Error("workflow not found", zap.String("name", wfName))
			return fmt.Errorf("workflow = %s not found", wfName)
		}
		flow := flow.Convert(wf, wfId, s.container)
		flowCtx, err := s.container.GetFlowDao().AddActionOutputToFlowContext(wfName, wfId, int(taskResult.ActionId), data)
		if err != nil {
			return err
		}
		s.taskExecutor.ExecuteAction(wfName, int(flowCtx.NextAction), flow, flowCtx)
	case api.TaskResult_FAIL:
		//retry logic
	}
	return nil
}
