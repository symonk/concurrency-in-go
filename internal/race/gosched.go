package race

import (
	"fmt"
	"runtime"
)

// TODO: This isn't reliable, Update this.
//
// ControllingContextSwitchingDemo demonstrates the capabilities of runtime.GoSched() to
// cause context switching.
//
// It sets the go max procs (cpu cores) to 1, this would prevent additional
// goroutines from running concurrently.  We can force it to happen with
// setting runtime.GOMAXPROCS(1) and then calling runtime.GoSched() to
// cause a context switch.
func ControllingContextSwitchingDemo() {
	runtime.GOMAXPROCS(1) // Only allow a single core of a (potentially) multi-core CPU.
	go func() { fmt.Println("Inside Goroutine") }()
	fmt.Println("Hello World")
}
