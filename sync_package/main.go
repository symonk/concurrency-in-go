package main

import (
	"fmt"
	"sync"
)

func main() {
	onceFunc()

}

// onceFunc demonstrates the usage of the sync.Once
// function.
func onceFunc() {
	/*
		sync.OnceFunc returns a function that when invoked will
		execute the initial function.  It guarantees that the func
		provided will be executed only once, regardless of how many
		goroutines asynchronously attempt to call it.
	*/
	var i int
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
