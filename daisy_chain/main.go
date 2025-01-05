package main

import "fmt"

// out simulates a goroutine that just adds one to the value
// and is used as the sole goroutine for all chains.
//
// a more realistic example would use different goroutines for
// various steps.
func f(left, right chan int) {
	left <- 1 + <-right
}

// main demonstrates a pattern to daisy chain through
// 10,000 individual goroutines, each doing some work
// on the initial start value, where the leftmost (final)
// routine produces our output and all others are transient
// mechanism to manipulate and forward.
//
// This is known as the daisy chain pattern and is easier to
// visualise as the game 'chinese whispers'.
func main() {
	const n = 10               // 10 chains in the flow
	leftmost := make(chan int) // the final channel, we get our output from here
	// assign two channels for now, left and right to leftmost.
	left := leftmost
	right := leftmost
	// iterate the number of chain/steps
	for i := 0; i < n; i++ {
		right = make(chan int)
		go f(left, right)
		left = right
	}

	// attempt to read a value from right
	go func(c chan int) { c <- 1 }(right)
	// print our final value
	fmt.Println(<-leftmost)

}
