# A4 · ISP: depend on what you use

## Questions first

Guess before you read on. Wrong guesses are fine. After the task, fill in the
answers and fix your guesses.

1. `WelcomeFat` calls one method of a 10-method interface. What does that cost you when you test it?
   - Guess: Mock rest of func while testing, testing radius increases if logic of the func changes
   - Answer: Same as above
2. Who should declare the small interface: the code that uses it, or the code that provides it? Why?
   - Guess: Code that uses it, consumer should defined needs not the provider.
   - Answer: Same as above
3. When is a big interface perfectly fine?
   - Guess: No idea
   - Answer: When there is only one function using it and it needs all the interface methods.

## What the principle says

**Interface Segregation Principle**, from Robert C. Martin: clients should not
be forced to depend on methods they do not use. A wide interface gets split into
the narrow roles that callers actually need.

Where a type must declare which interfaces it implements, that split is the
provider's decision and costs real work. Go satisfies interfaces implicitly, so
the *consumer* can declare the narrow interface it needs and every existing type
fits it unchanged. That makes this principle close to free in Go, and it is the
idea behind the Go proverb "the bigger the interface, the weaker the
abstraction".

The usable question is: *do I depend on more than I use?*

Cross-principle view: [SOLID.md](../SOLID.md).

## The idea

**Problem.** A function asks for a big interface but uses one method of it.
Anyone who wants to stand in for it, such as a test, must supply all of it.

**Fix.** Let the function declare the small interface it really needs. Any
type with that method fits, including the big one, so the provider pays
nothing for the change.

**Analogy (the gate register).** The guard only needs your name. The register
has ten columns, so you fill nine of them for nobody to read. A slip asking only
for your name would do, and anyone who can fill the whole register can certainly
write their name on a slip.

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

## FAQ

**If a function really needs all ten methods, is that not an SRP problem?**

Yes, and spotting that matters. ISP and SRP look at different things, so fixing
one can leave the other standing.

| Depends on | Actually uses | Verdict |
|---|---|---|
| 10 | 1 | ISP violation. Narrow the interface |
| 10 | 10 | ISP is satisfied. Now ask SRP's question |
| 1 | 1 | Both fine |

The middle row is the trap. Declare a consumer-side interface with exactly the
ten methods you call and ISP is satisfied, because you now depend on precisely
what you use. Nothing got smaller. You made the dependency honest without making
the function simpler, and the fix there is to split the function, not the
interface.

There are also two separate things to point SRP at in `before.go`:

- **The function.** How many different requirements would make you edit it? Ten
  means split it.
- **The interface.** `FatStore` holds user methods, order methods and `Ping`.
  That is at least three reasons to change, so the interface breaks SRP too.

The two principles join up by reasoning: SRP asks for one reason to change, and
an interface mixing users with orders has two. Point SRP at an interface and you
get narrow role interfaces, which is what ISP asks for. So `before.go` shows both
smells at once, ISP from the caller's side and SRP from the interface's.

**What makes a good wrapped error message?**

The message is a path. Each layer puts its own context in front, separated by
`: `, so you read it left to right from outermost to innermost:

```
handle GET /welcome: welcome "u1": query users: not found
```

- **Name the operation, not the failure.** The value is already an error, so
  "error while getting user" repeats the word at every layer and adds nothing.
  The same chain written badly: `error in handle GET /welcome: error while
  getting user u1: error in query users: not found`.
- **Start lowercase.** Go's code review comments say error strings "should not
  be capitalized (unless beginning with proper nouns or acronyms) or end with
  punctuation", because they are "usually printed following other context". A
  capital letter lands in the middle of the path above.
- **Use `%q` for values** such as an id or a key. With `%s`, an empty or
  whitespace value disappears from the message. With `%q` you still see `""`.
- **Why `%w` and not `%v`** is the subject of [A6](../a6-errors/THEORY.md).

## In the wild

- Standard library: `go doc io.Reader`, then `go doc io.ReadWriter`. How many methods does each have, and who builds the bigger one from the smaller ones?
- Source: [src/io/io.go](https://github.com/golang/go/blob/master/src/io/io.go) (master may differ from your Go 1.26)
- Book: *100 Go Mistakes*, #5 "Interface pollution" ([summary](https://100go.co/#interface-pollution-5)).

## My delete-test result

**Try:** change `Welcome` back to take `FatStore`. Open the fake in
`isp_test.go`. How many methods must it have?

What got harder or impossible? One line:
