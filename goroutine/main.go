package main

import "fmt"

func main() {
	go HelloWorld()

	// Anonymous functions work too!
	go func() {
		fmt.Println("Hello World!")
	}()

	// Remember, goroutines are considered daemon threads
	// This can exit without printing anything, depending on scheduling
	// semantics!
}

func HelloWorld() {
	fmt.Println("Hello")
}
