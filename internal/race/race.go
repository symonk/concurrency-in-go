package race

import "fmt"

// main demonstrates a typical race condition where an expectation
// of execution order is expected but not guaranteed.
// What can happen here?
// Depending on goroutine scheduling, which is largely non-deterministic:
// 1. Printing the value `0`.  The goroutine never got CPU time and incremented the number.
//
// 2. Printing nothing and exiting.  The goroutine fired before the number == 0 check where number was set to `1`.
//
// 3. Printing the number was 0 or was it, followed by the number `1`.  The goroutine fired but AFTER line 17,
// Since this is not an atomic operation, it's completely possible!
//
// Running this may require many iterations to see all of these cases, but no they are there.  Often it takes a
// slight change to some parameters, such as bigger scale, more users, delays, latency etc for these kinds of bugs
// to appear, but they can be very costly and harder to fix/debug.
func RaceConditionDemo() {
	var number int // create a new int (size dependent on architecture) with the default falsy value (0).
	go func() {
		number++
	}()
	if number == 0 {
		fmt.Println("number was 0, or was it...? ", number)

		// This demonstrates the non atomic nature of this code.
		if number != 0 {
			fmt.Println("This is possible! number was incremented after the == 0 check, but before the check above!")
		}
	}

}
