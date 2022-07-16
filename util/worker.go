package util

import (
	"sync"

	"github.com/mohitkumar/finch/logger"
	"go.uber.org/zap"
)

type TaskStop struct{}

type Task interface{}

type Worker struct {
	name     string
	capacity int
	stop     chan struct{}
	wg       *sync.WaitGroup
	taskChan chan func() error
}

func (w *Worker) Start() {
	w.wg.Add(1)
	go func() {
		defer w.wg.Done()

		for {
			select {
			case taskFn := <-w.taskChan:
				err := taskFn()
				if err != nil {
					logger.Error("error in executing task in worker", zap.String("worker", w.name), zap.Any("fn", taskFn))
				}
			case <-w.stop:
				logger.Info("stopping worker", zap.String("worker", w.name))
				return
			}
		}
	}()
}

func (w *Worker) Stop() {
	w.stop <- struct{}{}
}

func NewWorker(name string, wg *sync.WaitGroup, capacity int) *Worker {
	ch := make(chan func() error, capacity)
	stop := make(chan struct{})
	return &Worker{
		taskChan: ch,
		name:     name,
		wg:       wg,
		stop:     stop,
	}
}
