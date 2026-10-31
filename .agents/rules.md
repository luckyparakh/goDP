# Rules for this repo

The learner is a backend services engineer on Go 1.26. They want plain English,
short sentences and tables over prose. They do not want ceremony. Background on
how the plan was reached is in `design_pattern_handoff.md`. It is long and
partly out of date, so read it only when you need the reason behind a decision.

## Shape of one pattern

```
patterns/<id>/THEORY.md   questions first, then the idea, analogy, "wrong when", optional FAQ, "in the wild", delete test
patterns/<id>/before.go   the problem, in comments, plus a TRY IT
patterns/<id>/after.go    a stub that compiles, with the GOAL in comments
patterns/<id>/<name>_test.go
solutions/<id>/after.go   sealed. First line is //go:build ignore, then a blank line. HOW / WHY / WHEN NOT in comments
```

A task about packages (B1, B2, B5) also has real subpackages. `before.go` and
`after.go` at the root then hold only the PROBLEM and GOAL comments, and the
code lives in `before/<pkg>/` (do not edit) and `after/<pkg>/` (the learner's
work). The tests sit at the root and import both. The solution mirrors
`after/` under `solutions/<id>/after/`, every file sealed the same way, and
`solutions/<id>/after.go` still holds HOW / WHY / WHEN NOT. Tests must compile
before and after the learner's change, so they never name a type that moves
between packages.

The task list is `BLUEPRINT.md`. The learner's steps are `HOWTO.md`. The log is `PROGRESS.md`.

## Rules

- **One pattern at a time.** Generate one folder, wait for the learner's attempt, review it.
- **Starting state:** every `TestBefore_*` passes and every `TestAfter_*` fails. Stubs must compile. A compile error is a bug.
- **Sealed solutions.** Do not reveal or hint at `solutions/` until the learner says they attempted. If asked for a hint, ask a question.
- **Do not write the learner's lines.** The `Guess`, `Answer`, delete-test and `PROGRESS.md` entries are theirs.
- **Run everything you generate:** `gofmt -l .`, `go vet ./...`, `go test`, and `.agents/scripts/check-pattern.sh <id> start`. Paste real output, never expected output. Concurrency tasks also need `go test -race`.
- **THEORY.md has 3 pointed questions at the top**, each with `Guess:` and `Answer:` lines. The learner guesses before reading and answers after the task.
- **Analogies:** made-up is fine, and each task picks whatever everyday situation fits it best. No central story, and no reusing a domain for its own sake. See the analogy rule in `skills/new-pattern.md` for what makes one good: pain first, one situation, three or four short sentences, no unfamiliar words, India-based where it fits, something people really do, and the bad version must really lose something in the story. An analogy is an illustration, so never claim a fact about the real world through it.
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
