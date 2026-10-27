# Skill: publish-post

Use when the learner asks for a blog post or write-up for a pattern id, for
example "blog post for a1-srp" or "write up a2-ocp".

Read first: [../rules.md](../rules.md), then `patterns/<id>/THEORY.md`,
`before.go`, the learner's finished `after.go`, `solutions/<id>/after.go`, the
test file, and the `<id>` row in `PROGRESS.md`.

Write to `blog/<id>-<slug>.md`, for example `blog/a1-srp-one-job-per-piece.md`.
Plain Markdown, no front matter. Create `blog/` if it is missing.

## Before you write: the guard

Stop and say so if any of these is true. Do not write the file.

| Check | Why |
|---|---|
| The pattern has no row in `PROGRESS.md` | Without "what I got wrong" the post is the same as every other post on this pattern |
| `go test ./patterns/<id>/` does not pass | There is no proof to show yet |
| The delete-test line in `THEORY.md` is empty | The post's best section would be missing |

## Structure

Target 700 to 1200 words. On a small pattern some sections are two sentences,
which is fine. Do not pad to reach a length.

| Section | What goes in it |
|---|---|
| `# <plain title>` | What it does, not the principle's name. For example "One job per piece" |
| `**TL;DR**` | Two lines: what breaks, what fixes it, and the one thing you can count |
| `## The everyday version` | The analogy, pain first. 3 to 4 sentences, same domain rules as `rules.md` |
| `## How you'd fix it in real life` | The analogy's fix in 2 to 3 sentences, then a `> Where this analogy breaks:` line |
| `## The same pain, in code` | Prose. Name each reason the code would change. End with the cost you can count |
| | A `go` snippet from `before.go`. Trim it to what the point needs |
| `## The fix` | Prose: how it works and why it helps |
| | A `go` snippet from the finished code |
| `## Proof` | The command, then the **real** pasted output: failing first, passing after |
| `## What I got wrong` | `TODO-yours`. See below |
| `## The delete test` | What was put back, and what became harder. Use the learner's line |
| `## Do and don't` | A two-column table, at most 4 rows |
| `## When not to use this` | From WHEN NOT TO USE in the sealed solution, in your own words |
| `## Further reading` | The three "In the wild" lines, plus a link to `patterns/<id>/` in the repo |

## The learner's voice

Two sections carry the learner's voice, and a post that sounds like you is a
worse post. Write them as:

```markdown
## What I got wrong

> TODO-yours. From PROGRESS.md: <paste their row verbatim>.
> Write 3 or 4 sentences: what you expected, what actually happened, and how
> you found it.
```

Do the same for the closing paragraph. Everything else you write in plain prose.
Tell the learner in your reply which sections are theirs to fill.

## Rules for the post

- **It publishes nothing.** You write a local file. Posting is the learner's call. Never post to any site, and never offer to.
- **Never invent terminal output.** Run the command and paste what it printed. To show the failing state, run the stub in a throwaway module outside the repo rather than editing `patterns/`. Say which state each block is.
- **Name the platform under pasted output.** For example "Go 1.26, Windows". The same code can behave differently elsewhere, so a reader needs to know where it ran.
- **Mark every trim.** If you shorten a snippet or an output line, put `...` where the cut is. Never present a trimmed block as the whole thing.
- **Quote the standard library, do not copy it.** Link to the source file. Keep every code block in the post original to this repo.
- **Nothing from work.** No company, customer or colleague names, and no code from a private repo. This repo is public.
- **No filler.** No "in today's fast-paced world". Start at the problem.
- **One post at a time.** Same as one task at a time.

## After writing

Reply with: the file path, the word count, which sections are `TODO-yours`, and
anything you could not verify. Do not commit.
