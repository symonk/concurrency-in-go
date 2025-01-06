package main

import (
	"fmt"
	"sync"
)

func main() {
	onceFunc()

}

// onceFunc utilises sync.OnceFunc
//
// sync.OnceFunc wraps a function, returning another function
// that when invoked is executed but guarantees synchronised
// run once only semantics, that is multiple goroutines attempting
// to run the returned func `f` will be blocked until the first
// routine to call it has finished.
//
// if the invocation of f panics, all subsequent calls by other
// goroutines will also panic.
func onceFunc() {
	var i int

	// NOTE: should f panic, all subsequent asynchronous calls to f will also panic.
	f := sync.OnceFunc(func() {
		fmt.Println("this will only run once")
		i++
	})
	var wg sync.WaitGroup
	wg.Add(3)
	for i := 0; i < 3; i++ {
		go func() {
			defer wg.Done()
			f()
		}()
	}
	wg.Wait()
	fmt.Println("all 3 routines finished, but i only incremented once!")
}
