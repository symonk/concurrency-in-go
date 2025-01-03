package main

import (
	"fmt"
	"sync"
	"time"
)

// main demonstrates the fan in pattern.
// consolidating data from multiple goroutines.
func main() {
	// invoke a long running io function, three times.
	a, b, c := someIO(), someIO(), someIO()
	fanned := fanIn(a, b, c)
	for element := range fanned {
		fmt.Println(element)
	}

}

// fanIn demonstartes a simple way to fan in multiple channels
// into one.
func fanIn(ch ...<-chan status) <-chan status {
	out := make(chan status)
	var wg sync.WaitGroup
	wg.Add(len(ch))
	for _, channel := range ch {
		go func(c <-chan status) {
			defer wg.Done()
			for {
				select {
				case v, ok := <-c:
					if !ok {
						return
					}
					out <- v
				}
			}
		}(channel)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

// status encapsulates some response from a server
type status struct {
	code    int
	message string
}

// someIO simulates a long running function.
func someIO() <-chan status {
	c := make(chan status, 20)
	go func() {
		defer close(c)
		for i := range 20 {
			time.Sleep(200 * time.Millisecond)
			c <- status{code: i, message: fmt.Sprintf("message %d", i)}
		}
	}()
	return c
}
