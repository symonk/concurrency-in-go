package main

import "fmt"

// main demonstrates the quit signal pattern.
// terminating a go routine by a send on a
// channel.
//
// The goroutine (similarly) to the context example
// waits for a receive on the quit channel to know
// when to terminate.
func main() {
	i := 100
	q := make(chan struct{})
	g := generator(i, q)

	for i := 0; i < 10; i++ {
		v := <-g
		fmt.Println(v)
	}

	// cause a termination event to occur
	q <- struct{}{}
}

// generator accepts an integer and a channel to quit on.
// it returns a channel that indefinitely yields messages
// containing i++ where i+1 is calculated for each read
// on the returned channel.
func generator(i int, quit chan struct{}) <-chan string {
	out := make(chan string)
	go func() {
		for {
			select {
			case out <- fmt.Sprintf("Message %d", i):
				i++
			case <-quit:
				fmt.Println("terminating")
				return
			}
		}
	}()
	return out
}
