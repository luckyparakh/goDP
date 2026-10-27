# Skill: new-pattern

Use when the learner asks for the next pattern or names an id from
`BLUEPRINT.md`. It has two phases, and you **stop between them**.

Read first: [../rules.md](../rules.md) and the row in `BLUEPRINT.md`. For shape,
look at an untouched finished pattern such as `patterns/a2-ocp/`, and at
`solutions/a1-srp/after.go` for the sealed-solution style.

## Phase 1: spec (write nothing to disk)

Answer in the chat:

1. **Problem** in two sentences, and **what varies**.
2. **before.go:** the types, the pain, and a TRY IT the learner can do without
   editing anything, if possible.
3. **after.go stub:** the exact signatures. It must compile.
4. **Tests:** names `TestBefore_*` and `TestAfter_*`, and what each checks. Say
   which fail at the start.
5. **Analogy:** 3 to 4 plain sentences, within the domain rules.
6. **3 questions** for the top of `THEORY.md`, each with the answer you expect.
   The expected answers stay in the spec and never go into the files.
7. **Wrong when**, and the **delete-test Try** line.
8. **In the wild:** symbol, source path and book entry. Verify each one by
   running `go doc`, checking the file under GOROOT, and checking 100go.co.
9. **Risks and unknowns.** Anything you have not verified.

Then stop and ask for approval. Do not create files.

## Phase 2: build (after approval)

1. Create the five files from the shape in `rules.md`. Write the sealed solution
   with `//go:build ignore` on the first line and a blank second line.
2. Run `gofmt -l .`, `go vet ./...` and
   `bash .agents/scripts/check-pattern.sh <id> start` until it says PASS.
3. For a concurrency task also run `go test -race -count=50 -shuffle=on ./patterns/<id>/`
   and report any flakiness.
4. Hand over with five lines or fewer: what was built and the real check
   output. Add the new row to `BLUEPRINT.md` only if it is missing.
5. Do not generate the next pattern.
