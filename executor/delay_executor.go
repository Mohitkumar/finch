package executor

import (
	"sync"

	api_v1 "github.com/mohitkumar/finch/api/v1"
	"github.com/mohitkumar/finch/container"
	"github.com/mohitkumar/finch/logger"
	"github.com/mohitkumar/finch/util"
	"go.uber.org/zap"
)

var _ Executor = new(DelayExecutor)

type DelayExecutor struct {
	container *container.DIContiner
	sync.WaitGroup
	stop          chan struct{}
	actionExector *ActionExecutor
}

func NewDelayExecutor(container *container.DIContiner, actionExector *ActionExecutor) *DelayExecutor {
	return &DelayExecutor{
		container:     container,
		actionExector: actionExector,
		stop:          make(chan struct{}),
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

			msg, err := ex.container.ActionExecutionRequestEncDec.Decode([]byte(r))
			if err != nil {
				logger.Error("can not decode action execution request")
				continue
			}
			err = ex.actionExector.Execute(*msg)
			if err != nil {
				logger.Error("error in executing workflow", zap.String("wfName", msg.WorkflowName), zap.String("flowId", msg.FlowId))
				continue
			}
		}
	}
	tw := util.NewTickWorker(1, ex.stop, fn, &ex.WaitGroup)
	tw.Start()
	logger.Info("delay executor started")
	return nil
}

func (ex *DelayExecutor) Stop() error {
	ex.stop <- struct{}{}
	return nil
}
