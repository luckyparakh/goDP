# Agent: check-pattern (independent QA)

You check one pattern folder and report. **You do not edit any file.** If your
tool can start a fresh session or sub-agent, run this there so you read the work
without the author's reasoning.

Input: a pattern id such as `a2-ocp`, and a mode: `start` (the folder was just
generated) or `done` (the learner finished it).

Read first: [../rules.md](../rules.md).

## 1. Mechanical checks

Run `bash .agents/scripts/check-pattern.sh <id> <mode>` and paste the output.

## 2. Content checks (the script cannot do these)

| # | Check |
|---|---|
| a | `THEORY.md` is plain English. Jargon is explained or avoided. The analogy is 3 to 4 sentences, stays within the allowed domains, and makes no claim about the real world. People really behave the way the story says. The "before" version loses something in the story that someone needs, and that loss matches what the `TestBefore_*` tests prove. |
| b | No spoilers. `after.go` and `THEORY.md` do not contain the solution's key code. Compare with `solutions/<id>/`. |
| c | `before.go`: the PROBLEM and TRY IT make sense, and every claim in the comments is true. Run what you can to prove it. |
| d | The 3 questions can each be answered with a guess, and the `Questions first` section does not give the answers away. |
| e | "In the wild": run `go doc <symbol>`. The book entry matches the title and number on https://100go.co if you can fetch it, or it says "none verified yet". |
| f | The pattern has a row in `BLUEPRINT.md` and its proof matches the test. |
| g | Concurrency only: `go test -race -count=50 -shuffle=on ./patterns/<id>/`. Report flaky tests. |
| h | `THEORY.md` explains the **principle**, not only this code's problem. It must name the principle, give the classic formulation, and say what changes in Go. Ask yourself: could a learner who searched the web for this principle's name map what they found onto this task? If the task's mechanism is not obviously the principle, the file must spell out the chain between them. Only skip this for a task with no named principle behind it. |

## 3. Report

A table with `check | result | evidence`. Keep what you **verified** apart from
what you **could not verify**. End with `PASS`, or a numbered fix list. Never
claim a result you did not run. Never fix anything yourself.
