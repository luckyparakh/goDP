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
6. **Analogy:** see the analogy rule below.
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

1. Create the five files from the shape in `rules.md`, or the multi-package
   shape described there when the task is about packages. Write every sealed
   solution file with `//go:build ignore` on the first line and a blank second
   line.
2. Run `gofmt -l .`, `go vet ./...` and
   `bash .agents/scripts/check-pattern.sh <id> start` until it says PASS.
3. For a concurrency task also run `go test -race -count=50 -shuffle=on ./patterns/<id>/`
   and report any flakiness.
4. Hand over with five lines or fewer: what was built and the real check
   output. Add the new row to `BLUEPRINT.md` only if it is missing.
5. Do not generate the next pattern.

## Analogy rule

Pick whatever everyday situation fits the task best. There is no central story
and no fixed set of domains.

- **Pain in the first sentence.** Never describe the fix before the reader has
  felt the problem. This is the mistake to check for first, and the one that
  spoiled the first versions of A4, A5 and A9.
- **One situation, not a list.** If it takes four example items to work, it is
  too complicated. The reader should picture one thing.
- **Three or four short sentences.** Longer means the idea is not clear yet.
- **Everyday things the learner actually touches.** The learner is in India, so
  prefer what is ordinary there: a gate register, a borewell or water tanker, an
  auto, a tiffin, an extension board, a chit sent with someone. Avoid anything
  regional to somewhere else. "Boiler" and "bin day" were real rejections.
- **No word should need explaining.** If you have to gloss a term, pick a
  different situation.
- **People already do this.** Every action in the story is something people
  really do in that situation. If you had to invent a behaviour so it looks like
  the code, pick another situation. A6's first version had people relaying "no
  rice" by sticking notes on top of notes. Nobody does that.
- **Carry the task's sharpest point, not just the shape.** Read the proof the
  test gives and make sure the analogy shows that same thing. A4's point is that
  the big implementation still satisfies the small interface, so the analogy has
  to say that anyone who can fill the register can write their name.
- **The bad version must really hurt in the story.** Tell the "before" version
  inside the story and check that someone actually loses something they need.
  If a person would shrug, the analogy teaches the opposite. A6's rice note
  failed this: reword "no rice" and the fact survives, so rewording looks
  harmless, when the whole point is that the text survives and the error does
  not.
- **Map it before you write it.** In the spec, give a small table: each key code
  piece (the bad move, the fix, what the test proves) against its part of the
  story. If a row is empty or forced, change the situation, not the wording.
- **Label it and keep it hypothetical.** It illustrates, and never asserts a
  fact about the real world.

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
