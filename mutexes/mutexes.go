package main

import (
	"fmt"
	"sync"
)

// main outlines the use cases and differences between
// the sync.Mutex and sync.RWMutex.
func main() {
	readWriteMutex()
	readOnlyMutex()
}

// readWriteMutex demonstrates an example of a RW mutex which
// is used when the routine wanting to acquire the exclusive
// lock has the intention of writing/updating a shared resource.
// By guarding it with a Mutex, we guarantee a synchronisation
// event and other goroutines in the program will have an
// accurate view of the world regarding it's changes.
func readWriteMutex() {
	var count int
	var mu sync.Mutex // A RW mutex

	incrementFunc := func() {
		mu.Lock()
		defer mu.Unlock()
		fmt.Println("incrementing count")
		count++
	}

	decrementFunc := func() {
		mu.Lock()
		defer mu.Unlock()
		fmt.Println("decrementing count")
		count--
	}

	var wg sync.WaitGroup
	for i := 0; i <= 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			incrementFunc()
			decrementFunc()
		}()
	}

	wg.Wait()
	fmt.Println("we guarantee that count is zero: ", count)

}

func readOnlyMutex() {

}
