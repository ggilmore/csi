package main

import (
	"sync"
)

// The issue with this one is a lock-discipline issue.
// We have two operations that acquire multiple locks
//
//
// c.Terminate(): First acquires its own write lock THEN
// 				  the state manager's write lock
//				  via RemoveConsumer()
//
// s.PrintState(): The state manager first acquires its own read
// 				   lock THEN the all of the read locks of its
// 				   consumers in some ordering.
//
// There exists an ordering where neither goroutine can make progress.
//
// The c.Terminate() goroutine is holding its own write lock, but it's waiting
// for the statemanager's write lock so that it can proceed with removeConsumer
//
// The s.PrintState() goroutine is holding its own read lock, but it's waiting
// for the consumer's read lock so that it can proceed.
//
// The solution here is to enforce an ordering when you need to acquire multiple
// locks. Put another way, ALWAYS get the statemanager's lock before you acquire a
// consumer lock.
//
// This can be accomplished multiple ways: having c.Terminate call RemoveConsumer first,
// having c.Terminate manually lock the statemanager, and demote removeConsumer to
// an unexported function that doesn't mess with locks directly, etc.
//
// I went with simply calling RemoveConsumer first, but I also think that it's strange
// to have the consumer directly call methods on the state manager.
// Perhaps it's better to not have Terminate() call RemoveConsumer(), and instead just
// require the caller to do the RemoveConsumer() cleanup operation call aftewards. (This
// would also avoid this deadlock issue since the lock would be released in between)
//

func main() {
	s := NewStateManager(10)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		c := s.GetConsumer(0)
		c.Terminate()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		s.PrintState()
	}()

	wg.Wait()
}
