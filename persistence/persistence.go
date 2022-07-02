package persistence

import (
	"time"

	api "github.com/mohitkumar/finch/api/v1"
	"github.com/mohitkumar/finch/model"
	rd "github.com/mohitkumar/finch/persistence/redis"
)

const WF_PREFIX string = "WF_"
const METADATA_CF string = "METADATA_"

type PersistenceImplementation string

const REDIS_PERSISTENCE_IMPL PersistenceImplementation = "redis"
const INMEMORY_PERSISTENCE_IMPL PersistenceImplementation = "inmemmory"

type WorkflowDao interface {
	Save(wf model.Workflow) error

	Delete(name string) error

	Get(name string) (*model.Workflow, error)
}

type FlowDao interface {
	CreateAndSaveFlowContext(wFname string, flowId string, action int, dataMap map[string]any) (*api.FlowContext, error)
	UpdateFlowContextData(wFname string, flowId string, action int, dataMap map[string]any) (*api.FlowContext, error)
	GetFlowContext(wfName string, flowId string) (*api.FlowContext, error)
	SaveFlowContext(wfName string, flowId string, flowCtx *api.FlowContext) error
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
type RedisConfig struct {
	Host      string
	Port      uint16
	Namespace string
}

type InmemConfig struct {
}
type Config struct {
	RedisConfig RedisConfig
	InmemConfig InmemConfig
}
type PersistenceFactory struct {
	initialized bool
	wfDao       WorkflowDao
	flowDao     FlowDao
	queue       Queue
}

func (p *PersistenceFactory) setInitialized() {
	p.initialized = true
}
func (p *PersistenceFactory) Init(config Config, pImpl PersistenceImplementation) {
	defer p.setInitialized()
	switch pImpl {
	case REDIS_PERSISTENCE_IMPL:
		rdConf := &rd.Config{
			Host:      config.RedisConfig.Host,
			Port:      config.RedisConfig.Port,
			Namespace: config.RedisConfig.Namespace,
		}
		p.wfDao = rd.NewRedisWorkflowDao(*rdConf)
		p.flowDao = rd.NewRedisFlowDao(*rdConf)
		p.queue = rd.NewRedisQueue(*rdConf)
	case INMEMORY_PERSISTENCE_IMPL:

	}
}

func (p *PersistenceFactory) GetWorkflowDao() WorkflowDao {
	if !p.initialized {
		panic("ersistence not initalized")
	}
	return p.wfDao
}

func (p *PersistenceFactory) GetFlowDao() FlowDao {
	if !p.initialized {
		panic("ersistence not initalized")
	}
	return p.flowDao
}

func (p *PersistenceFactory) GetQueue() Queue {
	if !p.initialized {
		panic("ersistence not initalized")
	}
	return p.queue
}
