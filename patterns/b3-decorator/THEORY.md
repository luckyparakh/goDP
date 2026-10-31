# B3 · Decorator: write it once, wrap every route

## Questions first

Guess before you read on. Wrong guesses are fine. After the task, fill in the
answers and fix your guesses.

1. In `before.go`, how many places do you edit to add request timing to every route, and what did `/export` get wrong?
   - Guess:
   - Answer:
2. In a chain of Logging, then Auth, is a request that Auth rejects still logged? What if you swap the two?
   - Guess:
   - Answer:
3. The log line needs the status code, but the inner handler is the one that writes it. How can the logging wrapper see it?
   - Guess:
   - Answer:

## What the principle says

**Decorator**, from *Design Patterns* by Gamma, Helm, Johnson and Vlissides
(1994): attach a responsibility to an object by wrapping it in another object
with the same interface. The textbook case puts a border or a scroll bar around
a window, without the window knowing.

**What changes in Go:** no class hierarchy is needed. A decorator is a function
that takes an `http.Handler` and returns an `http.Handler`, usually a closure.
Because a wrapper is the same type as what it wraps, wrappers stack.

Once they stack, **order is the design**. The outermost wrapper sees the
request first and the response last, and any wrapper can stop the request
before it reaches the ones inside.

The usable question is: *is this the same for every route? Then it is a
wrapper, not a copy.*

Cross-principle view: [SOLID.md](../SOLID.md), the "O" row: a new job is a new
wrapper, and no handler is edited.

## The idea

**Problem.** The same jobs are pasted into every handler. A change means
editing all of them, and one copy quietly drifts.

**Fix.** Write each job once as a wrapper, and put one chain of wrappers in
front of all the routes. The order of the chain is written in one place.

**Analogy (the society gate).** Every flat in the building does its own visitor
check, and one flat forgot, so anyone who says "plumber" walks straight in.
Move the check to the gate. The guard writes the name in the register, then
rings the flat, then lets the visitor up, so one check covers every flat,
including the one that opens next month. Because the register comes first, it
also shows the visitors a flat turned away.

## How the code shows it

| File | What it is |
|---|---|
| `before.go` | The problem. Three handlers with the same pasted code, one of them missing auth |
| `after.go` | Your job: `Chain`, `Logging`, `Auth`, `Timing` and `Routes` |
| `decorator_test.go` | `TestBefore_*` pass today. `TestAfter_*` pass when you are done |

Proof: an execution-order test shows `A> B> C> h <C <B <A`, and a rejected
request still appears in the log.

## Wrong when

The job differs per route, or there is only one route. A shared wrapper then
hides a one-off behind indirection.

## Read three

The standard library has three decorators in one file. Read each with
`go doc -src`, and write one line: what it does before calling your handler,
what it does after, and what it wraps besides the handler.

- `net/http.StripPrefix`:
- `net/http.TimeoutHandler`:
- `net/http.MaxBytesHandler`:

## In the wild

- Standard library: `go doc net/http.StripPrefix`. It takes a `Handler` and returns a `Handler`. What does it change about the request before yours sees it?
- Source: [src/net/http/server.go](https://github.com/golang/go/blob/master/src/net/http/server.go) (master may differ from your Go 1.26)
- Book: none verified yet.

## My delete-test result

**Try:** in `Routes`, swap `Logging` and `Auth` in the chain, and run
`TestAfter_RejectedRequestIsLogged`.

What got harder or impossible? One line:
