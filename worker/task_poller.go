package worker

type workerWithStopChannel struct {
	worker Worker
	stop   chan struct{}
}
type TaskPoller struct {
	Workers []*workerWithStopChannel
	Config  WorkerConfiguration
}

func NewTaskPoller(conf WorkerConfiguration) *TaskPoller {
	return &TaskPoller{
		Config: conf,
	}
}

func (tp *TaskPoller) RegisterWorker(worker Worker) {
	stopc := make(chan struct{})
	tp.Workers = append(tp.Workers, &workerWithStopChannel{worker: worker, stop: stopc})
}

func (tp *TaskPoller) Start() {
	for _, w := range tp.Workers {
		client, err := NewClient(tp.Config.ServerUrl)
		if err != nil {
			panic(err)
		}
		pw := &PollerWorker{
			worker: w.worker,
			stop:   w.stop,
			client: client,
		}
		pw.Start()
	}
}

func (tp *TaskPoller) Stop() {
	for _, w := range tp.Workers {
		w.stop <- struct{}{}
	}
}
