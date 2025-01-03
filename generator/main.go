package main

import "fmt"

func main() {
	g := generator(10, 20)
	for i := range g {
		fmt.Println(i)
	}
}

// generator yields the integers between start (inclusive)
// and end (exclusive) through the channel it returns.
func generator(start, end int) <-chan int {
	/*
		A new channel is created (unbuffered) in this case as
		the receiver is looking one number at a time.
		A goroutine is responsible for sending the numbers in
		and defering the channel close when finished.

		This allows the range loop on line #7 to iterate all
		values and stop iterating when all have been processed.
	*/
	c := make(chan int)
	go func() {
		defer close(c)
		for i := start; i < end; i++ {
			c <- i
		}
	}()
	return c
}
