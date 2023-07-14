package pool

import (
	"blog/pkg/v"
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
		sync.WaitGroup
	}
)

var (
	gPool *Pool
	once  sync.Once
)

func (p *Pool) Init() {
	for i := 0; i < p.maxWorkers; i++ {
		p.WaitGroup.Add(1)
		go func(id int) {
			defer p.WaitGroup.Done()
			p.work(id)
		}(i)
	}
}

func GetPool() *Pool {
	once.Do(func() {
		gPool = &Pool{
			maxWorkers: v.GetViper().GetInt("pool.maxWorkers"),
			jobs:       make(chan *Job),
			results:    make(chan *Result),
		}
	})
	return gPool
}

func (p *Pool) work(id int) {

}
