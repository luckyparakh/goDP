# BLUEPRINT: the 40 tasks

How to do a task is in [HOWTO.md](HOWTO.md). This file is only the list.

One task is one folder under `patterns/`, with the same five files every time:

```
patterns/<id>/THEORY.md      questions, theory, analogy, "wrong when"
patterns/<id>/before.go      the problem, explained in comments
patterns/<id>/after.go       a stub: you fill it in
patterns/<id>/<name>_test.go the proof: Before passes, After fails until done
solutions/<id>/after.go      sealed: the finished code and why it works
```

A pattern counts as learned when a test proves it. After it works you delete it
and write down what broke. If nothing broke, it was decoration.

## Block A: foundations

| ID | Task | Proof |
|---|---|---|
| A1 | SRP: one job per piece | The sum is tested with no temp folder |
| A2 | OCP: add a case without editing working code | A new tier is 1 file added, 0 modified |
| A3 | LSP: same shape is not the same promise | One checklist catches a lying store |
| A4 | ISP: depend on what you use | A test stand-in needs 1 method, not 10 |
| A5 | DIP: the user of a thing says what it needs | The failure path is one line in a test |
| A6 | Error taxonomy: sentinel, typed, wrapped | `errors.Is` / `errors.As` work through 2 wrap layers |
| A7 | Context: propagate, timeout, cancel across 3 call levels | `ctx.Err()` is `DeadlineExceeded` |
| A8 | Embedding is not inheritance | A test that breaks your OO intuition |
| A9 | The same job as a func type and as a 1-method interface | A benchmark and a written "which would I ship?" |
| A10 | A type whose zero value is usable | `var c Counter` works with no constructor |

A1 to A10 are ready. A8 kills the most Java habits, so do not skip it.

A1 to A5 are the five SOLID principles, one task each. Each task explains its own
principle, so start with the task. For the view across all five, which is where
the blast-radius framing lives, see [patterns/SOLID.md](patterns/SOLID.md).

## Block B: the daily core

| ID | Task | Proof |
|---|---|---|
| B1 | The interface lives in the consumer's package | `go list -deps` shows the dependency arrow reversed |
| B2 | Adapter and error translation (driver errors to domain errors) | The service compiles with zero `database/sql` imports |
| B3 | Decorator: a middleware chain (log, auth, timing) | An execution-order test, then read 3 stdlib decorators |
| B4 | Strategy: replace a growing `switch` with a func value | Add a 4th strategy without touching the caller |
| B5 | Composition root: wire handler, service and store in `main` | `grep` shows no package builds its own dependencies |
| B6 | Null object: a no-op logger or metrics recorder | Zero nil checks, and a benchmark shows the no-op is cheap |

## Block C: construction

| ID | Task | Proof |
|---|---|---|
| C1 | `NewX` with validation: invalid state cannot be built | Bad config returns an error, never a half-built struct |
| C2 | Functional options | Adding a 5th option leaves the old call site compiling |
| C3 | Options vs a config struct: a written decision for 3 cases, no code | Written justification (a config struct under 3 options; options for 3+ that may grow) |

## Block D: concurrency (the critical path for services)

| ID | Task | Proof |
|---|---|---|
| D1 | A goroutine-leak harness, reused in D2 to D9 | `runtime.NumGoroutine()` before and after, with a settle period |
| D2 | A bounded worker pool | An atomic gauge: max in flight never exceeds N |
| D3 | Pipeline with fan-out and fan-in | `-race` clean, and every stage exits on cancel |
| D4 | `errgroup.WithContext` | The first error cancels siblings |
| D5 | Graceful shutdown of an HTTP server | An in-flight request completes and new ones are rejected |
| D6 | singleflight | 50 identical concurrent calls make exactly 1 downstream call |
| D7 | A rate limiter (token bucket) | The rate is not exceeded over a measured window |
| D8 | A circuit breaker: closed, open, half-open | A table-driven transition test |
| D9 | Retry, backoff, jitter and an idempotency key | No double write after a retried, ambiguous failure |

