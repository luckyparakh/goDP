# HOWTO: how to do any task

The same six steps for every pattern. Replace `<id>` with the folder name, for
example `a1-srp`.

| # | Do this | Where |
|---|---|---|
| 1 | Answer the 3 questions with a **guess**. Do not look anything up. Wrong is fine | `patterns/<id>/THEORY.md` |
| 2 | Read the rest of the theory. About 3 minutes | same file |
| 3 | Read the PROBLEM comment, do TRY IT, then run `go test ./patterns/<id>/` | `before.go` |
| 4 | Write `after.go`, or the files its GOAL names, until every test passes | `after.go` |
| 5 | Open the solution and compare. Read HOW, WHY and WHEN NOT | `solutions/<id>/after.go` |
| 6 | Fill in the **Answer** lines and fix your guesses. Do the delete test (one line). Add a row to the log | `THEORY.md`, then `PROGRESS.md` |

## What you should see at step 3

`TestBefore_*` pass and `TestAfter_*` fail. That is the correct starting line.
A red test is a good start. If the tests do not even compile, that is a bug in
the task, so tell Claude.

## The delete test (step 6)

It shows what your fix actually bought. You put the "before" shape back and see
what gets harder.

1. Commit your work first, so you can undo (`git add -A` then `git commit`).
2. Do the **Try** line in the pattern's `THEORY.md`. It says how to put the old
   shape back and what to attempt.
3. Write one line: the first thing that got harder or impossible. If nothing did,
   write "nothing" and tell Claude. The fix may have been decoration.
4. Undo with `git checkout -- .`. Only do this after step 1, because it throws
   away uncommitted changes in tracked files.

The `TestBefore_*` tests are your evidence. They still pass on the old shape,
and each one shows something awkward.

## Rules

- **Do not open `solutions/` before you have tried step 4.** That is the point.
- Stuck? Ask Claude: "hint for `<id>`". You get a question back, not an answer.
- Finished? Tell Claude: "I attempted `<id>`, review it".

## Commands

| Command | What it does |
|---|---|
| `go test ./patterns/<id>/` | Run one pattern's tests. Use this loop |
| `go test ./patterns/<id>/...` | The same, for a task with subfolders (B1, B2, B5) |
| `go test -race ./patterns/<id>/` | Also required for every concurrency task |
| `gofmt -l .` | Lists badly formatted files. Should print nothing |
| `go vet ./...` | Should print nothing |
| `go doc -src <symbol>` | Read the real standard library, for "In the wild" |

`go test ./...` from the repo root fails while any pattern is unfinished. That
is expected, so run the single-pattern command.

## Done means

- [ ] Every test in `patterns/<id>/` passes, and `gofmt -l .` and `go vet ./...` print nothing
- [ ] The three **Answer** lines are filled in and your wrong guesses are corrected
- [ ] The delete-test line is filled in
- [ ] A row is added to `PROGRESS.md`, including what you got wrong
