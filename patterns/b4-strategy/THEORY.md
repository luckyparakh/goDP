# B4 · Strategy: hand over the how, not its name

## Questions first

Guess before you read on. Wrong guesses are fine. After the task, fill in the
answers and fix your guesses.

1. To add a "jitter" policy in `before.go`, what must you edit, and what happens if a caller misspells it?
   - Guess:
   - Answer:
2. A2 also replaced a `switch` with funcs. What is different here?
   - Guess:
   - Answer:
3. After the change, who chooses which backoff is used, and where would that happen in a real service?
   - Guess:
   - Answer:

## What the principle says

**Strategy**, from *Design Patterns* by Gamma, Helm, Johnson and Vlissides
(1994): define a family of algorithms, put each one behind the same interface,
and let the client be configured with one of them. The textbook cases are
sorting and line-breaking algorithms.

The textbook violation is the one in `before.go`: a `switch` on a mode or type
name that grows a case for every new variant.

**What changes in Go:** a strategy is usually just a func value. A constructor
returns a closure that has captured its own settings, such as
`Exponential(base, limit)`. The client holds the func and calls it, and never
learns which strategies exist. A9 covered when an interface is the better
shape.

The usable question is: *does the caller need to know every option, or only
to run the one it is given?*

Cross-principle view: [SOLID.md](../SOLID.md), the "O" row.

## The idea

**Problem.** The caller picks an algorithm by name through a growing `switch`.
A new one means editing working code, and a misspelt name compiles and
silently gets the default.

**Fix.** Make the algorithm a value. The caller is handed one and runs it.
New algorithms are new funcs, written anywhere, and a misspelling is a compile
error.

**Analogy (the route).** The driver knows routes only by name, so a new
shortcut means explaining it first, and when you say a name the driver hasn't
heard, you get the usual route and notice halfway. Send the route on your phone
instead. The driver still drives, but which way to go is yours to hand over,
and any route works, even one the driver has never taken. A route on a map
can't be misheard.

## How the code shows it

| File | What it is |
|---|---|
| `before.go` | The problem. `Delay` is a `switch` on a policy name |
| `after.go` | Your job: `Backoff`, three constructors, and `Retry` |
| `strategy_test.go` | `TestBefore_*` pass today. `TestAfter_*` pass when you are done |

Proof: a Backoff written only in the test file works with `Retry`, and
`Retry` is not edited.

## Wrong when

There are two options that never change, chosen in one place. A `switch` you
can read top to bottom is then clearer than funcs spread around.

## FAQ

**The config file says "exponential". Don't I need the string anyway?**

Yes, once, at the edge. `main` reads the config and turns the name into a
`Backoff`. Below that, nothing sees the string. The `switch` does not vanish,
it moves to the one place that has to deal with names, and a typo there fails
at startup instead of in the middle of a retry.

## In the wild

- Standard library: `go doc slices.SortFunc`. It never knows what order you want. Where does that decision live, and what would the version with a mode string look like?
- Source: [src/slices/sort.go](https://github.com/golang/go/blob/master/src/slices/sort.go) (master may differ from your Go 1.26)
- Book: none verified yet.

## My delete-test result

**Try:** delete `Backoff`, make `Retry` take a policy name and call `Delay`
again, then try to make `TestAfter_FourthStrategyNeedsNoEdit` pass.

What got harder or impossible? One line:
