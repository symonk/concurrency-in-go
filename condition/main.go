package main

import (
	"fmt"
	"sync"
	"time"
)

// main demonstrates the use of sync.Cond to create a rendevous point
// for multiple goroutines.  It outlines a naive approach that can
// achieve somewhat of the same thing then suggests improvements.
//
// Note: There are other ways to achieve similar things here and
// often the sync.Cond is not really needed, in fact there are
// proposals to remove it from the language.
func main() {
	fmt.Println("Running the naive CPU spinning example")
	firstNaiveImplementation()
	fmt.Println("Running the naive CPU spinning example with sleep")
	secondStillNaiveImplementation()
	fmt.Println("Running the sync.Cond example")
	thirdUtilisingSyncCond()
}

// firstNaiveImplementation demonstrates a naive approach you might
// use when trying to create some sort of throttle/rendevous point.
//
// It is terrible and will cause a CPU spin as the condition is constantly
// evaluated.
func firstNaiveImplementation() {
	defer fmt.Println("condition was finally true")
	var condition bool
	go func() {
		time.Sleep(5 * time.Second)
		condition = true
	}()
	for condition == false {
		continue
	}
}

// secondStillNaiveImplementation demonstrates a naive approach you might
// add a sleep within the core cycle loop to ease up the CPU, but this is
// still inefficient.
func secondStillNaiveImplementation() {
	defer fmt.Println("condition was finally true")
	var condition bool
	go func() {
		time.Sleep(5 * time.Second)
		condition = true
	}()
	for condition == false {
		time.Sleep(1 * time.Second)
	}
}

// thirdUtilisingSyncCond demonstrates the use of sync.Cond to create a
// rendevous point.  The routines are put to sleep by the scheduler and
// are awoken when the condition is met.
//
// Because we used sync.Cond here, the solution is much more efficient
// because it allows other goroutines to utilise the underlying OS thread
// to run, as the runtime scheduler has put the goroutine into the state
// of sleep.
//
// Understand `Wait` carefully, it actually Unlocks the underlying mutex
// of the Cond when entering it, and then relocks it when it is awoken.
// this can seem like an odd side effect at first.
//
// This uses .Signal() to alert a single goroutine, however you can alert
// all by using the .Broadcast() method.
func thirdUtilisingSyncCond() {
	defer fmt.Println("condition was finally true")
	var condition bool
	var m sync.Mutex
	c := sync.NewCond(&m)

	go func() {
		time.Sleep(5 * time.Second)
		condition = true
		c.Signal()
	}()
	c.L.Lock() // Wait does not do this, we must lock first!
	for condition == false {
		c.Wait()
	}
	c.L.Unlock() // Wait locks the lock when leaving, this is also necessary!
}
