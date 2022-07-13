package worker

import (
	"context"
	"sync"
	"time"

	backoff "github.com/cenkalti/backoff/v4"
	api_v1 "github.com/mohitkumar/finch/api/v1"
	"github.com/mohitkumar/finch/logger"
	"github.com/mohitkumar/finch/util"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type pollerWorker struct {
	worker                   Worker
	client                   *client
	stop                     chan struct{}
	maxRetryBeforeResultPush int
	retryIntervalSecond      int
	wg                       *sync.WaitGroup
}

func (pw *pollerWorker) PollAndExecute() error {
	ctx := context.Background()
	req := &api_v1.TaskPollRequest{
		TaskType: pw.worker.GetName(),
	}
	task, err := pw.client.GetApiClient().Poll(ctx, req)
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				return nil
			case codes.Unavailable:
				pw.client.Refresh()
			}
		}
		return err
	}
	result, err := pw.worker.Execute(util.ConvertFromProto(task.Data))
	var taskResult *api_v1.TaskResult
	if err != nil {
		taskResult = &api_v1.TaskResult{
			WorkflowName: task.WorkflowName,
			FlowId:       task.FlowId,
			ActionId:     task.ActionId,
			Status:       api_v1.TaskResult_SUCCESS,
		}
	} else {
		taskResult = &api_v1.TaskResult{
			WorkflowName: task.WorkflowName,
			FlowId:       task.FlowId,
			ActionId:     task.ActionId,
			Data:         util.ConvertToProto(result),
			Status:       api_v1.TaskResult_SUCCESS,
		}
	}
	b := backoff.WithMaxRetries(backoff.NewConstantBackOff(time.Duration(pw.retryIntervalSecond)*time.Second), uint64(pw.maxRetryBeforeResultPush))
	err = backoff.Retry(func() error {
		_, err := pw.client.GetApiClient().Push(ctx, taskResult)
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.Unavailable:
				pw.client.Refresh()
			}
		}
		if err != nil {
			return err
		}
		return nil
	}, b)
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
