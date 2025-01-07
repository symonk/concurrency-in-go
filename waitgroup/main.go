package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// main demonstrates the abilities of a `sync.WaitGroup`.
//
// A Waitgroup has three methods:
// Add(n int) - Increment the counter by n
// Done() - Decrement the counter by 1
// Wait() - Wait until the counter reaches zero
//
// The use of a WaitGroup is traditionally to fan out
// multiple goroutines with some work, have them internally
// defer a Done() call and in the spawning code have it
// wait for .Wait().
func main() {
	var wg sync.WaitGroup
	routines := 10
	wg.Add(routines) // Add 10 to the internal counter, we are waiting for 10 routines
	start := time.Now()
	for i := 0; i < routines; i++ {
		go func(wg *sync.WaitGroup) {
			defer wg.Done() // Mark a single routine as 'Done', decrementing by 1
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(400)))
		}(&wg)
	}
	wg.Wait() // Wait for all goroutines
	fmt.Printf("finished in %0.2f seconds\n", time.Since(start).Seconds())
}
