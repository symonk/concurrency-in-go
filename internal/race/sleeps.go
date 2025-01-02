package race

import (
	"fmt"
	"time"
)

// AvoidSleeps demonstrates the naive approach people often take when
// starting out to try and avoid a data race.
//
// This example mirrors that of the one in race.go, however it sprinkles
// a time.Sleep between the goroutine creation and the check for number.
//
// This makes the race here significantly harder to detect, but it's still
// there and will present itself later at some point.
//
// Trying to control the scheduler is a complete anti pattern and pain will
// often soon follow.
//
// runtime.GoSched() would be an alternative here, but again do NOT do it!
func AvoidSleeps() {
	var number int
	go func() { number++ }()
	time.Sleep(time.Second)
	if number == 0 {
		fmt.Println("number was 0", number)
	}
}
