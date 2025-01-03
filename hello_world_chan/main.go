package main

import "fmt"

// main demonstrates one of the core go synchronisation
// primitves, the channel.  There is a lot of complexity
// to channels in go, see the `channels` documentation in
// this repository for a deep dive.
//
// This example fixes the race condition in the hello_world_goroutine
// example.
func main() {
	/*
		We make an (unbuffered) channel, one with no length/cap.
		After the spawned goroutine is finished, we close the channel
		which causes the receive <- on it to finally be unblocked.
	*/
	c := make(chan struct{}) // We use struct{} as it uses zero memory.
	go func() {
		saySomething("foo")
		close(c)
	}()
	<-c // This will block until the channel is closed, causing the program not to exit.
}

// saySomething prints a simple phrase to stdout.
func saySomething(phrase string) {
	fmt.Println(phrase)
}
