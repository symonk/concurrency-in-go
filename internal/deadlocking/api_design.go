package deadlocking

// A dummy Pi struct
type Pi struct{}

// CalculateDigitsOfPiBad is a bad example of calculating the digits of Pi.
func CalculateDigitsOfPiBad(start, finish int64, pi *Pi) {
	/*
		The documentation of this function serves relatively little purpose.
		The implementation utilises concurrency within, but to the caller they
		are not clearly aware and may in their code, add additional guards etc.

		This can raise some questions like:

		How do I calculate Pi with this function?
		Should i run multiple instances myself or is it concurrent internally?
		Should i guard my pi instances internally, are multiple routines accessing/writing to it?
	*/

}

// CalculateDigitsOfPiGood calculates the digits of pi between start and finish.
//
// Internally, CalculateDigiestsOfPiGood will create Floor((end-begin)/2) concurrent
// processes which recursively call CalculateDigitsOfPiGood.  T
//
// Synchronization of writes to the pi instance are handled internally by the Pi struct
func CalculateDigitsOfPiGood(start, finish int64, pi *Pi) {
	/*
		A slightly better docstring, makes it clear the responsibilities of the func
		and answers any questions outlined above in the bad example.
	*/
}

// CalculatePiBest calculates the digits of pi between start and finish.
// It returns a channel that will be populated with the calculated pi digits.
func CalculatePiBest(start, finish int64) <-chan int {
	/*
		An even simpler example that indicates the implementation takes care
		of concurrency and provides a channel to interact with the calculated
		pi digits.
	*/
	return make(chan int)
}
