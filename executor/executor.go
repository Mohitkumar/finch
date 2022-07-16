package executor

import (
	"github.com/mohitkumar/finch/logger"
	"go.uber.org/zap"
)

type Executor interface {
	Start() error
	Stop() error
	Name() string
}

type Executors struct {
	exArray []Executor
}

func (e *Executors) Register(ex Executor) {
	e.exArray = append(e.exArray, ex)
}

func (e *Executors) Start() {
	for _, ex := range e.exArray {
		err := ex.Start()
		if err != nil {
			logger.Error("error starting executor", zap.String("name", ex.Name()), zap.Error(err))
		}
	}
}

func (e *Executors) Stop() {
	for _, ex := range e.exArray {
		err := ex.Stop()
		if err != nil {
			logger.Error("error stoping executor", zap.String("name", ex.Name()), zap.Error(err))
		}
	}
}
