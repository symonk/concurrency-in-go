package main

import (
	"fmt"
	"strings"
	"sync"
)

func main() {
	onceFunc()
	onceValue()
	onceValues()

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

// onceFuncValue utilises `sync.OnceValue`
//
// Similarly to the sync.OnceFunc function, this offers
// guaranteed run-once semantics with a subtle difference of
// returning a single value to the caller.
//
// subsequent asynchronous calls yield the same value and
// panics in the main invocation of f() will panic all
// other calls.
//
// sync.OnceValue is a generic function.
func onceValue() {
	lower := "foo"
	calls := 0
	var wg sync.WaitGroup
	wg.Add(1000)

	f := sync.OnceValue[string](func() string {
		calls++
		return strings.ToUpper(lower)
	})

	for i := 0; i < 1000; i++ {
		go func() {
			defer wg.Done()
			f()
		}()
	}
	wg.Wait()
	fmt.Println(f())   // every call yields the result from the first invocation.
	fmt.Println(calls) // is only 1
}

// onceValues utilises `sync.OnceValues`
//
// sync.OnceValues is a counter part to sync.OnceValue except
// it can return two value(s).
//
// The same semantics followed that of OnceFunc and OnceValue
// in that the function is only executed once regardless, it's
// return value 'cached' and returned on subsequent calls.
//
// panics in the returned func f(), will panic all subsequent
// calls.
func onceValues() {
	result := make(chan bool, 1)
	invoked := 0

	f := sync.OnceValues(func() (int, int) {
		result <- true
		invoked++
		return 100, 200
	})

	var wg sync.WaitGroup
	wg.Add(500)
	for i := 0; i < 500; i++ {
		go func() {
			defer wg.Done()
			one, two := f()
			_, _ = one, two
		}()
	}

	wg.Wait()
	fmt.Println("called times: ", invoked)
	fmt.Println("it was true:", <-result)
}
