# BLUEPRINT: the 43 tasks

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

## Capstones: one per block

Separate tasks do not teach assembly. You can finish ten patterns and still not
know how they fit together. Each block ends with one small program that only
works if most of that block works.

A capstone is scaffolded differently from a pattern:

```
capstones/<id>/README.md        the proofs, written as commands
capstones/<id>/<name>_test.go   the same proofs, as tests that fail
capstones/<id>/types.go         interfaces and structs, declared and empty
capstones/<id>/fake.go          a fake for the slow dependency, so tests cost nothing
capstones/<id>/*.go             stubs: panic("write me")
```

**No sealed solution.** A pattern has one shape worth copying; a capstone has
many. A sealed answer here would say there is one right design, and there is not.

**Three capstones, not six.** Not every block has enough meat for its own, and a
capstone invented to fill a slot is decoration.

| # | After | Name | What it is |
|---|---|---|---|
| **K1** | A + B + C | **The service** | One HTTP service: interfaces owned by the caller, a swappable store behind an adapter, a middleware chain, everything wired in `main`, nothing constructible in an invalid state |
| **K2** | D | **The keeper** | A program that keeps something the way you want it. Defined below |
| **K3** | E + F | **The restructure** | Take K1 back out. Add explicit states, a handler chain and a cache in front of the slow step. Then do F4's hexagonal restructure on it |

Why these groupings:

- **A cannot stand alone.** Its patterns are all shape, with no runtime behaviour.
  They only prove themselves once something uses them.
- **C is three tasks**, one of them written-only. Folding it into K1 is honest;
  giving it a capstone is ceremony.
- **D stands alone.** Twelve tasks, and the only block about behaviour over time.
- **F4 already is a capstone.** Pointing it at K1 means K1 was never throwaway,
  and E's patterns are what you add while restructuring it.

**You build three things, but only two are new.** K3 is K1 revisited months later
by a better engineer. That is the one worth showing somebody.

**Only K2 is designed below.** Designing K1 and K3 before you have done their
blocks is planning, and planning is cheaper than building, which is why it wins.
One sentence each when you get there.

**Cost, estimated and not measured.** Roughly 78 h for the 39 remaining tasks and
40 to 46 h for the three capstones, after subtracting F4 which is already counted
as a task. Call it 120 to 125 h: about 13 weeks at 10 h/week, or 21 at 6. The
capstones are half the budget. That is the point of them, not an overrun.

### K2 — the Block D capstone

You write down what you **want** in a small config file. The program looks at
what is **actually** there. Where they differ, it makes one slow call that
sometimes fails, to close the gap. It runs in a loop, and **you can kill it at
any moment and restart it without breaking anything.**

That last sentence is the point. A demo runs once on a good day. This is defined
by what happens when things go wrong.

The slow call is anything slow and unreliable: an HTTP fetch, an image resize, a
shell command, a model. A fake is the default, so the tests need no network and
no account.

Pick the thing that annoys you most in real life. Every folder in your notes has
a current summary · every saved link has a title and tags · every receipt in
Downloads is renamed and filed. The shape is the same in all of them.

| What the program forces you to get right | Task |
|---|---|
| Restart safely at any moment: compare wanted against actual, do not react to events | D11 |
| The call is expensive and sometimes times out: never redo finished work | D9 |
| Many items, but the far end has a limit: a fixed number of workers | D2 |
| Ctrl-C finishes cleanly; every call has a deadline | A7 |
| Prove the workers actually stopped | D1 |
| "Busy, retry" and "bad input, stop" are different | A6 |
| Two runs writing the same state file | D10 |
| Two copies running, only one allowed to act | D12 |
| Swap the slow thing for a fake | A5 |

Nine for nine. If a piece could be removed and it still worked, it was decoration.

**Done means these. Each is a command, not a judgement:**

1. Kill it halfway, restart: it finishes what is left and does **not** redo what is done
2. Run two copies at once: only one does work
3. Make the fake fail half the time: it still finishes, with no duplicates
4. Delete one output: the next pass regenerates **only that one**
5. Change one input: **only that one** is redone
6. `go test -race` clean; the goroutine count returns to where it started
7. The whole suite runs against the fake: no network, no account, no cost

My estimate is 18 to 20 hours for all seven. The smallest honest version is about
6 hours — config, compare, fake, idempotency, one worker, no leader election —
which already passes proofs 1, 4 and 5, and is already more than most demos do.
Grow it one proof at a time. Not measured.

**Do not design it before D12 passes.** Write one sentence naming what it keeps
in order, and nothing else.

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

B1 to B6 are ready. B1, B2 and B5 have `before/` and `after/` subfolders, because
they are about imports between packages: run them with `go test ./patterns/<id>/...`.

## Block C: construction

| ID | Task | Proof |
|---|---|---|
| C1 | `NewX` with validation: invalid state cannot be built | Bad config returns an error, never a half-built struct |
| C2 | Functional options | Adding a 5th option leaves the old call site compiling |
| C3 | Options vs a config struct: a written decision for 3 cases, no code | Written justification (a config struct under 3 options; options for 3+ that may grow) |

## Block D: concurrency (the critical path for services)

| ID | Task | Proof |
|---|---|---|
| D1 | A goroutine-leak harness, reused in D2 to D12 | `runtime.NumGoroutine()` before and after, with a settle period |
| D2 | A bounded worker pool | An atomic gauge: max in flight never exceeds N |
| D3 | Pipeline with fan-out and fan-in | `-race` clean, and every stage exits on cancel |
| D4 | `errgroup.WithContext` | The first error cancels siblings |
| D5 | Graceful shutdown of an HTTP server | An in-flight request completes and new ones are rejected |
| D6 | singleflight | 50 identical concurrent calls make exactly 1 downstream call |
| D7 | A rate limiter (token bucket) | The rate is not exceeded over a measured window |
| D8 | A circuit breaker: closed, open, half-open | A table-driven transition test |
| D9 | Retry, backoff, jitter and an idempotency key | No double write after a retried, ambiguous failure |
| D10 | Optimistic concurrency: a version, not a lock | A stale write returns `ErrConflict`; no update is silently lost |
| D11 | Reconcile: compare desired to actual, do not react to events | State is still correct after an event is dropped |
| D12 | Leader election with an expiring lease | One worker acts per tick; another takes over when the holder stops renewing |

D1 comes first. Without a leak test every later proof here is a guess.
`errgroup` and `singleflight` live in `golang.org/x/sync`, which is not in the
standard library.

D10 to D12 are the three ideas any system uses to coordinate more than one copy
of itself: a version check, a loop that closes the gap, and a lease that expires
on its own. Plain Go, standard library only, nothing to install. Do them in
order — D12's failover test is easier once D10 has taught you what a conflict is.

Each one starts from a bug you can watch happen: D10 loses an update, D11 is
wrong forever after one missed event, D12 does the work three times. The failing
behaviour is the specification.

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
- D10, D11 and D12 were added on 29/10/2026, taking the count from 40 to 43.
  Block D stopped at retry and idempotency, which left out the three things a
  service needs once there is more than one copy of it running. My estimate is
  2 to 3 hours for D10 and 4 to 5 each for D11 and D12. Not measured.

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
| Two writers, last one wins, and nothing returns an error | Optimistic concurrency: a version on the record |
| A handler that patches state in response to each event | Reconcile: compare desired to actual instead |
| Every replica running the same periodic job | Leader election: one lease, one actor |
| **An interface with one implementation and one caller** | **Nothing. Delete the interface.** |
