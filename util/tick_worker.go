package util

import (
	"sync"
	"time"
)

type TickWorker struct {
	stop         chan struct{}
	tickInterval int
	wg           *sync.WaitGroup

	fn func()
}

func NewTickWorker(interval int, stop chan struct{}, fn func(), wg *sync.WaitGroup) *TickWorker {
	return &TickWorker{
		stop:         stop,
		tickInterval: interval,
		wg:           wg,
		fn:           fn,
	}
}

func (tw *TickWorker) Start() {
	ticker := time.NewTicker(time.Duration(tw.tickInterval) * time.Second)
	tw.wg.Add(1)
	go func() {
		defer tw.wg.Done()
		for {
			select {
			case <-ticker.C:
				tw.fn()
			case <-tw.stop:
				ticker.Stop()
				return
			}
		}
	}()
}

func (tw *TickWorker) Stop() {
	tw.stop <- struct{}{}
}
