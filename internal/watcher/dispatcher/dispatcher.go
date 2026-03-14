package dispatcher

import (
	"context"
	"sync"
	"sync/atomic"
)

type Dispatcher struct {
	targetConcurrency int32
	activeWorkers     int32

	minLimit int32
	maxLimit int32

	mu   sync.Mutex
	cond *sync.Cond
}

func New(min, max int32) *Dispatcher {
	d := &Dispatcher{
		targetConcurrency: min,
		minLimit:          min,
		maxLimit:          max,
	}
	d.cond = sync.NewCond(&d.mu)
	return d
}

func (d *Dispatcher) Execute(ctx context.Context, task func()) {
	d.mu.Lock()
	for atomic.LoadInt32(&d.activeWorkers) >= atomic.LoadInt32(&d.targetConcurrency) {
		select {
		case <-ctx.Done():
			d.mu.Unlock()
			return
		default:
			d.cond.Wait()
		}
	}
	d.mu.Unlock()

	atomic.AddInt32(&d.activeWorkers, 1)

	go func() {
		defer func() {
			atomic.AddInt32(&d.activeWorkers, -1)
			d.cond.Broadcast()
		}()

		task()
	}()
}

func (d *Dispatcher) SetLimit(newLimit int32) {
	if newLimit > d.maxLimit {
		newLimit = d.maxLimit
	}
	if newLimit < d.minLimit {
		newLimit = d.minLimit
	}

	atomic.StoreInt32(&d.targetConcurrency, newLimit)

	d.cond.Broadcast()
}

func (d *Dispatcher) Wait() {
	d.mu.Lock()
	defer d.mu.Unlock()

	for atomic.LoadInt32(&d.activeWorkers) > 0 {
		d.cond.Wait()
	}
}
