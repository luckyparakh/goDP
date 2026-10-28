# Skill: new-pattern

Use when the learner asks for the next pattern or names an id from
`BLUEPRINT.md`. It has two phases, and you **stop between them**.

Read first: [../rules.md](../rules.md) and the row in `BLUEPRINT.md`. For shape,
look at an untouched finished pattern such as `patterns/a2-ocp/`, and at
`solutions/a1-srp/after.go` for the sealed-solution style.

## Phase 1: spec (write nothing to disk)

Answer in the chat:

1. **The principle in general terms**, for `## What the principle says`. Name
   it, say who stated it and roughly when, give the classic formulation, and say
   what changes in Go. See the general-layer rule below. Skip only when the task
   teaches a Go mechanism with no named principle behind it.
2. **Problem** in two sentences, and **what varies**.
3. **before.go:** the types, the pain, and a TRY IT the learner can do without
   editing anything, if possible.
4. **after.go stub:** the exact signatures. It must compile.
5. **Tests:** names `TestBefore_*` and `TestAfter_*`, and what each checks. Say
   which fail at the start.
6. **Analogy:** 3 to 4 plain sentences, within the domain rules.
7. **3 questions** for the top of `THEORY.md`, each with the answer you expect.
   The expected answers stay in the spec and never go into the files.
8. **Wrong when**, and the **delete-test Try** line.
9. **FAQ (optional):** one or two questions a learner is likely to ask about a
   known trade-off, such as "is a global variable bad?". Skip it if there is no
   real question. See the FAQ rule below.
10. **In the wild:** symbol, source path and book entry. Verify each one by
    running `go doc`, checking the file under GOROOT, and checking 100go.co.
11. **Risks and unknowns.** Anything you have not verified.

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

## General-layer rule

A learner who searches the web for the principle's name must be able to map what
they find back onto the task. If `THEORY.md` only describes this code's problem,
they cannot, and the task reads as a trick with no concept behind it. That
happened with `a3-lsp`, where the principle is about substitution but the task
asks for a test, and nothing joined the two.

`## What the principle says` goes first, before `## The idea`, and holds:

- the principle's name in full, who stated it and roughly when
- the classic formulation, in one or two sentences
- the textbook violation, if there is a well-known one
- **what changes in Go**, especially when the classic form assumes inheritance
- the usable question, the one you ask real code
- a link to `../SOLID.md` for the SOLID five

Where the task's mechanism is not obviously the principle, spell out the chain
from one to the other step by step. **Use bullets, never a numbered list**, for
the reason in the FAQ rule below.

Keep it to one short screen. It explains the principle, not the code: the code
belongs in `## The idea` and the tables after it.

## FAQ rule

A `## FAQ` section is optional. It also gets added later when a learner asks a
question worth keeping, as happened with the global variable in `a2-ocp`.

- Place it between `## Wrong when` and `## In the wild`.
- Each entry is a bold question followed by a short answer.
- **Use bullets, never numbered lists.** `check-pattern.sh` counts lines that
  start with `1. ` to confirm there are exactly 3 questions, so a numbered list
  anywhere in `THEORY.md` would break that check.
- Only state what you checked, with `go doc`, the source, or a run. Say so when
  it is an opinion, and give both sides of a trade-off.
