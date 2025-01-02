# concurrency-in-go

This repository shares my learning materials from study `concurrency-in-go` written by
[Katherine Cox Buday](https://github.com/kat-co).  The book is excellent and using this
repository along side studying the book can add some extra benefit.

The book in question can be purchases to support Katherine via `oreilly` at [Concurrency in Go](https://www.oreilly.com/library/view/concurrency-in-go/9781491941294/)

-----

### :tent: Table of Contents



-----

### :tent: Race Conditions

A `race condition` occurs when code written has a naive expectation on execution
order.  Often a develop expects the code written to execute as it is written.
These kinds of bugs can often be hard(er) to debug and can lie hidden until 
things are scaled up.  

> [!Caution]
> Attempting to manually force goroutine scheduling / context switching is considered
> an anti-pattern and should strongly be avoided.

[A basic data race](internal/race/race.go)

[A Naive Fix](internal/race/sleeps.go)

[Controlling Context Switching Manually](internal/race/gosched.go)

-----

### :tent: Atomicity

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

[Understanding simple atomicity](internal/atomicity/simple_increment.go)

[Ensuring atomicity (Naive)](internal/atomicity/ensuring_atomicity.go)


-----

### :tent: Dead Locking & Starvation

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

[Deadlocking Mutexes (Coffman Conditions Explained)](internal/deadlocking/deadlock.go)

[Livelocking](internal/deadlocking/livelock.go)

[Starvation](internal/deadlocking/starvation.go)

Smart abstractions and documentation are **vital** when concurrency is involved.  An example of how to
make things easier for developers consuming (or maintaining) your code in future is displayed in the
`api_design.go` file:

[Smart Concurrency API Design](internal/deadlocking/api_design.go)

-----

### Communication Sequential Processes



-----
