package main

import "fmt"

// main demonstrates a typical race condition where an expectation
// The following output is possible:
// 1: The goroutine write did not occur before the first read, i was zero and i was printed as 0.
// 2: The goroutine write occurred before checking i == 0, nothing is printed.
// 3: The goroutine write occurred between both reads, i was printed as 0, then also as 1.
//
// Running this may require many iterations to see all of these cases, but no they are there.  Often it takes a
// slight change to some parameters, such as bigger scale, more users, delays, latency etc for these kinds of bugs
// to appear, but they can be very costly and harder to fix/debug.
func main() {
	var i int
	go func() { i++ }() // read & write
	if i == 0 {         // read
		fmt.Println(i) // read
		if i != 0 {    //  read
			fmt.Println(i) // read
		}
	}
}
