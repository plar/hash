package pool

import (
	"log"
	"sync"
	"sync/atomic"
	"time"
)

type Task struct {
	ID      int64
	Handler func()
}

type TaskQueue chan Task

type Worker struct {
	id         uint
	workerPool chan TaskQueue
	taskQueue  TaskQueue
	quitCh     chan struct{}

	debug               bool
	totalProcessedTasks int
	totalTime           time.Duration
}

type Dispatcher struct {
	wg         sync.WaitGroup
	maxWorkers uint
	workerPool chan TaskQueue
	taskQueue  TaskQueue
	quitCh     chan struct{}

	shutdown int32
	debug    bool
}

// NewWorker create a new worker
func NewWorker(id uint, workerPool chan TaskQueue, quitCh chan struct{}) Worker {
	return Worker{
		id:         id,
		workerPool: workerPool,
		taskQueue:  make(TaskQueue),
		quitCh:     quitCh,
		debug:      false,
	}
}

func (w Worker) start() {
	for {
		// Register the worker task queue in the worker pool
		w.workerPool <- w.taskQueue

		select {
		case task := <-w.taskQueue:
			// Dispatcher has added a task to my taskQueue.
			if w.debug {
				log.Printf("pool worker #%d: started hash %v", w.id, task.ID)
			}

			start := time.Now()
			task.Handler()
			end := time.Since(start)

			w.totalTime += end
			w.totalProcessedTasks++

			if w.debug {
				log.Printf("pool worker #%d: finished hash %v", w.id, task.ID)
			}

		case <-w.quitCh:
			log.Printf("pool worker #%d: quit, processed tasks=%v, total exec time=%v", w.id, w.totalProcessedTasks, w.totalTime)
			return
		}
	}
}

// NewDispatcher creates and returns a new Dispatcher.
func NewDispatcher(taskQueue TaskQueue, maxWorkers uint) *Dispatcher {
	return &Dispatcher{
		maxWorkers: maxWorkers,
		workerPool: make(chan TaskQueue, maxWorkers),
		taskQueue:  taskQueue,
		quitCh:     make(chan struct{}),
		debug:      false,
	}
}

// Run executes the main loop.
func (d *Dispatcher) Run() {
	var i uint
	for i = 1; i <= d.maxWorkers; i++ {
		worker := NewWorker(i, d.workerPool, d.quitCh)
		go worker.start()
	}

	go d.dispatch()
}

// Stop the main loop.
func (d *Dispatcher) Stop() {
	if d.debug {
		log.Println("the dispatcher is stopping")
	}
	atomic.StoreInt32(&d.shutdown, 1)
	d.wg.Wait()
	close(d.quitCh)

	if d.debug {
		log.Println("the dispatcher has been stopped")
	}
}

// Dispatch queues a new task into the task queue.
func (d *Dispatcher) Dispatch(task Task) {
	if atomic.LoadInt32(&d.shutdown) != 0 {
		return // pool is stopped
	}

	if d.debug {
		log.Printf("dispatch a new task %#v", task.ID)
	}

	d.wg.Add(1)
	d.taskQueue <- Task{
		ID: task.ID,
		Handler: func() {
			defer d.wg.Done()
			task.Handler()
		},
	}
}

func (d *Dispatcher) dispatch() {
	if d.debug {
		log.Printf("The dispatcher is ready, total workers=%v\n", d.maxWorkers)
	}
	for {
		select {
		case task := <-d.taskQueue:
			go func() {
				taskQueue := <-d.workerPool // get an available task queue from the worker pool
				taskQueue <- task           // schedule the task
			}()

		case <-d.quitCh:
			return
		}
	}
}
