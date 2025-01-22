package main

import "fmt"

func main() {
	go HelloWorld()
}

func HelloWorld() {
	fmt.Println("Hello")
}
