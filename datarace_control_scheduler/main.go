package main

import (
	"fmt"
	"runtime"
)

// main demonstrates an attempt at controlling the scheduler to
// erradicate a data race, this is an awful idea and should be avoided.
func main() {
	var i int
	go func() { i++ }()
	runtime.Gosched()
	if i == 0 {
		fmt.Println(i)
	}

}
