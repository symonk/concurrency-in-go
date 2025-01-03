package main

import (
	"fmt"
	"time"
)

// main demonstrates the simple concept of a goroutine.
// It does not use any synchonisation mechanisms other
// than a very naive time.Sleep() call to cause a ctx switch
// and an arbitrary delay.
//
// Do not ever do this kind of approach.
func main() {

	// saySomething will (at some non deterministic) point in the future be ran in
	// a goroutine which are multiplexed on an arbitrary number of OS threads.
	// The go runtime will not wait for this to complete before exiting.
	go saySomething("foo")
	time.Sleep(time.Second)

	// A race condition exists here; this could exit before printing anything.
	// Goroutines are considered daemon threads and waiting for them to finish
	// is not something the go runtime will do.

}

// saySomething prints a simple phrase to stdout.
func saySomething(phrase string) {
	fmt.Println(phrase)
}
