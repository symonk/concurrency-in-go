package main

import (
	"fmt"
	"time"
)

// main demonstrates a naive solution to the data race outlined
// in 01_datarace. The idea is simple, if we introduce more time
// the goroutine (write) operation should be performed.
//
// While this can make it seem like the data race is no longer
// present, it is a horrible solution and should be avoided.
func main() {
	var i int
	go func() { i++ }()
	time.Sleep(time.Second)
	if i == 0 {
		fmt.Println(i)
	}
}
