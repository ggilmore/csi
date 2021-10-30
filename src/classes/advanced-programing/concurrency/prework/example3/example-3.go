package main

import (
	"fmt"
	"net/http"
	"sync"
)

// Original
//
// func main() {
// 	var urls = []string{
// 		"https://bradfieldcs.com/courses/architecture/",
// 		"https://bradfieldcs.com/courses/networking/",
// 		"https://bradfieldcs.com/courses/databases/",
// 	}
// 	var wg sync.WaitGroup
// 	for i := range urls {
// 		go func(i int) {
// 			wg.Add(1)
// 			// Decrement the counter when the goroutine completes.
// 			defer wg.Done()

// 			_, err := http.Get(urls[i])
// 			if err != nil {
// 				panic(err)
// 			}

// 			fmt.Println("Successfully fetched", urls[i])
// 		}(i)
// 	}

// 	// Wait for all url fetches
// 	wg.Wait()
// 	fmt.Println("all url fetches done!")
// }

// You need to do the "wg.Add" outside of the goroutine - otherwise
// there is an ordering where the "wg.Wait" runs before
// all of the "wg.Add" have even executed.
//
// You also need to close the response body
func main() {
	var urls = []string{
		"https://bradfieldcs.com/courses/architecture/",
		"https://bradfieldcs.com/courses/networking/",
		"https://bradfieldcs.com/courses/databases/",
	}

	var wg sync.WaitGroup
	for i := range urls {
		wg.Add(1)
		go func(i int) {
			// Decrement the counter when the goroutine completes.
			defer wg.Done()

			r, err := http.Get(urls[i])
			if err != nil {
				panic(err)
			}

			defer r.Body.Close()

			fmt.Println("Successfully fetched", urls[i])
		}(i)
	}

	// Wait for all url fetches
	wg.Wait()
	fmt.Println("all url fetches done!")
}
