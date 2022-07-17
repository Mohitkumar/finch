package persistence

import (
	"time"

	"github.com/mohitkumar/finch/model"
)

const WF_PREFIX string = "WF_"
const METADATA_CF string = "METADATA_"

type WorkflowDao interface {
	Save(wf model.Workflow) error

	Delete(name string) error

	Get(name string) (*model.Workflow, error)
}

type FlowDao interface {
	SaveFlowContext(wfName string, flowId string, flowCtx *model.FlowContext) error
	CreateAndSaveFlowContext(wFname string, flowId string, action int, dataMap map[string]any) (*model.FlowContext, error)
	AddActionOutputToFlowContext(wFname string, flowId string, action int, dataMap map[string]any) (*model.FlowContext, error)
	GetFlowContext(wfName string, flowId string) (*model.FlowContext, error)
}

type Queue interface {
	Push(queueName string, mesage []byte) error
	Pop(queuName string) ([]byte, error)
}

type PriorityQueue interface {
	Queue
	PushPriority(queueName string, priority int, mesage []byte) error
}

type DelayQueue interface {
	Push(queueName string, mesage []byte) error
	Pop(queueName string) ([]string, error)
	PushWithDelay(queueName string, delay time.Duration, message []byte) error
}
