package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	sayHello := func() {
		defer wg.Done()
		fmt.Println("hello")
	}

	wg.Add(1)
	sayHello()
	// Time.Sleep(time.Second) <-- Time.Sleep might make this appear to work, but don't do that ever.
	wg.Wait()
}
