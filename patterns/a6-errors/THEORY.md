# A6 · Errors: sentinel, typed, wrapped

## Questions first

Guess before you read on. Wrong guesses are fine. After the task, fill in the
answers and fix your guesses.

1. In `LoadBefore` the message clearly says "not found". Why does `errors.Is(err, ErrNotFound)` still say no?
   - Guess:
   - Answer:
2. When would you return a fixed sentinel error, and when a typed error with fields?
   - Guess:
   - Answer:
3. Is wrapping with `%w` always the right choice? What does it commit you to?
   - Guess:
   - Answer:

## The idea

**Problem.** An error passes through several layers, and each layer wants to add
where it happened. If a layer rewrites the error as text, the original is gone.
The caller can only search the message, which breaks the day someone rewords it.

**Fix.** Use three kinds of error on purpose. A **sentinel** is a fixed value for
a known condition. A **typed** error carries data the caller may need. A
**wrapped** error adds context on top and keeps the original inside it, using
`%w`. `errors.Is` and `errors.As` then look through the layers.

**Analogy (home kitchen).** You find there is no rice and tell the next person.
They stick their own note on top of yours instead of rewriting it, and the next
person does the same. Anyone can still peel the notes back to your original.
If each person rewrote the message in their own words, the one fact that
mattered would be lost.

## How the code shows it

| File | What it is |
|---|---|
| `before.go` | The problem. `LoadBefore` wraps with `%v` and loses both errors |
| `after.go` | Your job: `Get`, `Load` and `Handle`, three layers that keep the errors |
| `errs_test.go` | `TestBefore_*` pass today. `TestAfter_*` pass when you are done |

Proof: `errors.Is` and `errors.As` still work after two wraps.

## Wrong when

A layer has nothing useful to add, or you do not want callers to depend on the
inner error. Wrapping with `%w` makes it part of your API, so hide it on purpose
with `%v` if it is an implementation detail.

## In the wild

- Standard library: `go doc errors.Is`. What does "tree" mean in its description, and when does an error have more than one child?
- Source: [src/errors/wrap.go](https://github.com/golang/go/blob/master/src/errors/wrap.go) (master may differ from your Go 1.26)
- Book: *100 Go Mistakes*, #49 "Ignoring when to wrap an error" ([summary](https://100go.co/#ignoring-when-to-wrap-an-error-49)), #50 "Comparing an error type inaccurately" ([summary](https://100go.co/#comparing-an-error-type-inaccurately-50)), #51 "Comparing an error value inaccurately" ([summary](https://100go.co/#comparing-an-error-value-inaccurately-51)).

## My delete-test result

**Try:** change `Load` to wrap with `%v` instead of `%w`, and run
`TestAfter_SentinelSurvivesTwoWraps`. What does it say?

What got harder or impossible? One line:
