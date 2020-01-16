package throttler

import (
	"errors"
	"sync"
)

// Runner of tasks with queue
type Throttler struct {
	tasks     []Runnable
	toRun     int
	isProcess bool
	sync.Mutex
}

func NewThrottler() Throttler {
	return Throttler{tasks: make([]Runnable, 0), toRun: 0, isProcess: false}
}

// Add one Task to the queue
// Return current queue len
func (t *Throttler) AddTask(task Runnable) int {
	t.Lock()
	t.tasks = append(t.tasks, task)
	t.Unlock()
	return len(t.tasks)
}

// Run count of tasks from queue.
// Return toRun that left after finish
// If count is not positive - return error
func (t *Throttler) Run(count int) (int, error) {
	if count < 0 {
		return t.toRun, errors.New("the runs count below 0")
	}
	t.Lock()
	defer t.Unlock()
	t.toRun += count
	if !t.isProcess {
		t.isProcess = true
		defer func() { t.isProcess = false }()
		for t.toRun > 0 {
			if len(t.tasks) <= 0 {
				return t.toRun, nil
			}
			t.toRun--
			t.tasks[0].Run()
			t.tasks = t.tasks[1:]
		}
	}
	return t.toRun, nil
}

// Return current queue len
func (t Throttler) QueueLen() int {
	return len(t.tasks)
}
