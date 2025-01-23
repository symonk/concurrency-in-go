package main

import (
	"fmt"
	"runtime"
	"sync"
)

// main demonstrates how lightweight a simple goroutine
// actually is.
func main() {
	memoryUsed := func() uint64 {
		runtime.GC() // Force the garabage collector to run upfront.
		var s runtime.MemStats
		runtime.ReadMemStats(&s)
		return s.Sys
	}

	var c <-chan any
	var wg sync.WaitGroup
	noop := func() { wg.Done(); <-c } // block indefinitely so we can keep all the spawned ones alive to check memory.

	const numGoroutines = 1e5
	wg.Add(numGoroutines)
	before := memoryUsed() // Before we spawn, capture the memory used.
	for i := numGoroutines; i > 0; i-- {
		go noop()
	}
	wg.Wait()
	after := memoryUsed() // After we finish spawning them all and they are all blocking, capture memory used.
	fmt.Printf("%.3fkb", float64(after-before)/numGoroutines/1000)
}
