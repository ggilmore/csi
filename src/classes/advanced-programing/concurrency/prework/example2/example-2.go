package main

import (
	"fmt"
)

const numTasks = 3

// Original:
//
// func main() {
// 	var done chan struct{}
// 	for i := 0; i < numTasks; i++ {
// 		go func() {
// 			fmt.Println("running task...")

// 			// Signal that task is done
// 			done <- struct{}{}
// 		}()
// 	}

// 	// Wait for tasks to complete
// 	for i := 0; i < numTasks; i++ {
// 		<-done
// 	}
// 	fmt.Printf("all %d tasks done!\n", numTasks)
// }

// The channel is unbuffered, so the goroutines
// block until someone reads from that channel.
// We can fix this by buffering "done" with numTask size.
func main() {
	done := make(chan struct{}, numTasks)
	for i := 0; i < numTasks; i++ {
		go func() {
			fmt.Println("running task...")

			// Signal that task is done
			done <- struct{}{}
		}()
	}

	// Wait for tasks to complete
	for i := 0; i < numTasks; i++ {
		<-done
	}
	fmt.Printf("all %d tasks done!\n", numTasks)
}
