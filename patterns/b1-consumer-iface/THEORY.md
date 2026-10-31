# B1 · The interface lives in the consumer's package

## Questions first

Guess before you read on. Wrong guesses are fine. After the task, fill in the
answers and fix your guesses.

1. In `before/`, billing imports store. You move `Invoice` into billing, where it belongs. What breaks?
   - Guess:
   - Answer:
2. In A5 the provider imported nothing. Why must store import billing here?
   - Guess:
   - Answer:
3. When should an interface live with neither the consumer nor the provider?
   - Guess:
   - Answer:

## What the principle says

**Dependency Inversion Principle** at package scale (A5 did it inside one
package), together with Robert C. Martin's **Acyclic Dependencies Principle**:
the dependency graph of packages must have no cycles. Martin states it in
*Clean Architecture* (2017); when he first wrote it down is not verified here.

The textbook violation is a cycle: A imports B and B imports A, so neither can
be built, tested or released without the other.

**What changes in Go:** the compiler enforces it. A cycle is not a smell, it is
a build error: `import cycle not allowed`. So when a provider needs the
consumer's types, there is only one legal place for the interface, and that is
with the consumer.

A5 showed the arrow disappearing, because its provider returned a plain
string. Here the provider returns the consumer's own type, so the arrow has to
turn round: store imports billing.

The usable question is: *if the provider needs my types, which way must the
import go?*

Cross-principle view: [SOLID.md](../SOLID.md).

## The idea

**Problem.** The code that decides imports the code that stores, only because
the storage package declares the shared type and the interface. The day the
type moves home, the build breaks with a cycle.

**Fix.** The consumer declares its own type and the small interface it calls.
The provider imports the consumer and returns the consumer's type. The
consumer imports nothing of the provider's.

**Analogy (fee challan).** The college won't confirm your seat without the
bank's fee receipt, and the bank won't take the fee without a letter from the
college saying which seat it is for. You go back and forth between two counters
and nothing moves. The college fixes it by printing its own fee challan with
your seat on it. The bank stamps the college's challan and hands it back, so
the bank works on the college's paper and the college never needs the bank's.

## How the code shows it

| File | What it is |
|---|---|
| `before.go` | The problem, in comments, and a TRY IT that ends in an import cycle |
| `before/billing`, `before/store` | billing imports store. Do not edit |
| `after.go` | The GOAL |
| `after/billing`, `after/store` | Your job. They start as a copy of `before/` |
| `consumer_test.go` | `TestBefore_*` pass today. `TestAfter_*` pass when you are done |

Proof: `go list -deps` shows store depending on billing, and billing depending
on nothing of store's.

## Wrong when

There is one consumer, one implementation and no fake, and the provider uses
only plain types. Then you need neither the package split nor the interface.

## FAQ

**A5's FAQ said Go deletes the arrow instead of reversing it. Here it
reverses. Which is it?**

Both, depending on what the provider hands back.

| The provider returns | The arrow |
|---|---|
| Plain types (`string`, `int64`) | Disappears. Neither package imports the other |
| The consumer's own types (`billing.Invoice`) | Reverses. The provider imports the consumer |

The second is the usual case in a real service, because the provider almost
always returns the domain's records.

**Where should `ErrNotFound` go?**

It stays in store for now, so billing can only pass it up without naming it.
Deciding what an error means, and translating it at the boundary, is B2.

## In the wild

- Standard library: `go doc database/sql/driver.Driver`. `database/sql` uses drivers. Which package declares the interface, and which side imports which?
- Source: [src/database/sql/driver/driver.go](https://github.com/golang/go/blob/master/src/database/sql/driver/driver.go) (master may differ from your Go 1.26)
- Book: *100 Go Mistakes*, #6 "Interface on the producer side" ([summary](https://100go.co/#interface-on-the-producer-side-6)) and #7 "Returning interfaces" ([summary](https://100go.co/#returning-interfaces-7)). The titles are checked; that they fit this lesson is my reading.

## My delete-test result

**Try:** in `after/`, move the interface back into store while `Invoice` stays
in billing, then run `go build ./patterns/b1-consumer-iface/...`.

What got harder or impossible? One line:
