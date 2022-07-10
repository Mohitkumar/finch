package worker

import (
	"context"
	"sync"
	"time"

	api_v1 "github.com/mohitkumar/finch/api/v1"
	"github.com/mohitkumar/finch/logger"
	"github.com/mohitkumar/finch/util"
	"go.uber.org/zap"
)

type pollerWorker struct {
	worker Worker
	client *client
	stop   chan struct{}
	wg     *sync.WaitGroup
}

func (pw *pollerWorker) PollAndExecute() error {
	ctx := context.Background()
	req := &api_v1.TaskPollRequest{
		TaskType: pw.worker.GetName(),
	}
	task, err := pw.client.GetApiClient().Poll(ctx, req)
	if err != nil {
		return err
	}
	result, err := pw.worker.Execute(util.ConvertFromProto(task.Data))
	if err != nil {
		return err
	}

	taskResult := &api_v1.TaskResult{
		WorkflowName: task.WorkflowName,
		FlowId:       task.FlowId,
		ActionId:     task.ActionId,
		Data:         util.ConvertToProto(result),
	}
	_, err = pw.client.GetApiClient().Push(ctx, taskResult)
	if err != nil {
		return err
	}
	return nil
}

func (pw *pollerWorker) Start() {
	ticker := time.NewTicker(time.Duration(pw.worker.GetPollInterval()) * time.Second)
	pw.wg.Add(1)
	go func() {
		defer pw.wg.Done()
		for {
			select {
			case <-ticker.C:
				err := pw.PollAndExecute()
				if err != nil {
					logger.Error("error while polling", zap.Error(err))
				}
			case <-pw.stop:
				ticker.Stop()
				return
			}
		}
	}()
}
