package main

// main demonstrates what at a glance looks like
// a potentially atomic operation, a simple increment.
//
// Under thorough inspection there is actually multiple operations
// here, and this is truly only atomic in the context of a single
// threaded / goroutine context OR this code is isolated within a
// single goroutine in a parallel context, it is also considered atomic
// as no other routines have access to the data, it is encapsulated within.
func main() {
	var i int
	// Three operations occur here as far as instructions go
	// 1. Retrieve the value in `i`.
	// 2. Increment the value of `i`.
	// 3. Store the new value of `i`.
	// At ANY of these three points, context switching could occur
	// causing a data race in code that is not expecting it.
	// This is a similar example we outlined as part of the race condition code.
	i++
}
