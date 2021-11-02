package main

import (
	"sync"
	"sync/atomic"
)

type idService interface {
	// Returns values in ascending order; it should be safe to call
	// getNext() concurrently without any additional synchronization.
	getNext() uint32
}

// type wildWestService struct {
// 	counter uint64
// }

// func (w *wildWestService) getNext() uint64 {
// 	w.counter++
// 	return w.counter
// }

type atomicService struct {
	counter uint32
}

func (a *atomicService) getNext() uint32 {
	return atomic.AddUint32(&a.counter, 1)
}

type mutexService struct {
	mu      sync.Mutex
	counter uint32
}

func (m *mutexService) getNext() uint32 {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.counter++
	return m.counter
}

type goroutineService struct {
	requests chan struct{}
	results  chan uint32
	counter  uint32
}

// revive:disable-next-line:unexported-return
func NewGoRoutineService() *goroutineService {
	return &goroutineService{
		requests: make(chan struct{}),
		results:  make(chan uint32),
	}
}

func (g *goroutineService) Start() {
	go func() {
		for range g.requests {
			g.counter++
			g.results <- g.counter
		}
	}()
}

func (g *goroutineService) Stop() {
	close(g.requests)
}

func (g *goroutineService) getNext() uint32 {
	g.requests <- struct{}{}
	return <-g.results
}
