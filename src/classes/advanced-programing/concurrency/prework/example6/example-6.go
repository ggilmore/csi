package main

import (
	"fmt"
	"sync"
)

// Original
// type coordinator struct {
// 	lock   sync.RWMutex
// 	leader string
// }

// func newCoordinator(leader string) *coordinator {
// 	return &coordinator{
// 		lock:   sync.RWMutex{},
// 		leader: leader,
// 	}
// }

// func (c *coordinator) logState() {
// 	c.lock.RLock()
// 	defer c.lock.RUnlock()

// 	fmt.Printf("leader = %q\n", c.leader)
// }

// func (c *coordinator) setLeader(leader string, shouldLog bool) {
// 	c.lock.Lock()
// 	defer c.lock.Unlock()

// 	c.leader = leader

// 	if shouldLog {
// 		c.logState()
// 	}
// }

type coordinator struct {
	lock   sync.RWMutex
	leader string
}

func newCoordinator(leader string) *coordinator {
	return &coordinator{
		lock:   sync.RWMutex{},
		leader: leader,
	}
}

func (c *coordinator) LogState() {
	c.lock.RLock()
	defer c.lock.RUnlock()

	c.doLogRaw()
}

func (c *coordinator) doLogRaw() {
	fmt.Printf("leader = %q\n", c.leader)
}

func (c *coordinator) SetLeader(leader string, shouldLog bool) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.leader = leader

	if shouldLog {
		c.doLogRaw()
	}
}

// The old setLeader function was trying to acquire a read lock
// when it already held a write lock. That operation would never
// complete, resulting in deadlock.
//
// I refactored the logging code to have an unexported function
// just write to stdout without trying to hold a lock. It's the
// caller's responsibility to hold the lock if necessary.
//
// The now-exported LogState and SetLeader are now the public facing
// ways trigger a log, and they each make sure to hold the lock.
func main() {
	c := newCoordinator("us-east")
	c.LogState()
	c.SetLeader("us-west", true)
}
