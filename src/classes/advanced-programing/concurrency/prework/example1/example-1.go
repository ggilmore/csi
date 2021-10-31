package main

import (
	"fmt"
	"time"
)

// Original:
//
// func main() {
// 	for i := 0; i < 10; i++ {
// 		go func() {
// 			fmt.Printf("launched goroutine %d\n", i)
// 		}()
// 	}
// 	// Wait for goroutines to finish
// 	time.Sleep(time.Second)
// }

// You need to pass "i" as a parameter
// the value of i changes during the loop
func main() {
	for i := 0; i < 10; i++ {
		go func(i int) {
			fmt.Printf("launched goroutine %d\n", i)
		}(i)
	}
	// Wait for goroutines to finish
	time.Sleep(time.Second)
}
