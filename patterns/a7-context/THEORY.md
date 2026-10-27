# A7 · Context: pass the stop signal down

## Questions first

Guess before you read on. Wrong guesses are fine. After the task, fill in the
answers and fix your guesses.

1. `FetchNoCtx` cannot be stopped. What keeps running after the caller gives up, and who pays for it?
   - Guess:
   - Answer:
2. Where does a `context.Context` go in a function's parameters, and should you keep one in a struct?
   - Guess:
   - Answer:
3. A timeout is set at the top. How does the bottom level find out it is time to stop?
   - Guess:
   - Answer:

## The idea

**Problem.** A caller sets a time limit or changes its mind, but the code it
called has no way to hear that. The work carries on for someone who has left,
and every connection, buffer and goroutine it holds stays in use.

**Fix.** Pass a `context.Context` as the first parameter through every level of
the call chain. The top creates it with a timeout or a cancel. The level that
actually waits listens for it ending, and stops. The error it returns says why.

**Analogy (home kitchen).** You start dinner and send a friend shopping. If you
decide to order takeaway, you stick a note on the fridge and everyone checks it
before starting the next step. An alarm that writes "stop" after 30 minutes
does the same job without you. Without the note, everyone finishes the job for
a dinner nobody wants.

## How the code shows it

| File | What it is |
|---|---|
| `before.go` | The problem. `FetchNoCtx` runs all 5 steps no matter what |
| `after.go` | Your job: `Handle`, `Service` and `Repo`, passing `ctx` down and stopping when it ends |
| `ctxflow_test.go` | `TestBefore_*` pass today. `TestAfter_*` pass when you are done |

Proof: a 30ms timeout makes `ctx.Err()` equal `DeadlineExceeded`, and the work
stops before the fifth step.

## Wrong when

The function does no waiting and cannot fail, so a context only adds a
parameter. Never use `context.Background()` in the middle of a chain, because it
cuts the stop signal off.

## In the wild

- Standard library: `go doc context.WithTimeout`. Why must you call the `cancel` function it returns, even if the timeout fires on its own?
- Source: [src/context/context.go](https://github.com/golang/go/blob/master/src/context/context.go) (master may differ from your Go 1.26)
- Book: *100 Go Mistakes*, #60 "Misunderstanding Go contexts" ([summary](https://100go.co/#misunderstanding-go-contexts-60)) and #61 "Propagating an inappropriate context" ([summary](https://100go.co/#propagating-an-inappropriate-context-61)).

## My delete-test result

**Try:** change `Service` so it calls `Repo(context.Background())` instead of
`Repo(ctx)`, and run `TestAfter_TimeoutReachesTheBottom`.

What got harder or impossible? One line:
