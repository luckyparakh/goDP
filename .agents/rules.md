# Rules for this repo

The learner is a backend services engineer on Go 1.26. They want plain English,
short sentences and tables over prose. They do not want ceremony. Background on
how the plan was reached is in `design_pattern_handoff.md`. It is long and
partly out of date, so read it only when you need the reason behind a decision.

## Shape of one pattern

```
patterns/<id>/THEORY.md   questions first, then the idea, analogy, "wrong when", "in the wild", delete test
patterns/<id>/before.go   the problem, in comments, plus a TRY IT
patterns/<id>/after.go    a stub that compiles, with the GOAL in comments
patterns/<id>/<name>_test.go
solutions/<id>/after.go   sealed. First line is //go:build ignore, then a blank line. HOW / WHY / WHEN NOT in comments
```

The task list is `BLUEPRINT.md`. The learner's steps are `HOWTO.md`. The log is `PROGRESS.md`.

## Rules

- **One pattern at a time.** Generate one folder, wait for the learner's attempt, review it.
- **Starting state:** every `TestBefore_*` passes and every `TestAfter_*` fails. Stubs must compile. A compile error is a bug.
- **Sealed solutions.** Do not reveal or hint at `solutions/` until the learner says they attempted. If asked for a hint, ask a question.
- **Do not write the learner's lines.** The `Guess`, `Answer`, delete-test and `PROGRESS.md` entries are theirs.
- **Run everything you generate:** `gofmt -l .`, `go vet ./...`, `go test`, and `.agents/scripts/check-pattern.sh <id> start`. Paste real output, never expected output. Concurrency tasks also need `go test -race`.
- **THEORY.md has 3 pointed questions at the top**, each with `Guess:` and `Answer:` lines. The learner guesses before reading and answers after the task.
- **Analogies:** made-up is fine. Use 3 to 4 plain sentences, two or three domains at most (a home kitchen, plugs and chargers, and one more only if needed such as train tickets), and never a story that continues across patterns. An analogy is an illustration, so never claim a fact about the real world through it.
- **"In the wild" has three lines:** a standard-library symbol with a question, its source path in `github.com/golang/go` (check the file exists under `$(go env GOROOT)`), and a book entry. Cite *100 Go Mistakes* by mistake number and title, linked as `https://100go.co/#<title-hyphenated>-<number>`. Never invent a chapter number. If no entry is verified, write "none verified yet".
- **Facts about code or the standard library:** check with `go doc` or by running it, or say it is unverified. Never state a result you did not run.
- **Every line of code is original.** Never copy Refactoring Guru code.
- **Public repo.** Never put company or customer code, names or secrets in it, including in `PROGRESS.md`.
- **Push back on ceremony.** An interface with one implementation and one caller should be deleted.
- **Out of scope:** system design, DDIA, CQRS, strategic DDD.

## Check a solution without a second tree

`.agents/scripts/check-pattern.sh <id> [start|done]` copies the pattern into a
throwaway module, swaps in the solution files with their first two lines
removed, and runs `go vet` and `go test -race`. It changes nothing in the repo.
