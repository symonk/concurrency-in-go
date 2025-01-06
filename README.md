# concurrency-in-go <!-- omit from toc -->

This repository shares my learning materials from study `concurrency-in-go` written by
[Katherine Cox Buday](https://github.com/kat-co).  The book is excellent and using this
repository along side studying the book can add some extra benefit.

The book in question can be purchased to support Katherine via `oreilly` at [Concurrency in Go](https://www.oreilly.com/library/view/concurrency-in-go/9781491941294/)

-----

# Table of Contents <!-- omit from toc -->

- [:mag\_right: Introduction Materials](#mag_right-introduction-materials)
- [:eyes: Caveats](#eyes-caveats)
- [:tent: Race Conditions](#tent-race-conditions)
- [:tent: Atomicity](#tent-atomicity)
- [:tent: Dead Locking \& Starvation](#tent-dead-locking--starvation)
- [:tent: Communication Sequential Processes](#tent-communication-sequential-processes)
- [:tent: Synchronisation Primities](#tent-synchronisation-primities)
- [:tent: Pipelining](#tent-pipelining)
- [:tent: Patterns](#tent-patterns)

-----

## :mag_right: Introduction Materials

The following materials are critical to understand when learning about concurrency in go:

    - Placeholder
    - Placeholder
    - Placeholder

## :eyes: Caveats

Understanding the golang memory model is important when dealing with concurrency to fully
understand why certain things may cause subtle bugs or be head scratching.  Let's take
one example:

```go
var a, b int

func f() {
	a = 1
	b = 2
}

func g() {
	print(b)
	print(a)
}

func main() {
	go f()
	g()
}
```

Seems relatively straight forward?  It is completely possible here that from the reading perspective
of one goroutine, lets call it goroutine `2` that the writes in another goroutine (`1`) - which invoked
`f()` asynchronously are not visible **OR** the write assigning b is visible, but a is not!

This code can print out `20` where b is `2` and a is `0`, you may be very confused by that as 
b was assigned after a, however for various reasons (and lack of synchronisation) such as the compiler
may even rewrite these lines or the processor may be smart at runtime.

Always use synchronisation primitives when required to ensure correctness.

## :tent: Race Conditions

A `race condition` occurs when code written has a naive expectation on execution
order.  Often a developer expects the code written to execute as it is written.
These kinds of bugs can often be hard(er) to debug and can lie hidden until 
things are scaled up.  

> [!Caution]
> Attempting to manually force goroutine scheduling / context switching is considered
> an anti-pattern and should strongly be avoided.

[Race Conditions: A basic Introduction](datarace_simple/main.go)

[Race Conditions: A Naive Fix](datarace_sleep/main.go)

[Race Conditions: Causing a context switch](datarace_control_scheduler/main.go)

-----

## :tent: Atomicity

`Atomicity` is the concept that something is indivisible or uninterruptable within
a particular `context`.  Context is **very** important here.  Something that is
`atomic` within your process (such as an atomic add leveraging CPU swap instructions)
is not `atomic` in the context of the operating system.

`Performance` plays a vital part in managing the parallelism of code and when using
various primitives to guard against race conditions, a performance penalty must be
considered.

> [!Note]
> Having an 'opt-in' convention for using an API that requires users to remember to
> guard the critical sections is error prone, try and build this into your APIs and
> have function docstrings articulate when this is (or isn't) the case

[Atomicity: A Basic Introduction](atomicity_simple/main.go)

[Atomicity: A Naive Solution](atomicity_naive/main.go)


-----

## :tent: Dead Locking & Starvation

At a basic level, ensuring atomicity with locking critical sections is not the be all and
end all.  All of this can be done however you can still run into other problems, such as
multiple blocks on locks.  Go is not by default re-entrant in terms of mutexes etc so
this is another case of problems that need to be considered.

This section covers case of dead locking, live locking and starvation.

In order to understand where deadlocking can occur, there are a few conditions we can
evaluate,  these are known as the `Coffman Conditions`:


* `Mutual Exclusion`: A concurrent process holds exclusive rights to a resource at any time.
* `Wait-For Condition`: A concurrent process must simultaneously hold a resource and wait for another.
* `No Premption`: A resource held by a concurrent process can only be released by that process itself.
* `Circular Wait`: A concurrent Process (P1) must be waiting on a chain of other concurrent processes
(P2, ...PN), which are in turn waiting on it (P1).

> [!Note]
> Preventing even one of the 4 conditions above, can help prevent deadlocking!

> [!Tip]
> Try to limit the scope of locking to critical sections to start, rather than being broad with locking
> see the starvation example.  It is much easier to widen the locking later, than to reduce it.

[Locking: Deadlock](locking_deadlock/main.go)

[Locking: Livelock](locking_livelock/main.go)

[Locking: Starvation](locking_starvation/main.go)

-----

## :tent: Communication Sequential Processes


-----

## :tent: Synchronisation Primities

-----

## :tent: Pipelining

While pipelining is considered just another pattern, I feel it warrants an individual 
topic of its own.  Pipelining is a concurrent program that has multiple `sequential`
stages that are parallelised internally.  Both a basic and complex example of pipelining
exist in:

 - [Basic Pipeline](basic_pipeline/main.go)
 - [Advanced Pipeline](advanced_pipeline/main.go)

Typically within a pipeline, both the first and last stages have a single entry point
(channel) i.e generating numbers (in) returning the transformed results (out).  Interim
stages typically take an upstream inbound channel and yield their results to and outbound
one.

-----

## :tent: Patterns

A collective of patterns with explanations can be found below:


| Pattern                                                   | Summary                                             |
|-----------------------------------------------------------|-----------------------------------------------------|
| [01 Basic Goroutine](hello_world_goroutine/main.go)       | A simple introduction to goroutines.                |
| [02 Basic Channel](hello_world_chan/main.go)              | A simple introduction to channels.                  |
| [03 Generator](generator/main.go)                         | A python like generator                             |
| [04 Fan In](fanin/main.go)                                | Fan in multiple goroutines                          |
| [05 Restore Sequence](restore_sequence/main.go)           | Fan in multiple goroutines with equal yielding      |
| [06 Select Timeout](select_timeout/main.go)               | Cause a goroutine to terminate conditionally        |
| [07 Quit Signal](quit_signal/main.go)                     | Cancel a goroutine with an channel send             |
| [08 Daisy Chain](daisy_chain/main.go)                     | A simulation of chinese whispers with goroutines    |
| [09 Basic Pipeline](basic_pipeline/main.go)               | A simple mathematical example of pipelining         |
| [10 Advanced Pipeline](advanced_pipeline/main.go)			| A smarter parallel pipeline						  |

-----
