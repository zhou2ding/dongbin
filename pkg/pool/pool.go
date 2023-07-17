package pool

import (
	"sync"
)

type (
	Job struct {
		id int
	}
	Result struct {
		job *Job
		sum int
	}
	Pool struct {
		maxWorkers int
		jobs       chan *Job
		results    chan *Result
		wg         sync.WaitGroup
	}
)

var (
	gPool *Pool
	once  sync.Once
)

func GetPool(maxWorkers int) *Pool {
	once.Do(func() {
		gPool = &Pool{
			maxWorkers: maxWorkers,
			jobs:       make(chan *Job),
			results:    make(chan *Result),
		}
	})
	return gPool
}

func (p *Pool) Start() {
	for i := 0; i < p.maxWorkers; i++ {
		go p.worker()
	}
}

func NewPool(maxWorkers int) *Pool {
	return &Pool{
		maxWorkers: maxWorkers,
		jobs:       make(chan *Job),
		results:    make(chan *Result),
	}
}

func (p *Pool) Submit(job *Job) {
	p.wg.Add(1)
	p.jobs <- job
}

func (p *Pool) Wait() {
	p.wg.Wait()
	close(p.results)
}

func (p *Pool) worker() {
	for job := range p.jobs {
		sum := processJob(job)
		result := &Result{
			job: job,
			sum: sum,
		}
		p.results <- result
		p.wg.Done()
	}
}

func processJob(job *Job) int {
	return job.id * 2
}
