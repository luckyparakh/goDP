# A3 · LSP: same shape is not the same promise

## Questions first

Guess before you read on. Wrong guesses are fine. After the task, fill in the
answers and fix your guesses.

1. Both stores in `before.go` compile as a `Store`. What does the compiler check, and what can it not check?
   - Guess: It can check signature.
   - Answer: It can check only sigranture but not the behaviour of implementation
2. Write down three promises any key-value store should keep.
   - Guess: Put should work with no error, Get should give error in case of no key and Update should work properly
   - Answer: Same as above
3. How can you catch a lying store without writing a different test for each store?
   - Guess: Write contarct as part of code
   - Answer: Write contract in test like style to check each case

## What the principle says

**Liskov Substitution Principle**, from Barbara Liskov (1987), formalised with
Jeannette Wing (1994). If code works with a type `T`, it must keep working when
handed any subtype `S` of `T`. Not merely compile. Keep working.

The textbook violation is Rectangle and Square. `Square` matches every method
signature of `Rectangle`, but setting its width also changes its height, so code
that sets the width and then reads the height breaks. The signatures were never
the problem.

### Why this becomes a test in Go

That formulation assumes inheritance, and Go has none. The chain runs like this:

- In Go, "S can stand in for T" means a type implements an interface, so it can
  be used wherever that interface is expected.
- The compiler checks method names and signatures. Nothing else.
- It cannot check behaviour, so it cannot tell an honest implementation from a
  lying one.
- Therefore every LSP violation that survives compilation is a behavioural one.
- So the only way to check substitutability is to run the same behavioural
  checks against every implementation.
- Those checks are a **contract test**. `Contract` in this task is one.

One thing follows from that chain: the promises have to be written down
somewhere, or "substitutable" has no meaning. Here they are the five promises in
the doc comment on `Store`. Writing them down is part of the discipline, not
paperwork.

The usable question is: *does every implementation keep the same promises?*

Cross-principle view: [SOLID.md](../SOLID.md).

## The idea

**Problem.** Two things have the same shape, so the compiler accepts both,
but one breaks a promise the rest of the code relies on. The mistake hides
until something goes wrong far away.

**Fix.** Write the promises down as a checklist in code, and run the same
checklist against every implementation. A broken promise shows up as a line
in a report.

**Analogy (phone charger).** Your charger dies, so you swap in a replacement.
Same plug, same cable, and the phone charges, so as far as you can tell it is
the same thing. Imagine this one quietly delivers the wrong power: the phone
cannot tell the difference, and the damage only shows up months later. The only
way to know before you trust it is to put every charger through the same checks.

## How the code shows it

| File | What it is |
|---|---|
| `before.go` | The problem: an honest store and a lying one, both accepted by the compiler |
| `after.go` | Your job: write `Contract`, the checklist |
| `lsp_test.go` | `TestBefore_*` pass today. `TestAfter_*` pass when you are done. One test also checks five broken stores, one per promise |

Proof: the honest store reports 0 broken promises, the lying store at least 1.

## Often confused with

Searching for "LSP in Go" turns up two things that are real but are different
lessons here. Each has a kernel of truth, which is exactly why they get mixed up.

| What you will read | What it actually is | Where it lives here |
|---|---|---|
| "Keep interfaces to one or two methods" | **ISP**, a different principle. It helps substitution as a side effect, because fewer promises means fewer ways to break them | [A4](../a4-isp/THEORY.md) |
| "Embedding abuse" | Embedding is not subtyping. People expect a drop-in substitute and do not get one, and the compiler proves it | [A8](../a8-embedding/THEORY.md) |
| A contract test | The only way to **check** LSP in Go. It is the tool, not the principle | this task |

## Wrong when

There is only one implementation, or nobody depends on the promises.

## FAQ

**Who calls `Contract`, and when?**

Tests call it. It never runs in production. The function lives in a normal
`.go` file rather than a `_test.go` file so that it can be exported and called
from other packages, but its only callers are tests.

- **Who:** the test for *every* implementation of the interface. The real one,
  and every fake you wrote for other tests.
- **When:** in CI, on every change. It matters most on two days: when someone
  adds a new implementation, and when someone changes the contract itself. Add a
  sixth promise and every implementation gets re-checked for free.
- **Why the fakes too:** this is the part people skip. If your in-memory fake
  does not keep the same promises as the real store, your unit tests are green
  against a fake that lies, and the truth only shows up in production. Task F2
  reuses this suite for exactly that.

**What does it look like in real code?**

The contract usually lives in its own small package beside the interface, so
many implementation packages can import it without an import cycle:

```
internal/store/store.go                  the Store interface, ErrNotFound, the promises
internal/store/storetest/contract.go     Contract, importing store
internal/store/postgres/postgres_test.go calls storetest.Contract
internal/store/memory/memory_test.go     calls storetest.Contract  (the fake)
```

Each implementation's test is then about three lines:

```go
func TestPostgresStore_Contract(t *testing.T) {
	storetest.Contract(t, func() store.Store { return postgres.New(newTestDB(t)) })
}
```

Note the `t` in that signature. In production code the contract usually takes
`*testing.T` and reports each promise as a subtest, so a failure names the
promise it broke. This task returns `[]string` instead for one specific reason:
the test suite here has to assert that a *bad* store is caught, and a function
that calls `t.Error` internally could not do that without turning your own suite
red on purpose.

Both shapes are real. The standard library's own contract test,
`fstest.TestFS`, returns an `error` rather than taking a `*testing.T`, and its
docs show the caller doing `if err := fstest.TestFS(...); err != nil { t.Fatal(err) }`.

## In the wild

- Standard library: `go doc testing/fstest.TestFS`. Whose promises does it check, and against which implementations?
- Source: [src/testing/fstest/testfs.go](https://github.com/golang/go/blob/master/src/testing/fstest/testfs.go) (master may differ from your Go 1.26)
- Book: none verified yet.

## My delete-test result

**Try:** ignore `Contract`. Run the lying store through only the normal
Put/Get test (`TestBefore_BothStoresSatisfyTheInterface`). Does anything fail?

What got harder or impossible? One line:
