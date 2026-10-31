# A6 · Errors: sentinel, typed, wrapped

## Questions first

Guess before you read on. Wrong guesses are fine. After the task, fill in the
answers and fix your guesses.

1. In `LoadBefore` the message clearly says "not found". Why does `errors.Is(err, ErrNotFound)` still say no?
   - Guess: Because error was converted to string using %v verb
   - Answer: %v renders the error to text. The result is a plain *errors.errorString with no Unwrap method, so errors.Is has no chain to walk.

     Checked with a sentinel named `ErrSugarHigh`, so its own text cannot be
     misread as "unwrap found nothing":

     ```
     %v  outer="load: sugar is high"   type=*errors.errorString
         Unwrap -> <nil>
         Unwrap == nil ?  true

     %w  outer="load: sugar is high"   type=*fmt.wrapError
         Unwrap -> &errors.errorString{s:"sugar is high"}
         Unwrap == nil ?  false
         Unwrap == ErrSugarHigh ?  true
     ```

     | | `%v` | `%w` |
     |---|---|---|
     | Outer message | identical | identical |
     | Outer type | `*errors.errorString` | `*fmt.wrapError` |
     | `errors.Unwrap` | `nil`, there is no Unwrap method | the original error value |
     | `== ErrSugarHigh` | nothing to compare | `true`, the same value, not a copy |

     The message is identical both ways, which is exactly why reading the text
     proves nothing. `errors.Is` walks the Unwrap chain, and `%v` leaves none.
2. When would you return a fixed sentinel error, and when a typed error with fields?
   - Guess: Typed error will returned when caller want more info from error and sentinal error when its fixed.
   - Answer: Sentinel when the condition is the whole message and the caller only needs to recognise it. Typed when the caller needs a value out of it — a field name, a retry-after, an HTTP status
3. Is wrapping with `%w` always the right choice? What does it commit you to?
   - Guess: No idea
   - Answer: No. %w makes the inner error part of your API: callers may write errors.Is/errors.As against it, so you cannot swap or drop it later without breaking them. Use %v when the inner error is an implementation detail you may change

## The idea

**Problem.** An error passes through several layers, and each layer wants to add
where it happened. If a layer rewrites the error as text, the original is gone.
The caller can only search the message, which breaks the day someone rewords it.

**Fix.** Use three kinds of error on purpose. A **sentinel** is a fixed value for
a known condition. A **typed** error carries data the caller may need. A
**wrapped** error adds context on top and keeps the original inside it, using
`%w`. `errors.Is` and `errors.As` then look through the layers.

**Analogy (blood report).** Your family doctor reads your blood test, writes
"sugar is high, see a specialist" on a fresh sheet, and that sheet is all you
carry. The specialist asks how high, and the sheet cannot say: the words are
right, but the report is gone. Had the doctor stapled the note on top of the
report, the specialist would read it and still turn back to the number. Every
doctor after can staple one more sheet, and the report stays at the bottom.

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

## FAQ

**Why does `Get` return bare `ErrNotFound`? Saying which ID was not found is
more useful.**

You can add the ID, but only by wrapping. New text kills the identity:

```
fmt.Errorf("order %q not found", id)   -> errors.Is = false
fmt.Errorf("get %q: %w", id, ErrNotFound) -> errors.Is = true
```

Then ask who needs the ID:

| Who needs it | Use |
|---|---|
| A human reading the log | Sentinel, with the ID in the wrapped text. The caller passed the ID in, so it is not news to the program |
| The program, as data | A typed error such as `NotFoundError{ID}`, read with `errors.As` |

The ID is not lost here anyway. `Load` adds it, so the final message is
`handle: load "nope": not found`.

**`Load` adds the ID. If `Handle` lives in another package and cannot see that,
it may add it again. What stops the same fact appearing twice?**

Nothing in the code. It does happen:

```
handle "nope": load "nope": not found   (errors.Is still true)
```

It is only noise, not a bug. The fix is a convention, because a caller cannot
see inside its callee: **each layer adds only what the layer below did not
know.** The function that received an argument has already seen it, so it owns
it. The standard library works this way. `os.Open` puts the path in its own
error (`*os.PathError`, field `Path`), so callers of `os.Open` do not add the
path again.

| Layer | Adds |
|---|---|
| The one that received the ID | the ID |
| Each layer above | only what it alone knows: the request, the user, the route |

This task uses the other convention: `Get` does not wrap, so `Load` adds the ID.
Both work if a codebase sticks to one. Callee-owns-it is a recommendation, not a
rule, chosen because it is what the standard library does.

The worse duplication is logging. If every layer logs and returns, one failure
shows up three times. *100 Go Mistakes* #52 "Handling an error twice"
([summary](https://100go.co/#handling-an-error-twice-52)): logging an error is
handling it. Wrap on the way up, log or respond once at the top.

**How does the top layer know to use `errors.As` or `errors.Is`? How does it
know an empty ID is a validation problem?**

It never works that out. `Get` decided it at the bottom, where the knowledge is,
and recorded the decision as the error's kind. The top only asks which kind it
got, and one `switch` covers both questions:

```go
var ve *ValidationError
switch {
case err == nil:
	return 200
case errors.As(err, &ve): // this type? ve.Field is now usable
	return 400
case errors.Is(err, ErrNotFound): // this exact value?
	return 404
default: // unknown: do not show the raw text to users
	return 500
}
```

Run against the three cases: `"42"` gives 200, `"nope"` gives 404, `""` gives
400 with field `id`.

The top layer knows which kinds exist because `ErrNotFound` and
`ValidationError` are exported. They are part of the lower layers' API. Case
order only matters when one error can match both, for example a typed error that
wraps the sentinel. Put the more specific case first.

## In the wild

- Standard library: `go doc errors.Is`. What does "tree" mean in its description, and when does an error have more than one child?
- Source: [src/errors/wrap.go](https://github.com/golang/go/blob/master/src/errors/wrap.go) (master may differ from your Go 1.26)
- Book: *100 Go Mistakes*, #49 "Ignoring when to wrap an error" ([summary](https://100go.co/#ignoring-when-to-wrap-an-error-49)), #50 "Comparing an error type inaccurately" ([summary](https://100go.co/#comparing-an-error-type-inaccurately-50)), #51 "Comparing an error value inaccurately" ([summary](https://100go.co/#comparing-an-error-value-inaccurately-51)).

## My delete-test result

**Try:** change `Load` to wrap with `%v` instead of `%w`, and run
`TestAfter_SentinelSurvivesTwoWraps`. What does it say?

What got harder or impossible? One line:
