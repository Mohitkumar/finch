package persistence

import (
	"time"

	api "github.com/mohitkumar/finch/api/v1"
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
	CreateAndSaveFlowContext(wFname string, flowId string, action int, dataMap map[string]any) (*api.FlowContext, error)
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
	Queue
	PushWithDelay(queueName string, delay time.Duration, message []byte) error
}
