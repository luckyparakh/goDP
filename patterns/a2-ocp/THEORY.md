# A2 · OCP: add a case without editing working code

## Questions first

Guess before you read on. Wrong guesses are fine. After the task, fill in the
answers and fix your guesses.

1. To add a student discount to `Price` in `before.go`, which existing code must you edit, and what could you break by accident?
   - Guess: I should be editing Price func esp. the swicth case and may existing func or its test case may break. Need to retest each.
   - Answer: Same as above
2. How could a new discount be added without touching `Price` at all? One sentence.
   - Guess: Create a map of tier name vs fucn name and add a register func. Create a global map and add all func there and call upon from price func depending on tier name.
   - Answer: Same as above
3. What does that approach cost you that the `switch` did not?
   - Guess: A global variable, extra boilerplate code, less readibilty etc.
   - Answer: Same as above and also a init func, may be extra files

## What the principle says

**Open/Closed Principle**, from Bertrand Meyer (1988): software should be open
for extension but closed for modification. You should be able to add behaviour
without editing what is already there.

Meyer's original version extended by subclassing. Go has no inheritance, so the
extension point has to be something else: a function value, an interface, or a
lookup table that new code adds itself to. This task uses the lookup table.

The cost is that you must build the extension point *before* you need it, which
is why chasing this principle directly tends to produce abstraction nobody asked
for. The usable question is: *to add a case, do I add or do I modify?*

Cross-principle view: [SOLID.md](../SOLID.md).

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

## FAQ

**Is a global variable bad? The lookup table is one.**

Usually yes. Mutable global state can be changed by any code, hides a
dependency, lets tests interfere with each other, depends on initialization
order, and can race if it is written while being read.

A registry is the accepted exception. The standard library does the same:
`database/sql` keeps a `drivers` map behind a mutex, and `image.RegisterFormat`
is a registry too. Treat it as acceptable when all of these hold:

- it is written only during `init()`
- it is read-only after that
- it is unexported, with a small `Register` function in front
- it is locked if writes could ever happen later
- tests do not need to change it

It is not acceptable for runtime business state, config, caches or counters.

**Can I avoid the global?** Yes. Make a value that owns its own map, build it in
`main`, and register each tier there. That is easier to test. But you lose the
part that makes this lesson work: a new file registers itself. A file that
registers itself needs a package-level table. So you trade the "add a file, edit
nothing" win for explicitness, and neither choice is free.

## In the wild

- Standard library: `go doc database/sql.Register`. How does a database driver get added without editing `database/sql`? What happens if two use the same name?
- Source: [src/database/sql/sql.go](https://github.com/golang/go/blob/master/src/database/sql/sql.go) (master may differ from your Go 1.26)
- Book: *100 Go Mistakes*, #3 "Misusing init functions" ([summary](https://100go.co/#misusing-init-functions-3)). Registries like this one rely on `init()`.

## My delete-test result

**Try:** put the `switch` back (use `Price`) and add `student` again. Count the
files you edit.

What got harder or impossible? One line: It increased test-radius as existing work also got edited.
