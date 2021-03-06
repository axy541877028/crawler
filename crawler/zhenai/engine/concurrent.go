package engine

import "fmt"

type ConcurrentEngine struct {
	Scheduler   Scheduler
	WorkerCount int
}
type Scheduler interface {
	Submit(Request)
	ConfigureMasterWorkerChan(chan Request)
}

func (e *ConcurrentEngine) Run(seeds ...Request) {
	in := make(chan Request)
	out := make(chan ParseResult)
	e.Scheduler.ConfigureMasterWorkerChan(in)
	for i := 0; i < e.WorkerCount; i++ {
		createWorker(in, out)
	}
	for _, r := range seeds {
		e.Scheduler.Submit(r)
	}
	itcount:=0
	for {
		result := <-out
		for _, item := range result.Items {
			fmt.Printf("GOT item #%d :%v \n", itcount,item)
			itcount++
		}
		for _, request := range result.Requests {
			//fmt.Printf("GOT request:%v \n", request)
			e.Scheduler.Submit(request)
		}
	}
}
func createWorker(
	in chan Request, out chan ParseResult) {
	go func() {
		for {
			request := <-in
			result, err := worker(request)
			if err != nil {
				continue
			}
			out <- result
		}
	}()
}
