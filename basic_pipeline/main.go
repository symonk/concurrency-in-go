package main

import "fmt"

// main demonstates a pipelining example with 4
// seperate stages.
//
// A generation stage
// A Doubling stage
// A multiplication stage
// A bitwise operating stage
//
// because they share channel types, we can easily compose them.
func main() {
	final := stageThree(stageTwo(stageOne(generator())))
	for v := range final {
		fmt.Println(v)
	}
}

// generator yields the values 1->1,000,000
// it returns a channel to consume its values
func generator() <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for i := range 1_000_000 {
			out <- i
		}
	}()
	return out
}

// stageOne doubles the numbers from the input stream.
func stageOne(upstream <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for in := range upstream {
			out <- in * 2
		}
	}()
	return out
}

// stageTwo multiples the number by 10
func stageTwo(upstream <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for in := range upstream {
			out <- in * 10
		}
	}()
	return out
}

// stageThree bit shifts the result
func stageThree(upstream <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for in := range upstream {
			out <- in << 1
		}
	}()
	return out
}
