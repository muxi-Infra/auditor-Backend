package pool

const (
	DefaultWorkers = 3
	DefaultTaskNum = 1000
)

type Task struct {
	Topic string
	Data  []byte
}

type Pool struct {
	PendingJobs chan *Task
	workers     int
}

func NewPool(workers, chanSize int) *Pool {
	if workers <= 0 {
		workers = DefaultWorkers
	}

	if chanSize <= 0 {
		chanSize = DefaultTaskNum
	}

	return &Pool{
		PendingJobs: make(chan *Task, chanSize),
		workers:     workers,
	}
}

func (p *Pool) GetWorkerNums() int {
	return p.workers
}

func (p *Pool) Submit(topic string, data []byte) {
	p.PendingJobs <- &Task{topic, data}
}
