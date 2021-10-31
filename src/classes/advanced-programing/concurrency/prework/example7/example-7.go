package main

import (
	"fmt"
	"sync"
)

const (
	numGoroutines = 100
	numIncrements = 100
)

// Original:
// type counter struct {
// 	count int
// }

// func safeIncrement(lock sync.Mutex, c *counter) {
// 	lock.Lock()
// 	defer lock.Unlock()

// 	c.count += 1
// }

// func main() {
// 	var globalLock sync.Mutex
// 	c := &counter{
// 		count: 0,
// 	}

// 	var wg sync.WaitGroup
// 	for i := 0; i < numGoroutines; i++ {
// 		wg.Add(1)
// 		go func() {
// 			defer wg.Done()

// 			for j := 0; j < numIncrements; j++ {
// 				safeIncrement(globalLock, c)
// 			}
// 		}()
// 	}

// 	wg.Wait()
// 	fmt.Println(c.count)
// }

type counter struct {
	count int
}

func safeIncrement(lock sync.Locker, c *counter) {
	lock.Lock()
	defer lock.Unlock()

	c.count++
}

// The old code was passing the globalLock value, not a pointer to it.
// This means that this was equivalent to each goroutine having their
// individual lock (which is effectively useless).
// Chaning this to a pointer reference fixes the issue.
func main() {
	var globalLock sync.Mutex
	c := &counter{
		count: 0,
	}

	var wg sync.WaitGroup
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for j := 0; j < numIncrements; j++ {
				safeIncrement(&globalLock, c)
			}
		}()
	}

	wg.Wait()
	fmt.Println(c.count)
}
