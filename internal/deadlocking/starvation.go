package deadlocking

import (
	"fmt"
	"sync"
	"time"
)

// Livelocking is a subset of `Starvation`.
// Starvation is a situation where concurrent processes cannot get
// all the resources it needs to perform work.
//
// The main difference between `livelocking` and `starvation` is typically
// in most cases, livelocked scenarios cause the resources (P1...PN) to be
// relatively 'equally' starved.
//
// In a starvation scenario, it can often be the case that only a few of
// the resources are excessively starved (and get very little runtime)
// or they may even get no runtime at all.
//
// We use the term `Greedy` to describe a process that is starving others.
// Often they have too wide critical sections that are very slow to finish
// where they could of released shared locks earlier.  Consider something
// like slow running HTTP requests that are happening inside a mutex lock
// and only release the lock when finished after potentially N seconds.
//
// Starvation outlines a scenario where a greedy process is hogging all the
// scheduled time and not allowing another routine to fire much at all.
//
// The starved fn is more polite and frees up the locks inbetween its smaller
// pieces of work, This allows other goroutines the chance to run and aquire
// the shared lock.
//
// The greedy fn just steals the lock for the entirety of its execution.
//
// Our output is as follows:
//
// Greedy Func Executed 1788731 cycles in the time limit
// Starved Func Executed 504043 cycles in the time limit
//
// NOTE: The order of the lines output may vary as there is no guarantee which
// routine will finish first and print its output.
func Starvation() {
	var wg sync.WaitGroup
	var m sync.Mutex

	const limit = time.Second

	greedyFn := func(wg *sync.WaitGroup) {
		defer wg.Done()
		var c int
		for begin := time.Now(); time.Since(begin) <= limit; {
			m.Lock()
			time.Sleep(3 * time.Nanosecond)
			m.Unlock()
			c++
		}
		fmt.Println(fmt.Sprintf("Greedy Func Executed %d cycles in the time limit", c))
	}

	starvedFn := func(wg *sync.WaitGroup) {
		defer wg.Done()
		var c int
		for begin := time.Now(); time.Since(begin) <= limit; {
			m.Lock()
			time.Sleep(time.Nanosecond)
			m.Unlock()

			m.Lock()
			time.Sleep(time.Nanosecond)
			m.Unlock()

			m.Lock()
			time.Sleep(time.Nanosecond)
			m.Unlock()
			c++
		}
		fmt.Println(fmt.Sprintf("Starved Func Executed %d cycles in the time limit", c))
	}

	// Add two expecations to the waitgroup and fire off two goroutines
	// one greedy and one starved.
	wg.Add(2)
	go greedyFn(&wg)
	go starvedFn(&wg)
	wg.Wait()

}
