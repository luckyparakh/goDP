# A10 · A type whose zero value is usable

## Questions first

Guess before you read on. Wrong guesses are fine. After the task, fill in the
answers and fix your guesses.

1. With the "before" `Counter`, what happens when you write `var c Counter` and then `c.Inc("a")`, and why?
   - Guess:
   - Answer:
2. Name two types in the standard library that work straight from their zero value.
   - Guess:
   - Answer:
3. Why does creating a map "when you first need it" need extra care once goroutines are involved?
   - Guess:
   - Answer:

## The idea

**Problem.** A type that only works after you call its constructor is a trap.
Forget the call, or create it somewhere that does not know about it, and it
fails at run time. The compiler will not warn you.

**Fix.** Design the type so that doing nothing is valid. Make every field usable
from its zero value, and create anything you need (like a map) at the moment you
first use it.

**Analogy (the extension board).** You plug in the new extension board and
nothing works, so you assume it is a bad piece. Nobody mentioned the small
switch on the back that has to be pressed once. A board that just works when
plugged in has nothing to forget.

## How the code shows it

| File | What it is |
|---|---|
| `before.go` | The trap. `Counter` crashes unless you call `NewCounter` |
| `after.go` | Your job: `Tally`, a counter that works as `var t Tally` and is safe for many goroutines |
| `zerovalue_test.go` | `TestBefore_*` pass today. `TestAfter_*` pass when you are done |

Proof: `var t Tally` followed by `t.Inc("a")` just works, and `go test -race`
stays clean.

## Wrong when

The type needs something it cannot invent: a connection string, a required
setting, a dependency. Then a zero value cannot be valid. Use a constructor that
checks its input, so an invalid value can never be built.

## In the wild

- Standard library: `go doc bytes.Buffer` and `go doc sync.Mutex`. Both say what their zero value is. What would using them cost if they did not?
- Source: [src/bytes/buffer.go](https://github.com/golang/go/blob/master/src/bytes/buffer.go) (master may differ from your Go 1.26)
- Book: none verified yet.

## My delete-test result

**Try:** give `Tally` a `NewTally` that makes the map, remove the lazy creation
in `Inc`, and run `TestAfter_ZeroValueIsReady`.

What got harder or impossible? One line:
