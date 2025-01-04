package main

import (
	"context"
	"fmt"
	"time"
)

// main demonstrates how to utilise the select timeout
// pattern.  Frequently done with a context in order
// to time-box a particular goroutine and have it
// terminate if it is running for too long.
//
// If contexts are new to you, just know a context can be
// used to time something out with time, or explicitly by
// calling the cancellation functions.
// They can also be used to share data across boundaries by
// storing data within them, but that is not something of use here.
//
// We give a goroutine upto 2 seconds to get through all the generated
// values.  Should it be too slow, we terminate it.
func main() {
	c := generator()
	ctx, cancelFn := context.WithDeadline(context.Background(), time.Now().Add(2*time.Second))
	defer cancelFn()
	/*
		Here we will demonstrate how a long running goroutine
		can be caused to terminate.  The done channel here is
		used to solely not exit immediately and give the goroutine
		some time to run.
	*/
	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			select {
			case i, ok := <-c:
				if !ok {
					fmt.Println("we generated all the values")
					return
				}
				fmt.Println("received value: ", i)
			case <-ctx.Done():
				fmt.Println("timed out, error:", ctx.Err())
				return
			}
		}
	}()
	<-done
}

// generator yields the integer values 1 through 5.
func generator() <-chan int {
	out := make(chan int)
	go func() {
		for i := range 5 {
			out <- i
			time.Sleep(500 * time.Millisecond) // Exceed the context deadline
		}
		close(out)
	}()
	return out
}
