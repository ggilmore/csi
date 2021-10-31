package main

import (
	"fmt"
	"sync"
	"time"
)

// Original:
//
// type dbService struct {
// 	lock       *sync.RWMutex
// 	connection string
// }

// func newDbService(connection string) *dbService {
// 	return &dbService{
// 		lock:       &sync.RWMutex{},
// 		connection: connection,
// 	}
// }

// func (d *dbService) logState() {
// 	d.lock.RLock()
// 	defer d.lock.RUnlock()

// 	fmt.Printf("connection %q is healthy\n", d.connection)
// }

// func (d *dbService) takeSnapshot() {
// 	d.lock.RLock()
// 	defer d.lock.RUnlock()

// 	fmt.Printf("Taking snapshot over connection %q\n", d.connection)

// 	// Simulate slow operation
// 	time.Sleep(time.Second)

// 	d.logState()
// }

// func (d *dbService) updateConnection(connection string) {
// 	d.lock.Lock()
// 	defer d.lock.Unlock()

// 	d.connection = connection
// }

// func main() {
// 	d := newDbService("127.0.0.1:3001")

// 	var wg sync.WaitGroup

// 	wg.Add(1)
// 	go func() {
// 		defer wg.Done()

// 		d.takeSnapshot()
// 	}()

// 	// Simulate other DB accesses
// 	time.Sleep(200 * time.Millisecond)

// 	wg.Add(1)
// 	go func() {
// 		defer wg.Done()

// 		d.updateConnection("127.0.0.1:8080")
// 	}()

// 	wg.Wait()
// }

type dbService struct {
	lock       *sync.RWMutex
	connection string
}

func newDBService(connection string) *dbService {
	return &dbService{
		lock:       &sync.RWMutex{},
		connection: connection,
	}
}

func (d *dbService) LogState() {
	d.lock.RLock()
	defer d.lock.RUnlock()

	d.dologRaw()
}

func (d *dbService) dologRaw() {
	fmt.Printf("connection %q is healthy\n", d.connection)
}

func (d *dbService) takeSnapshot() {
	d.lock.RLock()
	defer d.lock.RUnlock()

	fmt.Printf("Taking snapshot over connection %q\n", d.connection)

	// Simulate slow operation
	time.Sleep(time.Second)

	d.dologRaw()
}

func (d *dbService) updateConnection(connection string) {
	d.lock.Lock()
	defer d.lock.Unlock()

	d.connection = connection
}

// The first goroutine was acquiring the same readlock twice (reentrant).
// Using the same technique as problem 6 fixes the issue.
func main() {
	d := newDBService("127.0.0.1:3001")

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		d.takeSnapshot()
	}()

	// Simulate other DB accesses
	time.Sleep(200 * time.Millisecond)

	wg.Add(1)
	go func() {
		defer wg.Done()

		d.updateConnection("127.0.0.1:8080")
	}()

	wg.Wait()
}
