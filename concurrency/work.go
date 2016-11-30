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

type Worker struct {
	ID   int
	Work chan WorkRequest
	Quit chan bool
}

func newWorker(id int) Worker {
	worker := Worker{
		ID:   id,
		Work: make(chan WorkRequest),
		Quit: make(chan bool),
	}
	return worker
}

func (w *Worker) Start() {
	go func() {
		for {
			select {
			case work := <-w.Work:
				//Work
				if work.output != nil {
					work.output <- work.f(work.params)
				} else {
					work.f(work.params)
				}
				time.Sleep(work.delay)
			case <-w.Quit:
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

type WorkersPool struct {
	container.Queue
	workQ *WorkQueue
	Quit  chan bool
	mutex *sync.Mutex
}

func newWorkQueue(max int) *WorkQueue {
	return &WorkQueue{
		q:         make(chan WorkRequest, max),
		maxLength: max,
	}
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
					worker.Work <- work
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
