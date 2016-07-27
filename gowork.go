package gowork

import (
	"fmt"
	"sync"
	"time"
)

type WorkFunction func(request interface{}, response interface{})

type gowork struct {
	routinenum int
	goworker   WorkFunction
	request    chan interface{}
	wg         *sync.WaitGroup
	timeout    int
	mu         sync.Mutex
}

func (g *gowork) workerpool(res interface{}) {
	g.mu.Lock()
	defer g.mu.Unlock()

	for i := 0; i < g.routinenum; i++ {
		g.waitgroupAdd(1)
		go g.worker(res)

	}
}

//when the worker timeout let waitgroup done but it still run a goroutine, I can't
//set the res nil due to it may be a IN/OUT param.
func (g *gowork) worker(res interface{}) {
	defer g.waitgroupDone()
	for req := range g.request {
		done := make(chan bool)
		go func(req, res interface{}) {
			g.goworker(req, res)
			done <- true
		}(req, res)

		select {
		case <-done:
			{
				close(done)
			}
		case <-time.After(time.Duration(g.timeout) * time.Second):
			{
				fmt.Println("timeout")
				//res = nil
			}
		}
	}
}

func (g *gowork) addrequest(req interface{}) {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.request <- req
}

func (g *gowork) close() {
	g.mu.Lock()
	defer g.mu.Unlock()

	close(g.request)
}

func (g *gowork) waitgroupAdd(delta int) {
	g.wg.Add(delta)
}

func (g *gowork) waitgroupDone() {
	g.wg.Done()
}

func (g *gowork) waitgroupWait() {
	g.wg.Wait()
}
