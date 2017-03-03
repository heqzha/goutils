package concurrency

import (
	"fmt"
	"sync"
	"time"

	"github.com/heqzha/goutils/container"
)

type WorkQueue struct {
	q         chan WorkRequest
	maxLength int
}

func newWorkQueue(max int) *WorkQueue {
	return &WorkQueue{
		q:         make(chan WorkRequest, max),
		maxLength: max,
	}
}

func (w *WorkQueue) push(work WorkRequest) error {
	if len(w.q) >= w.maxLength {
		return fmt.Errorf("WorkQueue is full, cannot add more works.")
	}
	w.q <- work
	return nil
}

func (w *WorkQueue) isFull() bool {
	return len(w.q) >= w.maxLength
}

func (w *WorkQueue) isEmpty() bool {
	return len(w.q) == 0
}

type WorkRequest struct {
	f      func(interface{}) interface{}
	params interface{}
	delay  time.Duration
	output chan interface{}
}

const (
	WorkerStateIdle = 0
	WorkerStateBusy = 1
	WorkerStateQuit = 2
)

type Worker struct {
	ID     int
	Work   *WorkQueue
	Status int
	Quit   chan bool
	mutex  *sync.RWMutex
}

func newWorker(id int) Worker {
	worker := Worker{
		ID:     id,
		Work:   newWorkQueue(1),
		Status: WorkerStateIdle,
		Quit:   make(chan bool),
		mutex:  &sync.RWMutex{},
	}
	return worker
}

func (w *Worker) Start() {
	go func() {
		for {
			select {
			case work := <-w.Work.q:
				w.Busy()
				//Work
				if work.output != nil {
					work.output <- work.f(work.params)
				} else {
					work.f(work.params)
				}
				w.Idle()
				time.Sleep(work.delay)
			case <-w.Quit:
				w.Unavailable()
				return
			}
		}
	}()
}

func (w *Worker) Stop() {
	go func() {
		w.Quit <- true
	}()
}

func (w *Worker) Unavailable() {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	w.Status = WorkerStateQuit
}

func (w *Worker) Busy() {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	w.Status = WorkerStateBusy
}

func (w *Worker) Idle() {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	w.Status = WorkerStateIdle
}

func (w *Worker) IsAvailable() bool {
	w.mutex.RLock()
	defer w.mutex.RUnlock()
	return w.Status == WorkerStateIdle
}

type WorkersPool struct {
	container.Queue
	workQ *WorkQueue
	Quit  chan bool
	mutex *sync.Mutex
}

func (wp *WorkersPool) Start(nWorkers int, maxBuffer int) {
	wp.Clear()
	wp.Quit = make(chan bool)
	wp.workQ = newWorkQueue(maxBuffer)
	wp.mutex = &sync.Mutex{}

	for i := 0; i < nWorkers; i++ {
		worker := newWorker(i)
		worker.Start()
		wp.Push(&worker)
	}

	go func() {
		for {
			select {
			case work := <-wp.workQ.q:
				go func() {
					wp.mutex.Lock()
					worker := wp.Pop().(*Worker)
					if worker.IsAvailable() {
						worker.Work.q <- work
					} else {
						time.Sleep(time.Millisecond * 50)
						wp.workQ.push(work)
					}
					wp.Push(worker)
					wp.mutex.Unlock()
				}()
			case <-wp.Quit:
				return
			}
		}
	}()
}

func (wp *WorkersPool) Collect(f func(interface{}) interface{}, params interface{}, delay time.Duration) error {
	if wp.workQ == nil {
		return fmt.Errorf("WorkQueue is nil.")
	}
	work := WorkRequest{
		f:      f,
		params: params,
		delay:  delay,
		output: nil,
	}
	return wp.workQ.push(work)
}

func (wp *WorkersPool) CollectWithOutput(f func(interface{}) interface{}, params interface{}, delay time.Duration, output chan interface{}) error {
	if wp.workQ == nil {
		return fmt.Errorf("WorkQueue is nil.")
	}
	work := WorkRequest{
		f:      f,
		params: params,
		delay:  delay,
		output: output,
	}
	return wp.workQ.push(work)
}

func (wp *WorkersPool) Stop() {
	for wp.Len() > 0 {
		w := wp.Pop().(*Worker)
		w.Stop()
	}
	wp.Quit <- true
}

func (wp *WorkersPool) IsFull() bool {
	return wp.workQ.isFull()
}

func (wp *WorkersPool) IsEmpty() bool {
	return wp.workQ.isEmpty()
}
