package executor

import (
	"encoding/json"
	"sync"

	api_v1 "github.com/mohitkumar/finch/api/v1"
	"github.com/mohitkumar/finch/container"
	"github.com/mohitkumar/finch/flow"
	"github.com/mohitkumar/finch/logger"
	"github.com/mohitkumar/finch/model"
	"github.com/mohitkumar/finch/util"
	"go.uber.org/zap"
)

var _ Executor = new(DelayExecutor)

type DelayExecutor struct {
	container *container.DIContiner
	sync.WaitGroup
	stop         chan struct{}
	taskExecutor *TaskExecutor
}

func NewDelayExecutor(container *container.DIContiner) *DelayExecutor {
	return &DelayExecutor{
		container:    container,
		taskExecutor: NewTaskExecutor(container),
		stop:         make(chan struct{}),
	}
}

func (ex *DelayExecutor) Name() string {
	return "delay-executor"
}

func (ex *DelayExecutor) Start() error {
	fn := func() {
		res, err := ex.container.GetDelayQueue().Pop("delay_action")
		if err != nil {
			_, ok := err.(api_v1.PollError)
			if !ok {
				logger.Error("error while polling delay queue", zap.Error(err))
			}
			return
		}
		for _, r := range res {
			var msg *model.FlowContextMessage
			json.Unmarshal([]byte(r), &msg)
			wf, err := ex.container.GetWorkflowDao().Get(msg.WorkflowName)
			if err != nil {
				logger.Error("workflow not found", zap.String("name", msg.WorkflowName))
				continue
			}
			flow := flow.Convert(wf, msg.FlowId, ex.container)
			flowCtx, err := ex.container.GetFlowDao().GetFlowContext(wf.Name, flow.Id)
			if err != nil {
				logger.Error("flow context not found", zap.String("name", msg.WorkflowName), zap.String("flowId", msg.FlowId))
				continue
			}
			err = ex.taskExecutor.ExecuteAction(wf.Name, msg.ActionId, flow, flowCtx)
			if err != nil {
				logger.Error("error in executing workflow", zap.String("wfName", msg.WorkflowName), zap.String("flowId", msg.FlowId))
				continue
			}
		}
	}
	tw := util.NewTickWorker(1, ex.stop, fn, &ex.WaitGroup)
	tw.Start()
	return nil
}

func (ex *DelayExecutor) Stop() error {
	ex.stop <- struct{}{}
	return nil
}
