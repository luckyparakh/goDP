# A1 · SRP: one job per piece

## Questions first

Guess before you read on. Wrong guesses are fine. After the task, fill in the
answers and fix your guesses.

1. Look at `Process` in `before.go`. How many different reasons could force you to edit it?
   - Guess: 3
   - Answer: 3
2. Why does the "before" test need a temp folder just to check one sum?
   - Guess: 
   - Answer: Beacuse sum lives in same func which creates file
3. What is the downside of splitting it into pieces?
   - Guess: More boilerplate 
   - Answer: Now some other func, has to create a pipeline to call total -> render -> save

## The idea

**Problem.** One function does several different jobs. Each job has its own
reason to change, so every change edits the same function and every check
runs all the jobs. You also cannot ask for one result without going through
the others, including their failures.

**Fix.** Split by *reason to change*. One piece adds up, one builds the text,
one saves. Each can be changed and tested alone.

**Analogy (home kitchen).** Imagine one recipe card that also holds your
shopping list and the washing-up rota. To change the shopping list you reprint
the whole card. To check the salt in the recipe you have to do the washing-up.
Three separate cards fix that.

## How the code shows it

| File | What it is |
|---|---|
| `before.go` | The problem. `Process` totals, writes text and saves, all in one |
| `after.go` | Your job: write `Total`, `Render` and `Save` |
| `srp_test.go` | `TestBefore_*` pass today. `TestAfter_*` pass when you are done |

Proof: the sum is tested with no temp folder.

## Wrong when

The pieces always change together and are never used apart.

## In the wild

- Standard library: `go doc log/slog.Handler`. What to check: a logger has two
  jobs, making a log record and deciding how it looks and where it goes. `slog`
  gives the second job to a `Handler`. Run this and compare the two lines. The
  call to log did not change; only the handler did.

  ```go
  text := slog.New(slog.NewTextHandler(os.Stdout, nil))
  json := slog.New(slog.NewJSONHandler(os.Stdout, nil))
  text.Info("order saved", "id", 7) // level=INFO msg="order saved" id=7
  json.Info("order saved", "id", 7) // {"level":"INFO","msg":"order saved","id":7}
  ```

  (Both also print a time field, which I left out. I ran this on Go 1.26.)
- Source: [src/log/slog/handler.go](https://github.com/golang/go/blob/master/src/log/slog/handler.go) (master may differ from your Go 1.26)
- Book: none verified yet.

## My delete-test result

**Try:** call `Process` instead of your three pieces, and try to check one sum
without making a folder.

What got harder or impossible? One line: I can't get sum w/o foler and a bad folder will give error.
