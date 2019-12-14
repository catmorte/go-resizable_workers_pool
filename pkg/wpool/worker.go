package wpool

type worker struct {
	jobCh  <-chan Job
	deadCh chan *worker
	killCh chan *worker
}

func newWorker(jobCh <-chan Job, deadCh chan *worker) *worker {
	w := &worker{
		jobCh:  jobCh,
		deadCh: deadCh,
		killCh: make(chan *worker),
	}
	w.work()
	return w
}

func (w *worker) work() {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				w.deadCh <- w
			}
		}()
	infinity:
		for {
			select {
			case <-w.killCh:
				break infinity
			case job := <-w.jobCh:
				job()
			}
		}
	}()
}
