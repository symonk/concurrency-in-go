package main

import (
	"sync"
	"testing"
)

/*
A simple benchmark of swapping context between goroutines
Run this with:

(from a teminal in the repository root)
pushd goroutine_context_switching && go test -bench=. -cpu=1 && popd
*/

func BenchmarkSwap(b *testing.B) {
	var wg sync.WaitGroup
	begin := make(chan struct{})
	c := make(chan struct{})

	var token struct{}

	sender := func() {
		defer wg.Done()
		<-begin
		for i := 0; i < b.N; i++ {
			c <- token
		}
	}

	receiver := func() {
		defer wg.Done()
		<-begin
		for i := 0; i < b.N; i++ {
			<-c
		}
	}

	wg.Add(2)
	go sender()
	go receiver()
	b.StartTimer()
	close(begin)
	wg.Wait()
}
