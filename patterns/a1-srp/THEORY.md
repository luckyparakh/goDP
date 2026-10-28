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

## What the principle says

**Single Responsibility Principle.** Usually stated as: a module should have one
reason to change. The catch is that nobody agrees what "a responsibility" is, so
the rule cannot be checked as written. The usable version is a question: *how
many different requirements would make me edit this?* Two different requirements
means two reasons to change, and that is where the split goes.

Nothing about this is Go-specific. The principle was written about classes, and
Go has no classes, so read "module" as a function, a type or a package.

Cross-principle view: [SOLID.md](../SOLID.md).

## The idea

**Problem.** One function does several different jobs. Each job has its own
reason to change, so every change edits the same function and every check
runs all the jobs. You also cannot ask for one result without going through
the others, including their failures.

**Fix.** Split by *reason to change*. One piece adds up, one builds the text,
one saves. Each can be changed and tested alone.

**Analogy (the shop bill).** The kirana shop writes your bill, files his copy
and updates your khata, all in one book, in one go. Ask him just what you owe
and he has to work through all three. Change the tax on rice and he rewrites the
page that also holds your khata. Three separate books, and each question takes
one.

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
