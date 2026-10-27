# A3 · LSP: same shape is not the same promise

## Questions first

Guess before you read on. Wrong guesses are fine. After the task, fill in the
answers and fix your guesses.

1. Both stores in `before.go` compile as a `Store`. What does the compiler check, and what can it not check?
   - Guess:
   - Answer:
2. Write down three promises any key-value store should keep.
   - Guess:
   - Answer:
3. How can you catch a lying store without writing a different test for each store?
   - Guess:
   - Answer:

## The idea

**Problem.** Two things have the same shape, so the compiler accepts both,
but one breaks a promise the rest of the code relies on. The mistake hides
until something goes wrong far away.

**Fix.** Write the promises down as a checklist in code, and run the same
checklist against every implementation. A broken promise shows up as a line
in a report.

**Analogy (phone charger).** Imagine a cheap charger whose plug fits the
socket and whose cable fits the phone. It has the right shape, but it quietly
gives the wrong power. You would only find out by running the same check on
every charger you own.

## How the code shows it

| File | What it is |
|---|---|
| `before.go` | The problem: an honest store and a lying one, both accepted by the compiler |
| `after.go` | Your job: write `Contract`, the checklist |
| `lsp_test.go` | `TestBefore_*` pass today. `TestAfter_*` pass when you are done. One test also checks five broken stores, one per promise |

Proof: the honest store reports 0 broken promises, the lying store at least 1.

## Wrong when

There is only one implementation, or nobody depends on the promises.

## In the wild

- Standard library: `go doc testing/fstest.TestFS`. Whose promises does it check, and against which implementations?
- Source: [src/testing/fstest/testfs.go](https://github.com/golang/go/blob/master/src/testing/fstest/testfs.go) (master may differ from your Go 1.26)
- Book: none verified yet.

## My delete-test result

**Try:** ignore `Contract`. Run the lying store through only the normal
Put/Get test (`TestBefore_BothStoresSatisfyTheInterface`). Does anything fail?

What got harder or impossible? One line:
