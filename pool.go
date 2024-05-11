package funcpool

import (
	"context"
	"sync"
)

type FuncDef interface {
	Call(context.Context)
}

type FuncPool struct {
	Ctx       context.Context
	Workers   int
	Funcs     chan FuncDef
	Results   chan struct{}
	WaitGroup sync.WaitGroup

	StatMutex   sync.Mutex
	AddCount    int
	ResultCount int
}

func NewFuncPool(ctx context.Context, workers, resultBufferSize int) *FuncPool {
	return &FuncPool{
		Ctx:       ctx,
		Workers:   workers,
		Funcs:     make(chan FuncDef),
		Results:   make(chan struct{}, resultBufferSize),
		WaitGroup: sync.WaitGroup{},
	}
}

func (obj *FuncPool) Start() {
	for i := 0; i < obj.Workers; i++ {
		obj.WaitGroup.Add(1)
		go func() {
			for f := range obj.Funcs {
				innerCtx, cancel := context.WithCancel(obj.Ctx)
				f.Call(innerCtx)
				obj.Results <- struct{}{}
				obj.incrementResultCount()
				cancel()
			}
			obj.WaitGroup.Done()
		}()
	}
}

func (obj *FuncPool) Stop() {
	close(obj.Funcs)
	obj.WaitGroup.Wait()
}

func (obj *FuncPool) AddFunc(f FuncDef) {
	obj.incrementAddCount()
	obj.Funcs <- f
}

func (obj *FuncPool) HasResults() bool {
	return len(obj.Results) > 0 || obj.ResultCount < obj.AddCount
}

func (obj *FuncPool) incrementAddCount() {
	obj.StatMutex.Lock()
	obj.AddCount++
	obj.StatMutex.Unlock()
}

func (obj *FuncPool) incrementResultCount() {
	obj.StatMutex.Lock()
	obj.ResultCount++
	obj.StatMutex.Unlock()
}
