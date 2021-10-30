package main

import (
	"fmt"
	"math/rand"
	"time"
)

var responses = []string{
	"200 OK",
	"402 Payment Required",
	"418 I'm a teapot",
}

func randomDelay(maxMillis int) time.Duration {
	return time.Duration(rand.Intn(maxMillis)) * time.Millisecond
}

func query(endpoint string) string {
	// Simulate querying the given endpoint
	delay := randomDelay(100)
	time.Sleep(delay)

	i := rand.Intn(len(responses))
	return responses[i]
}

// Original:
// // Query each of the mirrors in parallel and return the first
// // response (this approach increases the amount of traffic but
// // significantly improves "tail latency")
// func parallelQuery(endpoints []string) string {
// 	results := make(chan string)
// 	for i := range endpoints {
// 		go func(i int) {
// 			results <- query(endpoints[i])
// 		}(i)
// 	}
// 	return <-results
// }

// func main() {
// 	var endpoints = []string{
// 		"https://fakeurl.com/endpoint",
// 		"https://mirror1.com/endpoint",
// 		"https://mirror2.com/endpoint",
// 	}

// 	// Simulate long-running server process that makes continuous queries
// 	for {
// 		fmt.Println(parallelQuery(endpoints))
// 		delay := randomDelay(100)
// 		time.Sleep(delay)
// 	}
// }

// Query each of the mirrors in parallel and return the first
// response (this approach increases the amount of traffic but
// significantly improves "tail latency")
func parallelQuery(endpoints []string) (endpoint, response string) {
	type result struct {
		endpoint string
		response string
	}

	results := make(chan result, len(endpoints))
	for i := range endpoints {
		go func(i int) {
			endpoint := endpoints[i]
			response := query(endpoint)
			results <- result{
				endpoint: endpoint,
				response: response,
			}
			fmt.Println("done sending ")
		}(i)
	}

	r := <-results
	return r.endpoint, r.response
}

// The results channel was unbuffered, which resulted in a goroutine leak.
// Buffering the channel fixed the issue (and lets the goroutines return).
// Another way to solve it would be to have some kind of context.Done or
// other signal that lets them return
func main() {
	var endpoints = []string{
		"https://fakeurl.com/endpoint",
		"https://mirror1.com/endpoint",
		"https://mirror2.com/endpoint",
	}

	// Simulate long-running server process that makes continuous queries
	for {
		endpoint, response := parallelQuery(endpoints)
		fmt.Printf("%q: %s\n", endpoint, response)
		delay := randomDelay(100)
		time.Sleep(delay)
	}
}
