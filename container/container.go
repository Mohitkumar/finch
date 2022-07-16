package container

import (
	"github.com/mohitkumar/finch/config"
	"github.com/mohitkumar/finch/model"
	"github.com/mohitkumar/finch/persistence"
	rd "github.com/mohitkumar/finch/persistence/redis"
	"github.com/mohitkumar/finch/util"
)

type DIContiner struct {
	initialized              bool
	wfDao                    persistence.WorkflowDao
	flowDao                  persistence.FlowDao
	queue                    persistence.Queue
	delayQueue               persistence.DelayQueue
	FlowContextEncDec        util.EncoderDecoder[model.FlowContext]
	FlowContextMessageEncDec util.EncoderDecoder[model.FlowContextMessage]
}

func (p *DIContiner) setInitialized() {
	p.initialized = true
}

func NewDiContainer() *DIContiner {
	return &DIContiner{
		initialized: false,
	}
}

func (d *DIContiner) Init(conf config.Config) {
	defer d.setInitialized()

	switch conf.StorageType {
	case config.STORAGE_TYPE_REDIS:
		rdConf := &rd.Config{
			Addrs:     conf.RedisConfig.Addrs,
			Namespace: conf.RedisConfig.Namespace,
		}
		d.wfDao = rd.NewRedisWorkflowDao(*rdConf)
		d.flowDao = rd.NewRedisFlowDao(*rdConf)

	case config.STORAGE_TYPE_INMEM:

	}
	switch conf.QueueType {
	case config.QUEUE_TYPE_REDIS:
		rdConf := &rd.Config{
			Addrs:     conf.RedisConfig.Addrs,
			Namespace: conf.RedisConfig.Namespace,
		}
		d.queue = rd.NewRedisQueue(*rdConf)
		d.delayQueue = rd.NewRedisDelayQueue(*rdConf)
	}
	switch conf.EncoderDecoderType {
	case config.PROTO_ENCODER_DECODER:
	default:
		d.FlowContextEncDec = util.NewJsonEncoderDecoder[model.FlowContext]()
		d.FlowContextMessageEncDec = util.NewJsonEncoderDecoder[model.FlowContextMessage]()
	}
}

func (d *DIContiner) GetWorkflowDao() persistence.WorkflowDao {
	if !d.initialized {
		panic("ersistence not initalized")
	}
	return d.wfDao
}

func (d *DIContiner) GetFlowDao() persistence.FlowDao {
	if !d.initialized {
		panic("ersistence not initalized")
	}
	return d.flowDao
}

func (d *DIContiner) GetQueue() persistence.Queue {
	if !d.initialized {
		panic("ersistence not initalized")
	}
	return d.queue
}

func (d *DIContiner) GetDelayQueue() persistence.DelayQueue {
	if !d.initialized {
		panic("ersistence not initalized")
	}
	return d.delayQueue
}
