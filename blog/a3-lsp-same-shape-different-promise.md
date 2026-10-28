# Same shape, different promise

**TL;DR** — Two key-value stores implement the same Go interface. The compiler
accepts both, but one returns "no error" for a key it does not have, and every
caller that relies on `ErrNotFound` quietly misbehaves. The fix is a contract
test: one checklist of promises, run against every implementation. The honest
store reports 0 broken promises, the liar reports at least 1.

## The everyday version

Your phone charger dies, so you swap in a replacement. Same plug, same cable,
and the phone charges, so as far as you can tell it is the same thing. Imagine
this one quietly delivers the wrong power: the phone cannot tell the difference,
and the damage only shows up months later.

## How you'd fix it in real life

Stop trusting the plug shape. Put every charger through the same set of checks
before you rely on it, and keep the checks written down so you can run them on
the next one too.

> Where this analogy breaks: you can measure a charger directly. For code, the
> "checks" are something you have to write yourself, and they only cover the
> promises you remembered to write down.

## The same pain, in code

Here is the interface. Note that the sentinel is part of the contract, not an
implementation detail:

```go
// Store is a string key-value store.
//
// The contract. Every implementation must keep all five promises:
//  1. Get on a key that was never stored returns ErrNotFound.
//  2. Put on a healthy store succeeds.
//  3. After Put(k, v), Get(k) returns v.
//  4. Put on an existing key succeeds (it overwrites).
//  5. After overwriting, Get returns the newest value.
type Store interface {
	Get(key string) (string, error)
	Put(key, value string) error
}
```

Now two implementations. The second one breaks promise 1 — a missing key comes
back as `("", nil)`:

```go
func (s *goodStore) Get(key string) (string, error) {
	v, ok := s.m[key]
	if !ok {
		return "", ErrNotFound
	}
	return v, nil
}

func (s *lyingStore) Get(key string) (string, error) {
	return s.m[key], nil
}
```

And this compiles, which is the whole problem:

```go
var (
	_ Store = (*goodStore)(nil)
	_ Store = (*lyingStore)(nil)
)
```

The compiler checks method names and signatures. That is all it can do. It has
no way to know that promise 1 exists, so it cannot tell these two apart. Pass
the liar to code that does `errors.Is(err, ErrNotFound)` and the branch simply
never runs.

## The fix

Write the promises as runnable checks, in a plain function that takes a factory
and returns one line per broken promise:

```go
func Contract(newStore func() Store) []string {
	var errs []string
	store1, store2 := newStore(), newStore()

	// Promise 2
	if err := store1.Put("k", "v"); err != nil {
		errs = append(errs, fmt.Sprintf("2. put error %s", err))
	}

	// Promise 1
	_, getNoKeyErr := newStore().Get("no")
	if !errors.Is(getNoKeyErr, ErrNotFound) {
		errs = append(errs, "1. no key found error")
	}

	// ... promises 3, 4 and 5 trimmed
	return errs
}
```

Two details matter more than they look.

**It returns strings instead of calling `t.Error`.** That is what lets a test
assert the *liar* gets caught, without shipping a deliberately failing test.

**Promise 1 uses `errors.Is`, not `err != nil`.** Checking only that some error
came back is not enough, and that distinction is the subject of the mistake
section below.

## Proof

Starting from the stub, `go test` reports:

```
--- FAIL: TestAfter_GoodStoreKeepsEveryPromise (0.00s)
    lsp_test.go:46: Contract never called newStore
--- FAIL: TestAfter_LyingStoreIsCaught (0.00s)
    lsp_test.go:52: violations: []
    lsp_test.go:54: lying store passed the contract
--- FAIL: TestAfter_CatchesEachBrokenPromise/promise_1_(missing_key) (0.00s)
        lsp_test.go:109: broken store passed the contract
... five more broken stores, one per promise
FAIL
```

After writing `Contract`, the same command:

```
--- PASS: TestBefore_BothStoresSatisfyTheInterface (0.00s)
--- PASS: TestBefore_LyingStoreHidesMissingKeys (0.00s)
--- PASS: TestAfter_GoodStoreKeepsEveryPromise (0.00s)
--- PASS: TestAfter_LyingStoreIsCaught (0.00s)
--- PASS: TestAfter_CatchesEachBrokenPromise (0.00s)
    --- PASS: TestAfter_CatchesEachBrokenPromise/promise_1_(missing_key) (0.00s)
    --- PASS: TestAfter_CatchesEachBrokenPromise/promise_1_(wrong_error) (0.00s)
    --- PASS: TestAfter_CatchesEachBrokenPromise/promise_2_(put_fails) (0.00s)
    --- PASS: TestAfter_CatchesEachBrokenPromise/promise_3_(get_after_put) (0.00s)
    --- PASS: TestAfter_CatchesEachBrokenPromise/promise_4_(overwrite_succeeds) (0.00s)
    --- PASS: TestAfter_CatchesEachBrokenPromise/promise_5_(newest_value_wins) (0.00s)
ok  	patterns/a3-lsp
```

Both runs on Go 1.26, Windows. The failing block is trimmed where marked; the
passing block is complete.

## What I got wrong

> TODO-yours. From PROGRESS.md: *"Missed to check sentinal error, wrote not a
> clean code, like error message can be more informative"*.
>
> Write 3 or 4 sentences. The interesting part is *why* `err != nil` felt like
> enough: a store that returns `errors.New("record absent")` for a missing key
> does return an error, passes that check, and still breaks every caller doing
> `errors.Is(err, ErrNotFound)`. Worth saying that the test suite did not catch
> this either — the mutant for it had to be added afterwards.

## The delete test

> TODO-yours. Ignore `Contract` and run the lying store through only the normal
> Put/Get test. Then write one line: what can you no longer catch?

## Do and don't

| Do | Don't |
|---|---|
| Write the promises in the interface's doc comment | Leave the contract in someone's head |
| Check the specific error with `errors.Is` | Settle for `err != nil` |
| Run the same checks against the fakes too | Trust a fake you never verified |
| Ask the factory for a fresh value per promise group | Let one broken promise cascade into four reports |

## When not to use this

If there is only one implementation and nobody depends on the promises, a normal
test is enough and a contract suite is machinery you do not need. It also only
catches promises somebody wrote down: a store can still break a sixth promise
nobody listed, and the checklist will pass it.

## Further reading

- `testing/fstest.TestFS` is the standard library's own contract test. It checks
  an `fs.FS` implementation and, like the one above, returns an error rather
  than taking a `*testing.T`. Its docs show callers writing
  `if err := fstest.TestFS(myFS, "..."); err != nil { t.Fatal(err) }`.
- Source: [src/testing/fstest/testfs.go](https://github.com/golang/go/blob/master/src/testing/fstest/testfs.go)
- The code, tests and task card for this post:
  [`patterns/a3-lsp`](../patterns/a3-lsp)
