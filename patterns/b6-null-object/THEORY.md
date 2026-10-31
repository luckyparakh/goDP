# B6 · Null object: something that is always there

## Questions first

Guess before you read on. Wrong guesses are fine. After the task, fill in the
answers and fix your guesses.

1. In `before.go`, someone adds a new method and forgets one nil check. What happens, and when does anyone find out?
   - Guess:
   - Answer:
2. After the change, where does the one remaining nil check go, and why there?
   - Guess:
   - Answer:
3. `New(nil, nil)` gets the no-op. What does `New(nil, p)` get, when `p` is a nil `*recorder`?
   - Guess:
   - Answer:

## What the principle says

**Null Object**, from Bobby Woolf's chapter in *Pattern Languages of Program
Design 3* (1998): instead of a null reference, give the client an object with
the same interface that does nothing. The client then treats "nobody is
listening" like any other listener, and never checks.

The textbook violation is the one in `before.go`: a null check before every
use, and the one place somebody forgot.

**What changes in Go:** a null object costs almost nothing. An empty struct
with empty methods takes no memory, and the standard library already ships
some: `io.Discard` for writers and `slog.DiscardHandler` for loggers. Go also
has a trap other languages do not: an interface that holds a nil pointer is
not equal to `nil`, so a nil check can let a nil through.

The usable question is: *how many places must remember that this might be
missing?*

## The idea

**Problem.** An optional dependency is `nil` when it is not wanted, so every
use needs a check. The checks bury the real code, and the one that gets
forgotten panics only on the path where the dependency is missing.

**Fix.** Check once, at the edge. The constructor swaps `nil` for an object
that does nothing. From then on the dependency is always there, and nothing
inside checks.

**Analogy (the letterbox).** The courier rings each flat and has to decide what
to do when nobody answers, and the one courier who forgets leaves your packet
on the stairs. Fit a letterbox to every door instead, even the empty flats. Now
no courier asks who is home: everything goes in the box, and an empty flat's
box is just cleared now and then. The one check happens when the box is fitted.

## How the code shows it

| File | What it is |
|---|---|
| `before.go` | The problem. Four nil checks, and a fifth that was forgotten |
| `after.go` | Your job: the same service, starting as a copy, with the checks gone |
| `nullobject_test.go` | `TestBefore_*` pass today. `TestAfter_*` pass when you are done. It also has two benchmarks |

Proof: no field is compared with `nil` in `after.go`, a service built with
nothing never panics, and the benchmark gives you the cost.

## Wrong when

The dependency is required. Missing metrics in production should fail at
startup, not be silently dropped. And a null object for something that returns
data, such as a store that always says "not found", hides real bugs.

## How cheap is the no-op?

Run this and paste your real numbers:

```
go test -bench . -benchmem ./patterns/b6-null-object/
```

- Numbers I measured:
- Is the no-op cheap enough that I would never guard it? One line:

## FAQ

**Why not make `nil` itself the no-op, with `if r == nil { return }` at the top
of each method?**

That is valid Go. A method on a pointer type can be called on a nil pointer,
and if it checks first it is safe. It also fixes question 3's trap for that one
type. But it only works for types you own, and every method needs the guard,
so the checks move into the type instead of disappearing.

## In the wild

- Standard library: `go doc log/slog.DiscardHandler`. Its `Enabled` returns false. Why does that make a discarded log call cheaper than building the message and throwing it away?
- Source: [src/log/slog/handler.go](https://github.com/golang/go/blob/master/src/log/slog/handler.go) (master may differ from your Go 1.26)
- Book: *100 Go Mistakes*, #45 "Returning a nil receiver" ([summary](https://100go.co/#returning-a-nil-receiver-45)). The title is checked; that it fits question 3 is my reading.

## My delete-test result

**Try:** in `New`, remove the line that swaps nil metrics for your no-op, and
run `TestAfter_NilDependenciesAreSafe`.

What got harder or impossible? One line:
