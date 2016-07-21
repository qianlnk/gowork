package gowork

import (
	"sync"
)

type WorkFunction func(request interface{}, response interface{})

type gowork struct {
	routinenum int
	goworker   WorkFunction
	request    chan interface{}
	wg         *sync.WaitGroup
}

func (g *gowork) workerpool(res interface{}) {
	for i := 0; i < g.routinenum; i++ {
		g.wg.Add(1)
		go g.worker(res)

	}
}

func (g *gowork) worker(res interface{}) {
	defer g.wg.Done()
	for req := range g.request {
		g.goworker(req, res)
	}
}

func (g *gowork) addrequest(req interface{}) {
	g.request <- req
}
