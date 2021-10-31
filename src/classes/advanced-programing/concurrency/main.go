package main

import (
	"sync"
	"sync/atomic"
)

type idService interface {
	// Returns values in ascending order; it should be safe to call
	// getNext() concurrently without any additional synchronization.
	getNext() uint64
}

// type wildWestService struct {
// 	counter uint64
// }

// func (w *wildWestService) getNext() uint64 {
// 	w.counter++
// 	return w.counter
// }

type atomicService struct {
	counter uint64
}

func (a *atomicService) getNext() uint64 {
	return atomic.AddUint64(&a.counter, 1)
}

type mutexService struct {
	mu      sync.Mutex
	counter uint64
}

func (m *mutexService) getNext() uint64 {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.counter++
	return m.counter
}

type goroutineService struct {
	requests chan struct{}
	results  chan uint64
	counter  uint64
}

// revive:disable-next-line:unexported-return
func NewGoRoutineService() *goroutineService {
	return &goroutineService{
		requests: make(chan struct{}),
		results:  make(chan uint64),
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

func (g *goroutineService) getNext() uint64 {
	g.requests <- struct{}{}
	return <-g.results
}