D1 comes first. Without a leak test every later proof here is a guess.
`errgroup` and `singleflight` live in `golang.org/x/sync`, which is not in the
standard library.

## Block E: structure

| ID | Task | Proof |
|---|---|---|
| E1 | Template method as a func parameter | Two algorithms share one skeleton |
| E2 | Chain of responsibility with early exit | The downstream call count is 0 after a rejection |
| E3 | Facade and `internal/` | Importing `internal/` from outside fails to compile |
| E4 | A caching proxy over a slow dependency | 3 calls make 1 downstream call |
| E5 | An explicit state machine (a map of state to event to state) | An illegal transition is rejected |
| E6 | Observer: callbacks vs channels | A leak test on an unsubscribed observer |

## Block F: pro and architecture

| ID | Task | Proof |
|---|---|---|
| F1 | Generics: a typed cache or store | Compile-time type safety, no `any` leaking out |
| F2 | Fakes over mocks | The contract suite from A3 is reused for the fake |
| F3 | Table-driven tests and golden files | One table, 12 cases, readable failure output |
| F4 | Hexagonal restructure of one real service (the capstone) | The domain package imports zero infrastructure |
| F5 | A package-boundary audit script | CI fails if the domain package imports `database/sql` |
| F6 | A review drill on a real open-source Go repo | Written: one missing pattern, one over-applied |

## Further reading (verified entries only)

*100 Go Mistakes*, with a free summary at [100go.co](https://100go.co). Titles
and numbers are checked against that page. Which task each one fits is my reading.

| Task | Entry |
|---|---|
| A2 | #3 Misusing init functions |
| A4 | #5 Interface pollution |
| A5 | #6 Interface on the producer side |
| A6 | #48 to #54 (Panicking, up to Not handling defer errors) |
| A7 | #60 Misunderstanding Go contexts, and #61 Propagating an inappropriate context |
| D1 | #62 Starting a goroutine without knowing when to stop it |
| D (race checks) | #58 Not understanding race problems |

Not mapped yet: A1, A3, and the other two books. The chapter list for
*Concurrency in Go* in the old handoff was not re-checked.

## What changed from the original 39-task plan

- The old A1 was five labs in one task. It is now A1 to A5, one pattern each.
- Three old tasks are gone because A1 to A5 already teach them: B1 (ISP), B4 (LSP)
  and C4 (registry, which is A2).
- The old B2 (DIP) stays as B1 here. It is about the *import direction across
  packages*, which A5 does not cover.
- The old estimate of 40 to 60 hours was for 39 tasks and has not been
  re-measured.

## Know the name, rarely write it in Go

Abstract Factory, Prototype, Bridge, Mediator, Memento, Interpreter. Visitor is
a type switch. Flyweight is `sync.Pool`. Singleton is dependency injection (or
`sync.OnceValue` if unavoidable). Command is a closure. Builder is a struct
literal or functional options. Template Method is a func you pass in.

## The reviewer's lens

| Symptom in the code | Likely missing pattern |
|---|---|
| `sql.ErrNoRows` or driver error codes outside the DB layer | Adapter and error translation |
| A growing `switch` on a type or mode string | Strategy, or a registry map |
| The same logging, auth or timing copied across 8 handlers | Decorator or middleware |
| A package creating its own DB or HTTP client | DI and a composition root |
| A constructor with 7 parameters, 5 optional | Functional options |
| Nil checks on a logger everywhere | Null object |
| N identical concurrent requests hitting the DB | singleflight or a cache |
| An unbounded `go func()` inside a loop | Worker pool |
| Retry wrapped around a write with no key | An idempotency key (fix that before retrying) |
| **An interface with one implementation and one caller** | **Nothing. Delete the interface.** |
