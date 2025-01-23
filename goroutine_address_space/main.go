package main

import (
	"fmt"
	"sync"
)

// main demonstrates that goroutines when spawned operate in the same
// address space they were created in.
//
// This means, closures can capture variables from the parent scope
// and use them accurately.
func main() {
	var wg sync.WaitGroup
	word := "hello"
	wg.Add(1)
	go func() {
		defer wg.Done()
		word = "world"
	}()
	wg.Wait()
	fmt.Println(word) // "world"

}
