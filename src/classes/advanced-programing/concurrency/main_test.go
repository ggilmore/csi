package main

import (
	"fmt"
	"testing"

	"golang.org/x/sync/errgroup"
)

func TestService(t *testing.T) {
	workers := 10
	callsPerWorker := 10000

	goroutineService := NewGoRoutineService()
	goroutineService.Start()
	defer goroutineService.Stop()

	for _, test := range []struct {
		name    string
		service idService
	}{
		// {"wild west service", &wildWestService{}},
		{"atomic service", &atomicService{}},
		{"mutex service", &mutexService{}},
		{"goroutine service", goroutineService},
	} {
		t.Run(test.name, func(t *testing.T) {
			validateService(t, test.service, workers, callsPerWorker)
		})
	}
}

func BenchmarkService(b *testing.B) {
	workers := 10
	callsPerWorker := 10000

	goroutineService := NewGoRoutineService()
	goroutineService.Start()
	defer goroutineService.Stop()

	for _, bench := range []struct {
		name       string
		newService func() (service idService, teardown func())
	}{
		{"atomic service", func() (idService, func()) {
			return &atomicService{}, func() {}
		}},
		{"mutex service", func() (idService, func()) {
			return &mutexService{}, func() {}
		}},
		{"goroutine service", func() (idService, func()) {
			goroutineService := NewGoRoutineService()
			goroutineService.Start()

			teardown := func() {
				goroutineService.Stop()
			}

			return goroutineService, teardown
		}},
	} {
		b.Run(bench.name, func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				service, teardown := bench.newService()
				defer teardown()

				validateService(b, service, workers, callsPerWorker)
			}
		})
	}
}

func validateService(t testing.TB, service idService, workers, callsPerWorker int) {
	t.Helper()

	maxID := workers * callsPerWorker
	maxIDChan := make(chan uint64, workers*callsPerWorker)

	var g errgroup.Group
	for i := 0; i < workers; i++ {
		workerID := i
		g.Go(func() error {
			lastID := uint64(0)

			for j := 0; j < callsPerWorker; j++ {
				id := service.getNext()

				if id <= lastID {
					return fmt.Errorf("(worker %d): ids aren't monotonically increasing (lastID: %d, nextID: %d)", workerID, lastID, id)
				}

				lastID = id
				maxIDChan <- id
			}

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		t.Fatalf(err.Error())
	}

	close(maxIDChan)

	maxIDSeen := uint64(0)
	for id := range maxIDChan {
		if maxIDSeen < id {
			maxIDSeen = id
		}
	}

	if maxIDSeen != uint64(maxID) {
		t.Errorf("expected maxID across all workers to be %d, got %d", maxID, maxIDSeen)
	}
}
