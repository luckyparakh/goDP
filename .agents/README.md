# .agents: instructions for any AI assistant

This folder is plain Markdown plus one shell script, so it works with any AI
tool. Nothing here depends on one vendor.

**Start by reading [rules.md](rules.md).** Then, when the user's message matches
a trigger below, open that file and follow it.

| The user says | Open | Kind |
|---|---|---|
| "I attempted `<id>`", "review `<id>`", "hint for `<id>`" | [skills/review-attempt.md](skills/review-attempt.md) | skill |
| "new pattern", "next pattern", "generate `<id>`" | [skills/new-pattern.md](skills/new-pattern.md) | skill |
| "check `<id>`", "QA `<id>`" | [agents/check-pattern.md](agents/check-pattern.md) | independent reviewer |
| "quiz me", "revise", "drill" | [skills/revise.md](skills/revise.md) | skill |
| "blog post for `<id>`", "write up `<id>`" | [skills/publish-post.md](skills/publish-post.md) | skill |

- A **skill** is a set of steps you follow in the current conversation.
- An **agent** file is a brief for a separate, fresh reviewer. If your tool can
  start a sub-agent or a new session, hand it the file so it reads the work
  without your reasoning in its head. If not, follow it yourself and be strict.
- `scripts/check-pattern.sh <id> [start|done]` runs the mechanical checks.

## How a tool finds this folder

`CLAUDE.md` at the repo root points here, and Claude Code reads that file by
itself. For any other tool, put one line in its own instructions file: "Read
`.agents/README.md` first." These files are not registered as native skills or
slash commands in any tool, so the user triggers them by the phrases above.
