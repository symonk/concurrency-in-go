package atomicity

import (
	"fmt"
	"sync"
)

// It is often the case that non atomic code must be made atomic.
// Go others a few options here.  The trick comes with minimising
// these concepts around `critical sections` of code to guarantee
// correctness with the least performance impact.
func ensuringAtomicity() {

	// A basic data race
	var n int
	go func() { n++ }() // critical section (1 - write)
	if n == 0 {         // critical section (2 - read)
		fmt.Println("n was 0 at some point, but may not be now", n) // critical section (3 - read)
	} else {
		fmt.Println("n was incremented before line #14 was evaluated")
	}

	// So how do we solve this?
	// note: This code is not idiomatic go and is solely
	// for demo purposes:
	// This is our first introduction to the `sync` package
	// and we utilise for now a read-write mutex.
	var m sync.Mutex // The falsy default value of the mutex is sensible/usable
	var v int

	go func() {
		// This fixes critical section (1).
		// The increment is protected by the mutex
		// even tho the increment is 3 operations
		m.Lock()
		defer m.Unlock()
		v++
	}()

	m.Lock()
	if v == 0 {
		fmt.Println("v was 0")
	} else {
		fmt.Println("v was not zero, it was 1")
	}
	m.Unlock()

	// Looking a t this code, it may seem to be fixed, however it is simply not
	// true.  This code while somewhat protected by the mutex it has no guarantee
	// that the goroutine (anonymous func) would fire before the v equality checks.

}
