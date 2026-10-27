# Skill: review-attempt

Use when the learner says they attempted a pattern (an id like `a2-ocp`), or
asks for a hint.

Read first: [../rules.md](../rules.md), then `patterns/<id>/THEORY.md`,
`before.go`, the learner's `after.go` and the test file.

Do not open or quote `solutions/<id>/` until the learner has said they attempted.

## Mode A: review

1. Run `go test -count=1 -race ./patterns/<id>/`, `gofmt -l patterns/<id>` and
   `go vet ./patterns/<id>/`. Report the real output.
2. Read `after.go` line by line. Look past the tests:
   - Does it behave the same on other operating systems? A test that passes on
     Windows can hide a bug elsewhere.
   - Is error handling and wrapping consistent with `before.go`?
   - Is anything unclear, or anything the tests do not cover?
   Before you say a line is a bug, prove it by running a tiny program.
3. Read the `THEORY.md` answers. For each of the 3 questions mark the `Answer` as
   Correct, Partly or Missing. If it is not fully right, give a sharper answer
   in one or two plain sentences. Leave the learner's `Guess` as it is.
4. Check the delete-test line. It should name one specific thing that became
   harder. "It throws an error" is not enough.
5. Do not write the learner's lines for them. Offer wording only if asked.
6. Tell them to open `solutions/<id>/after.go`, compare, and read HOW, WHY and
   WHEN NOT. After they have, talk about the differences.

If the learner raises a fair objection to the task materials, say so and say
what the task got wrong. Do not defend the materials.

### Reply format

Start with a one-line verdict. Then a table of code pieces with a verdict each,
a table of the 3 answers, and at most 4 next steps.

## Mode B: hint

Ask **one** question that points at the gap. Never show code from the solution.
If they ask again, ask a narrower question. After three hints you may describe
the shape in words, still without code.
