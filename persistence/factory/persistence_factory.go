package factory

import (
	"github.com/mohitkumar/finch/persistence"
	rd "github.com/mohitkumar/finch/persistence/redis"
)

type PersistenceImplementation string

const REDIS_PERSISTENCE_IMPL PersistenceImplementation = "redis"
const INMEMORY_PERSISTENCE_IMPL PersistenceImplementation = "inmemmory"

type RedisConfig struct {
	Host      string
	Port      int
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
	wfDao       persistence.WorkflowDao
	flowDao     persistence.FlowDao
	queue       persistence.Queue
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

func (p *PersistenceFactory) GetWorkflowDao() persistence.WorkflowDao {
	if !p.initialized {
		panic("ersistence not initalized")
	}
	return p.wfDao
}

func (p *PersistenceFactory) GetFlowDao() persistence.FlowDao {
	if !p.initialized {
		panic("ersistence not initalized")
	}
	return p.flowDao
}

func (p *PersistenceFactory) GetQueue() persistence.Queue {
	if !p.initialized {
		panic("ersistence not initalized")
	}
	return p.queue
}
