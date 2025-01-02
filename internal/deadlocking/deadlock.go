package deadlocking

import (
	"fmt"
	"sync"
	"time"
)

// SimpleDeadLock demonstrates a very basic `deadlock“ scenario.
// The go runtime will actually protect us from some types of
// deadlocking scenarios (this is outlined in runtime protections)
//
// In this case, two goroutines are spawned, each battling over the same
// locks, slightly out of order.  This causes both routines to be blocked
// and the goruntime will detect and panic with a fatal error
//
// fatal error: all goroutines are asleep - deadlock!
//
// the sum function will attempt to lock v1.m and then v2.m
// the additional goroutine (also running sum) will attempt to
// lock v2.m and then v1.m
// This creates a scenario where v1 is locked, trying to lock v2
// v2 is locked, trying to lock v1
// Neither will concede control and both remain blocked indefinitely.
//
// The `Coffman Conditions` are outlined in the documentation for areas
// which are meeting all (4) of it's checks.
//
// Again of note, it is critical to understand that in most cases this will
// deadlock, it is still possible that the first invocation of sum() finishes
// before the second one is even invoked/switched too by the runtime scheduler.
// This demonstrates the complexity of getting concurrent code correct.
func SimpleDeadLock() {
	type v struct {
		m     sync.Mutex
		value int
	}

	// This is the first introduction to the `WaitGroup` type of the
	// `sync` package.  This is essentially a counter to manage the
	// finalization of multiple spawned goroutines.
	var wg sync.WaitGroup

	// create a function that accepts two pointers of type `v`
	// Coffman[1]: This code has mutual exclusion via mutexes.
	sum := func(v1, v2 *v) {
		defer wg.Done()     // defer the wg.Done() call to decrement the waitgroup
		v1.m.Lock()         // Lock the v1 instances mutex.
		defer v1.m.Unlock() // Defer the call to unlock after the function returns.

		time.Sleep(2 * time.Second) // simulate some IO bound work

		// Coffman[2]: This process has a resource and is waiting for another.
		v2.m.Lock()
		defer v2.m.Unlock()

		fmt.Println("sum: ", v1.value+v2.value)

		// Coffman[3]: There is no way for this function (goroutine) when spawned in one
		// to be pre-empted from outside.
	}

	var one, two v
	wg.Add(2) // Set the waitgroup count to 2, we want to wait for 2 routines to finish
	// Coffman [4]: The first invocation of sum() is waiting on the subsequent one.
	go sum(&one, &two)
	go sum(&two, &one)

	wg.Wait() // Wait for the sync.Waitgroup to be 'done', this means both goroutines above have finished.

}
