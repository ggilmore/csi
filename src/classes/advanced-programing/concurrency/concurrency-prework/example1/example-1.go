package main

import (
	"fmt"
	"time"
)

func main() {
	for i := 0; i < 10; i++ {
		go func() {
			fmt.Printf("launched goroutine %d\n", i)
		}()
	}
	// Wait for goroutines to finish
	time.Sleep(time.Second)
}
