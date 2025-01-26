package main

import (
	"fmt"
	"math"
	"os"
	"sync"
	"text/tabwriter"
	"time"
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

// readonlyMutex demonstrates the uplift you can get if you have a
// scenario where reads are more frequent than writes.  An unlimited
// number of readers can access the RW mutex, it is only locked if
// a writer is acquiring the lock.
func readOnlyMutex() {

	producerFunc := func(wg *sync.WaitGroup, mu sync.Locker) {
		defer wg.Done()
		for i := 5; i > 0; i-- {
			mu.Lock()
			defer mu.Unlock()
			time.Sleep(time.Second)
		}
	}

	observerFunc := func(wg *sync.WaitGroup, mu sync.Locker) {
		defer wg.Done()
		mu.Lock()
		defer mu.Unlock()
	}

	test := func(count int, mutex, rwMutex sync.Locker) time.Duration {
		var wg sync.WaitGroup
		start := time.Now()
		go producerFunc(&wg, mutex)
		for i := count; i > 0; i-- {
			go observerFunc(&wg, rwMutex)
		}
		wg.Wait()
		return time.Since(start)
	}

	tw := tabwriter.NewWriter(os.Stdout, 0, 1, 2, ' ', 0)
	defer tw.Flush()

	var rw sync.RWMutex
	fmt.Fprintf(tw, "Readers\tRWMutex\tMutex\t\n")
	for i := 0; i < 20; i++ {
		count := int(math.Pow(2, float64(i)))
		fmt.Fprintf(
			tw,
			"%d\t%v\t%v\n",
			count,
			test(count, &rw, rw.RLocker()),
			test(count, &rw, &rw),
		)
	}
}
