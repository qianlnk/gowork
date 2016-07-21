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

func NewWorkManager() *WorkManager {
	return &WorkManager{
		goworks: make(map[string]*gowork),
	}
}

func (w *WorkManager) register(name string, gtnum int, f WorkFunction) error {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	if _, ok := w.goworks[name]; ok {
		return errors.New(fmt.Sprintf("goworker: %s exist.", name))
	}

	work := &gowork{
		routinenum: gtnum,
		goworker:   f,
		request:    make(chan interface{}),
		wg:         new(sync.WaitGroup),
	}

	w.goworks[name] = work
	return nil
}

func (w *WorkManager) unregister(name string) error {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	if _, ok := w.goworks[name]; !ok {
		return errors.New(fmt.Sprintf("goworker: %s not exist.", name))
	}

	delete(w.goworks, name)
	return nil
}

func (w *WorkManager) NewGoroutine(name string, gtnum int, f WorkFunction, res interface{}) error {
	err := w.register(name, gtnum, f)
	if err != nil {
		return err
	}
	w.goworks[name].workerpool(res)
	return nil
}

func (w *WorkManager) AddRequest(name string, req interface{}) error {
	w.mutex.Lock()
	w.mutex.Unlock()

	if _, ok := w.goworks[name]; !ok {
		return errors.New(fmt.Sprintf("goworker: %s not exist.", name))
	}

	w.goworks[name].addrequest(req)
	return nil
}

func (w *WorkManager) Done(name string) error {
	w.mutex.Lock()
	if _, ok := w.goworks[name]; !ok {
		return errors.New(fmt.Sprintf("goworker: %s not exist.", name))
	}

	close(w.goworks[name].request)
	w.goworks[name].wg.Wait()
	w.mutex.Unlock()
	return w.unregister(name)

}