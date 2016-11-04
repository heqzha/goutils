package goutils

import(
	"time"
	"sync"

	"github.com/heqzha/goutils/container"
)

var(
	mutex = &sync.Mutex{}
	workQ = make(chan WorkRequest, 1024)
)

type WorkRequest struct{
	f func(interface{})interface{}
	params interface{}
	delay time.Duration
	output chan interface{}
}

type Worker struct{
	ID int
	Work chan WorkRequest
	Quit chan bool
}

func newWorker(id int) Worker{
	worker := Worker{
		ID:          id,
		Work:        make(chan WorkRequest),
		Quit:    make(chan bool),
	}
	return worker
}

func (w *Worker) Start(){
	go func(){
		for {
			select {
			case work := <-w.Work:
				//Work
				if work.output != nil{
					work.output<-work.f(work.params)
				}else{
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

type WorkerQueue struct{
	container.Queue
	Quit chan bool
}

func (wq *WorkerQueue) WorkStart(nWorkers int){
	wq.Clear()
	wq.Quit = make(chan bool)

	for i := 0; i<nWorkers; i++ {
		worker := newWorker(i)
		worker.Start()
		wq.Push(&worker)
	}

	go func() {
		for {
			select {
			case work := <- workQ:
				go func() {
					mutex.Lock()
					worker := wq.Pop().(*Worker)
					worker.Work <- work
					wq.Push(worker)
					mutex.Unlock()
				}()
			case <- wq.Quit:
				return
			}
		}
	}()
}

func (wq *WorkerQueue) WorkCollect(f func(interface{})interface{}, params interface{}, delay time.Duration){
	work := WorkRequest{
		f:f,
		params:params,
		delay:delay,
		output:nil,
	}
	workQ <- work
}

func (wq *WorkerQueue) WorkCollectWithOutput(f func(interface{})interface{}, params interface{}, delay time.Duration, output chan interface{}){
	work := WorkRequest{
		f:f,
		params:params,
		delay:delay,
		output:output,
	}
	workQ <- work
}

func (wq *WorkerQueue) WorkStop(){
	for wq.Len() > 0{
		w := wq.Pop().(*Worker)
		w.Stop()
	}
	wq.Quit<-true
}
