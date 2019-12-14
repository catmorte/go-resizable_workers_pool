package wpool

import (
	"context"
	"fmt"
	"sync"
)

type Job func()

type WorkersPool interface {
	Do(Job)
	GetSize() uint
	Resize(uint) error
}

type workersPool struct {
	sync.RWMutex
	inJobCh  chan Job
	outJobCh chan Job
	deadCh   chan *worker
	size     uint
	init     sync.Once
	ctx      context.Context
	workers  []*worker
	queue    []Job
}

func NewPool(ctx context.Context, size uint) WorkersPool {
	wp := &workersPool{
		ctx:      ctx,
		size:     size,
		workers:  make([]*worker, size),
		inJobCh:  make(chan Job),
		outJobCh: make(chan Job),
		deadCh:   make(chan *worker),
		queue:    make([]Job, 0),
	}
	for i := uint(0); i < size; i++ {
		wp.workers[i] = newWorker(wp.outJobCh, wp.deadCh)
	}
	go func() {
	infinity:
		for {
			select {
			case <-wp.ctx.Done():
				wp.killSomeWorkers(0, wp.size)
				break infinity
			case v := <-wp.inJobCh:
				wp.queue = append(wp.queue, v)
			case wp.getOutputChannel() <- wp.getJobToDo():
				wp.queue = wp.queue[1:]
			case deadW := <-wp.deadCh:
				wp.Lock()
				i := uint(0)
				for ; i < size; i++ {
					if wp.workers[i] == deadW {
						wp.workers[i] = newWorker(wp.outJobCh, wp.deadCh)
						break
					}
				}
				wp.Unlock()
			}
		}
	}()
	return wp
}

func (wp *workersPool) getOutputChannel() chan Job {
	if len(wp.queue) == 0 {
		return nil
	}
	return wp.outJobCh
}

func (wp *workersPool) getJobToDo() Job {
	if len(wp.queue) == 0 {
		return nil
	}
	return wp.queue[0]
}

func (wp *workersPool) Resize(newSize uint) error {
	wp.Lock()
	defer wp.Unlock()
	if newSize <= 0 {
		return fmt.Errorf("pool size should be more than zero")
	} else if newSize == wp.size {
		return nil
	}
	newWorkers := make([]*worker, newSize)
	copy(newWorkers, wp.workers)
	if newSize > wp.size {
		for i := wp.size; i < newSize; i++ {
			newWorkers[i] = newWorker(wp.outJobCh, wp.deadCh)
		}
	} else if newSize < wp.size {
		wp.killSomeWorkers(newSize, wp.size)
	}
	wp.workers = newWorkers
	wp.size = newSize
	return nil
}

func (wp *workersPool) killSomeWorkers(from, to uint) {
	for i := from; i < to; i++ {
		w := wp.workers[i]
		w.killCh <- w
	}
}

func (wp *workersPool) GetSize() uint {
	return wp.size
}

func (wp *workersPool) Do(job Job) {
	wp.inJobCh <- job
}
