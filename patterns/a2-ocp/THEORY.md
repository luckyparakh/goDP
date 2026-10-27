# A2 · OCP: add a case without editing working code

## Questions first

Guess before you read on. Wrong guesses are fine. After the task, fill in the
answers and fix your guesses.

1. To add a student discount to `Price` in `before.go`, which existing code must you edit, and what could you break by accident?
   - Guess:
   - Answer:
2. How could a new discount be added without touching `Price` at all? One sentence.
   - Guess:
   - Answer:
3. What does that approach cost you that the `switch` did not?
   - Guess:
   - Answer:

## The idea

**Problem.** A new case means editing code that already works. You then have
to re-check every old case, in case the edit broke one.

**Fix.** Keep the cases in a lookup table that other code adds to. A new case
is a new entry in a new file, and no existing code is edited.

**Analogy (power strip).** A lamp wired straight into the wall: to add a fan
you open the wall. A power strip with free sockets: the new gadget just plugs
in and nothing is rewired.

## How the code shows it

| File | What it is |
|---|---|
| `before.go` | The problem. `Price` is a `switch` |
| `after.go` | Your job: `Quote` with a lookup table, then a new `student` file |
| `ocp_test.go` | `TestBefore_*` pass today. `TestAfter_*` pass when you are done |

Proof: after you commit part 1, adding the `student` tier shows one added file
and no modified files (see the commands in `after.go`).

## Wrong when

New cases almost never arrive. The table is extra machinery, and the rules are
spread over files instead of sitting in one switch.

## In the wild

- Standard library: `go doc database/sql.Register`. How does a database driver get added without editing `database/sql`? What happens if two use the same name?
- Source: [src/database/sql/sql.go](https://github.com/golang/go/blob/master/src/database/sql/sql.go) (master may differ from your Go 1.26)
- Book: *100 Go Mistakes*, #3 "Misusing init functions" ([summary](https://100go.co/#misusing-init-functions-3)). Registries like this one rely on `init()`.

## My delete-test result

**Try:** put the `switch` back (use `Price`) and add `student` again. Count the
files you edit.

What got harder or impossible? One line:
