package deadlocking

import (
	"fmt"
	"sync"
	"time"
)

// LiveLocking demonstrates a `livelock` scenario.
// Where a deadlock locks up the program, a live lock is
// slightly different and can be harder to detect.  The
// go runtime also offers next to no protection.
// By definition a livelock is where code is running, potentially
// in many goroutines, but no progress within the system is actually
// being made.  This multiple goroutines waiting for an update that
// never arrives that are performance some logic in a loop internally.
func LiveLocking() {

	var wg sync.WaitGroup
	wg.Add(2)

	// clause is solely used to tick and get the code to exit.
	// otherwise this code is livelocked and would never exit
	// but instead churn the switch default indefinitely.
	//
	// NOTE: go runtime cannot detect this, as from a code POV it
	// is technically correct.
	clause := time.NewTicker(time.Second * 10)

	c := make(chan struct{})

	tryFn := func(wg *sync.WaitGroup, c chan struct{}, clause *time.Ticker) {
		defer wg.Done()
		for {
			select {
			case <-clause.C:
				fmt.Println("Exiting. We have been stuck in a live lock long enough.")
				clause.Stop()
				return
			case <-c:
				fmt.Println("This will never happen, nothing is closing the channel")
			default:
				// To avoid completely burning the CPU, wait a second but this code
				// is the only code that will ever fire here.
				time.Sleep(time.Second)
			}
		}
	}

	go tryFn(&wg, c, clause)
	go tryFn(&wg, c, clause)
	wg.Wait()
}
