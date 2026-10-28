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

**Analogy (the cancelled order).** You order a cake, so the shop starts baking
and sends someone out for candles. You cancel ten minutes later, but nobody
tells the kitchen, so the cake is finished and the candles are bought for a
party that is off. One cancel slip that every person checks before their next
step would have stopped all of it at once.

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
