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
- [:tent: Concurrency Building blocks](#tent-concurrency-building-blocks)
	- [:one: Goroutines](#one-goroutines)
	- [Context Switching](#context-switching)
- [:tent: Synchronisation Primities](#tent-synchronisation-primities)
	- [:one: Sync Package](#one-sync-package)
	- [:two: Placeholder](#two-placeholder)
- [:tent: Pipelining](#tent-pipelining)
- [:tent: Patterns](#tent-patterns)

-----

## :mag_right: Introduction Materials

The following materials are critical to understand when learning about concurrency in go:

    - Placeholder
    - Placeholder
    - Placeholder

## :eyes: Caveats

> [!Important]
> Please read [Go Memory Model](https://go.dev/ref/mem)

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

Go is modelled (but not entirely) on the ideas of Tony Hoare's 
`CSP` (Communication Sequential Processing).

[CSP](https://en.wikipedia.org/wiki/Communicating_sequential_processes)

>[!Note]
> Don't communicate by sharing memory, share memory by communicating!

>[!Note]
> Go channels/select are powerful CSP primitives, go still offers 
> typical mutexes etc via the sync package.


-----

## :tent: Concurrency Building blocks

### :one: Goroutines

The `goroutine` is the core building block behind go's excellent concurrency model.  A goroutine is a function that
is running `concurrently` (maybe in parallel).

> [!Note]
> A goroutine is not guaranteed to be running in parallel, you may have a single core machine!

 - [Hello World Goroutine](goroutine/main.go)

 The `goroutine` in go is **NOT** an OS thread, but is also not considered a `green` thread (A thread that
 is managed by the languages runtime).  They are infact a higher level of abstraction known as `coroutines`.

 A `coroutine` is simply a concurrent subroutine (think function, method, closures etc) that are not 
 `preemptive` (they cannot be interrupted).  Instead there are various points in their lifecycle where
 they can be stopped or enables re-entry.

 One of the beauties of go is that the goroutines are heavily managed by the go runtime and is abstracted
 away from the users.  The go runtime can inspect goroutines at runtime to detect operations that would
 block (and resume) and suspend/re-enter respectively.

 Go's mechanism for hosting goroutines is an implementation of whats called an `M:N scheduler`.  This means
 it maps `M` _green-threads_ -> `N` _operating system threads_.  When more goroutines are spawned than there
 are _green threads_ the go scheduler will handle the distribution of the goroutines across the available
 threads and for ones that block, ensure others can run.

 > [!Note]
 > More on the go scheduler later, see the topic -> The Go Runtime Scheduler.

 > [!Note]
 > Go uses the fork-join model of concurrency.

 > [!Caution]
 > Time.Sleep does NOT create a join point, avoid attempting to use it as one.

 - [A simple join point](fork_join_point/main.go)

`Goroutines` operate in the same address space in which they were created, this means they can 
access / read variables in their scope when running closures for exmaple:

 - [Goroutine Address Space](goroutine_address_space/main.go)

 Go's compiler will take care of keeping variables in scope to avoid goroutines accessing garbage
 collected memory.

`Goroutines` are extremely efficient, a fresh goroutine is only given a few kilobytes (2KB at
the time of writing this) which in most cases is often enough for the task.  If it isn't  the
go runtime will grow (and shrink) the memory for storing the goroutines stack automatically.

This allows go programs to use very little memory in comparison and it is practical to create
hundreds if not thousands of goroutines in the same address space.  If goroutines were threads
the system would be completely overloaded much faster.

 - [Goroutine Bootstrap Memory](goroutine_memory/main.go)

 In the example above, we created 10k goroutines on our system, using a simple 64 bit CPU with 32GB
 of memory.  The 10k routines took up a total of `2.59KB` per goroutine.  Without using swap space
 on this system, it would in theory be possible (if their stacks didn't need to grow ofcourse) to
 spawn a total of `2^5 (~32GB)` of ram => `12_350_000+` (yes, 12.3 **million**) goroutines without
 using swap space!

 > [!Caution]
 > Because you can, doesn't mean you should! Switching between this many routines will have a heavy penalty!

 ### Context Switching

 In the world of os threads, context switching can be pretty costly.  It's not to say that it also cannot hurt
 goroutines when the coroutines must be suspended/reentered etc, but its another win for go in that it is less
 costly.

 Here we can see a performance increase of over `90%` compared to switching OS threads on my machine on linux:

 - [Goroutine Context Switching Performance](goroutine_context_switching/main_test.go)
  
### :two: Waitgroups

A waitgroup is an atomic counter that allows waiting for a collection of goroutines to finish.
This is useful only if you do not care about routine return values, or if you do you have another
mechanism for collecting them:

 - [Sync WaitGroup](waitgroup/main.go)
  
### :three: Mutexes

Another synchronisation primitive, most familiar to those from other languages that handle
synchronisation outside of CSP.

 - [Sync Mutex & RWMutex](mutexes/main.go)

-----

## :tent: Synchronisation Primities

### :one: Sync Package

The `sync` package offers primarily low level synchronisation primitives.
Some of the common usages of the `sync` package are `WaitGroup`, various
flavours of `Mutex (Locker)` types and the `Once` variants aswell as other
low level primitivies.

 - [sync.OnceFunc](sync_package/main.go)
 - [sync.OnceValue](sync_package/main.go)
 - [sync.OnceValues](sync_package/main.go)
 - [sync.Cond](sync_package/main.go)

### :two: Placeholder


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
