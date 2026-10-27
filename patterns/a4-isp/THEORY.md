# A4 · ISP: depend on what you use

## Questions first

Guess before you read on. Wrong guesses are fine. After the task, fill in the
answers and fix your guesses.

1. `WelcomeFat` calls one method of a 10-method interface. What does that cost you when you test it?
   - Guess:
   - Answer:
2. Who should declare the small interface: the code that uses it, or the code that provides it? Why?
   - Guess:
   - Answer:
3. When is a big interface perfectly fine?
   - Guess:
   - Answer:

## The idea

**Problem.** A function asks for a big interface but uses one method of it.
Anyone who wants to stand in for it, such as a test, must supply all of it.

**Fix.** Let the function declare the small interface it really needs. Any
type with that method fits, including the big one, so the provider pays
nothing for the change.

**Analogy (home kitchen).** You ask a friend to cook the rice while you are
out. You leave a note: "rinse, 2 cups water, 15 minutes". You do not hand
over the 10-page house manual. Anyone who can follow the note can stand in.

## How the code shows it

| File | What it is |
|---|---|
| `before.go` | The problem. `FatStore` has 10 methods and `WelcomeFat` uses 1 |
| `after.go` | Your job: declare `UserGetter` and write `Welcome` |
| `isp_test.go` | `TestBefore_*` pass today. `TestAfter_*` pass when you are done. The test file shows a 10-method fake next to a 1-method fake |

Proof: a test stand-in needs 1 method instead of 10.

## Wrong when

Only one function uses the interface and it really needs everything. Splitting
then only adds names to read.

## In the wild

- Standard library: `go doc io.Reader`, then `go doc io.ReadWriter`. How many methods does each have, and who builds the bigger one from the smaller ones?
- Source: [src/io/io.go](https://github.com/golang/go/blob/master/src/io/io.go) (master may differ from your Go 1.26)
- Book: *100 Go Mistakes*, #5 "Interface pollution" ([summary](https://100go.co/#interface-pollution-5)).

## My delete-test result

**Try:** change `Welcome` back to take `FatStore`. Open the fake in
`isp_test.go`. How many methods must it have?

What got harder or impossible? One line:
