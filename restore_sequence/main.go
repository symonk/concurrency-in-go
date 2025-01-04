package main

import (
	"fmt"
	"sync"
	"time"
)

// main demonstrates the request sequence pattern.
// by restore sequence, we mean that with multiple
// goroutines performing work, we can ensure that they
// each take their turn and are given a fair chance to
// provide their results.
//
// Fan in has no guarantee and you could have two goroutines
// where goroutine A yields all values before any of Goroutine B's
//
// Restore Sequence forces turn based results by utilising a blocking
// chan shared between messages of each goroutine.
// resulting in A, B, A, B and so on and so forth.
func main() {
	fanned := fanInMerge(respond(1), respond(2), respond(3))

	/*
		Even tho we have 3 invocations of fanInMerge, it actually spawns 5
		working routines internally, so we have 15 goroutines total (3x5)
		all being fanned in through the fanned channel.

		Each call to respond(1), response(2), respond(3) will yield 5 responses
		and internally shares a done channel to ensure only any one of its internal
		goroutines can yield a result when expected.

		Here we are attempting to achieve that every 3 reads of the channel (fanned)
		comes from each of the individual channels we originally made.  Order is not
		entirely guaranteed due to the nature of scheduling but there should always be
		a case of A, B, C results in each block of 3 reads.
	*/

	expected := 5 // 3 x 5 responses

	for i := 0; i < expected; i++ {
		// for each of the results, attempt to get a value
		// from each
		first := <-fanned
		second := <-fanned
		third := <-fanned
		fmt.Println(first)
		fmt.Println(second)
		fmt.Println(third)
		// signal to each in the order we received, it's ok to yield again
		first.wait <- struct{}{}
		second.wait <- struct{}{}
		third.wait <- struct{}{}
	}
}

// reply encapsulates some response from a server
// each reply stores an id to easily demonstrate
// the sequence restoration in practice.
type reply struct {
	id       int
	duration time.Duration
	message  string
	wait     chan struct{}
}

// String implements fmt.Stringer and returns a string
// representation of the reply instance.
func (r reply) String() string {
	return fmt.Sprintf("id: %d, duration: %s, message: %s", r.id, r.duration, r.message)
}

// fanInMerge merges the multiple channels of work
func fanInMerge(ch ...<-chan reply) <-chan reply {
	out := make(chan reply)
	var wg sync.WaitGroup
	wg.Add(len(ch))
	for _, channel := range ch {
		go func(c <-chan reply) {
			defer wg.Done()
			for v := range channel {
				out <- v
			}
		}(channel)
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

// respond simulates IO bound calls to a server and returns a number
// of responses on the out channel.
// Each goroutine spawned by this function will push a new instance
// onto the out channel, then wait until it is signalled to yield
// another value.
//
// This ensures that all the goroutines from multiple invocations of
// respond operate in a turn based manner and the order of responses
// across multiple calls to this function are controlled.
func respond(id int) <-chan reply {
	out := make(chan reply)
	waiter := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(5)

	go func() {
		for i := range 5 {
			out <- reply{duration: time.Duration(1000 * i), message: fmt.Sprintf("message %d", i), id: id, wait: waiter}
			time.Sleep(100 * time.Millisecond)
			<-waiter
		}
	}()

	go func() {
		wg.Wait()
		close(out)
		close(waiter)
	}()
	return out
}
