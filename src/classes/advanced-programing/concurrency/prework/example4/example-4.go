package main

import (
	"fmt"
)

// Original:
// func main() {
// 	done := make(chan struct{}, 1)
// 	go func() {
// 		fmt.Println("performing initialization...")
// 		<-done
// 	}()

// 	done <- struct{}{}
// 	fmt.Println("initialization done, continuing with rest of program")
// }

// You need to "send" a done message inside the init goroutine, not try to
// receive one.
func main() {
	done := make(chan struct{}, 1)
	go func() {
		fmt.Println("performing initialization...")
		done <- struct{}{}
	}()

	<-done
	fmt.Println("initialization done, continuing with rest of program")
}
