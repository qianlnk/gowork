package gowork

import (
	"errors"
	"fmt"
	"sync"
)

type WorkManager struct {
	goworks map[string]*gowork
	mutex   sync.Mutex
}

const (
	worker_not_exist = "goworker: %s not exist."
)

func NewWorkManager() *WorkManager {
	return &WorkManager{
		goworks: make(map[string]*gowork),
	}
}

func (w *WorkManager) register(name string, gtnum int, f WorkFunction) error {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	if _, ok := w.goworks[name]; ok {
		return errors.New(fmt.Sprintf(worker_not_exist, name))
	}

	work := &gowork{
		routinenum: gtnum,
		goworker:   f,
		request:    make(chan interface{}),
		wg:         new(sync.WaitGroup),
		timeout:    30,
	}

	w.goworks[name] = work
	return nil
}

func (w *WorkManager) unregister(name string) error {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	if _, ok := w.goworks[name]; !ok {
		return errors.New(fmt.Sprintf(worker_not_exist, name))
	}

	delete(w.goworks, name)
	return nil
}

/***************************************************
Function:	create a new goroutine
Parameters:	name	[IN]	the worker name
		gtnum	[IN]	the goroutine number
		f	[IN]	worker function
		res	[OUT]	store the worker result
****************************************************/
func (w *WorkManager) NewGoroutine(name string, gtnum int, f WorkFunction, res interface{}) error {
	err := w.register(name, gtnum, f)
	if err != nil {
		return err
	}
	w.goworks[name].workerpool(res)
	return nil
}

/***************************************************
Function:	add a request param to the specified worker
Parameters:	name	[IN]	the worker name
		req	[IN]	the request param
****************************************************/
func (w *WorkManager) AddRequest(name string, req interface{}) error {
	w.mutex.Lock()
	if _, ok := w.goworks[name]; !ok {
		w.mutex.Unlock()
		return errors.New(fmt.Sprintf(worker_not_exist, name))
	}
	w.mutex.Unlock()

	w.goworks[name].addrequest(req)
	return nil
}

/***************************************************
Function:	close the specified worker and unregister it
Parameters:	name	[IN]	the worker name
****************************************************/
func (w *WorkManager) Done(name string) error {
	w.mutex.Lock()
	if _, ok := w.goworks[name]; !ok {
		w.mutex.Unlock()
		return errors.New(fmt.Sprintf(worker_not_exist, name))
	}

	w.goworks[name].close()
	w.goworks[name].waitgroupWait()
	w.mutex.Unlock()
	fmt.Println(name, "done")
	return w.unregister(name)

}

/***************************************************
Function:	set a goroutine timeout, default 30s
Parameters:	name	[IN]	the worker name
		timeout	[IN]	a goroutine timeout
****************************************************/
func (w *WorkManager) SetTimeout(name string, timeout int) error {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	if _, ok := w.goworks[name]; !ok {
		return errors.New(fmt.Sprintf(worker_not_exist, name))
	}
	w.goworks[name].timeout = timeout
	return nil
}
