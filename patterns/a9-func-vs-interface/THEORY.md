# A9 · A function or a one-method interface?

## Questions first

Guess before you read on. Wrong guesses are fine. After the task, fill in the
answers and fix your guesses.

1. In `before.go`, how many lines does it take to add one new rule, such as "divisible by 3"? What would a plain function need?
   - Guess:
   - Answer:
2. What can an interface do that a plain function cannot?
   - Guess:
   - Answer:
3. You have a function, but the API asks for the interface. What small thing bridges the two?
   - Guess:
   - Answer:

## The idea

**Problem.** One decision is being passed around, and it is wrapped in a type
with a method just so it can be passed. For a one-line rule that is more
ceremony than the rule itself.

**Fix.** Offer the plain function form for the common case, and add a small
adapter, a function type with the one method, so a function can still be used
wherever the interface is expected. Then decide which form you would ship.

**Analogy (one errand).** You need one packet of milk, so you interview and hire
a full-time cook for it. That is a type with a method for a one-line rule. A chit
sent with whoever is going out would have done it. When the job really is every
meal, and it has to remember what is already in the fridge, then you do want the
cook.

## How the code shows it

| File | What it is |
|---|---|
| `before.go` | The interface form. Each rule is a named type with a method |
| `after.go` | Your job: `FilterFunc` and the `RuleFunc` adapter |
| `funciface_test.go` | `TestBefore_*` pass today. `TestAfter_*` pass when you are done. It also has two benchmarks |

Proof: both forms give the same result, and the benchmarks give you numbers to
compare.

## Wrong when

The thing has several methods, needs identity, or has to be inspected. Then it
needs an interface. A function is the wrong shape for that.

## Which would I ship?

Run this and paste your real numbers:

```
go test -bench . -benchmem ./patterns/a9-func-vs-interface/
```

- Numbers I measured:
- The one I would ship, and why (one or two lines):

## In the wild

- Standard library: `go doc net/http.HandlerFunc`. How does that one method let any function be an HTTP handler?
- Source: [src/net/http/server.go](https://github.com/golang/go/blob/master/src/net/http/server.go) (master may differ from your Go 1.26)
- Book: *100 Go Mistakes*, #5 "Interface pollution" ([summary](https://100go.co/#interface-pollution-5)). The title is checked; that it fits this lesson is my reading.

## My delete-test result

**Try:** remove the `RuleFunc` adapter, then try to pass the closure
`func(n int) bool { return n%3 == 0 }` to `Filter`.

What got harder or impossible? One line:
