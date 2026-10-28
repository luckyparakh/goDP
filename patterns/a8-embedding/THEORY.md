# A8 · Embedding is not inheritance

## Questions first

Guess before you read on. Wrong guesses are fine. After the task, fill in the
answers and fix your guesses.

1. `Derived` embeds `Base` and replaces `Name`. What does `Derived{}.Hello()` return, "hello from base" or "hello from derived"?
   - Guess:
   - Answer:
2. Can you pass a `Derived` to a function that takes a `Base`? Why or why not?
   - Guess:
   - Answer:
3. If embedding is not inheritance, what is it useful for?
   - Guess:
   - Answer:

## The idea

**Problem.** Coming from Java or C++, embedding looks like inheritance. It is
not. The embedded value is a separate field that Go forwards methods to. A
method written on the inner type still sees the inner type, so it cannot see
what the outer type replaced.

**Fix.** Do not expect the inner method to find your change. Pass the varying
behaviour in explicitly, for example as an interface. Use embedding for what it
does well: forwarding many methods for free, and replacing just one.

**Analogy (the tiffin).** A tiffin carrier has a dabba of dal inside. You eat
the dal straight from the carrier, so it starts to feel like the carrier *is*
the dal. It is not. Ask the dal what it is and it answers "dal", whatever you
wrote on the sticker outside the carrier.

## How the code shows it

| File | What it is |
|---|---|
| `before.go` | The surprise. `Derived` replaces `Name`, but `Hello` still uses `Base`'s |
| `after.go` | Your job: `Greet`, which takes the varying part as a `Namer`, and `Shout`, which replaces one method |
| `embedding_test.go` | `TestBefore_*` pass today. `TestAfter_*` pass when you are done |

Proof: `Greet(Derived{})` says "derived", which `Derived{}.Hello()` never did.

## Wrong when

You embed only to get an "is-a" relationship, or you embed a type with many
exported methods and do not want all of them in your API. Embedding promotes
every one of them.

## In the wild

- Standard library: `go doc bufio.ReadWriter`. Which methods does it get without writing any of them, and from what?
- Source: [src/bufio/bufio.go](https://github.com/golang/go/blob/master/src/bufio/bufio.go) (master may differ from your Go 1.26)
- Book: *100 Go Mistakes*, #10 "Not being aware of the possible problems with type embedding" ([summary](https://100go.co/#not-being-aware-of-the-possible-problems-with-type-embedding-10)).

## My delete-test result

**Try:** make `Greet` call `Base{}.Hello()` instead of using its argument, and
run `TestAfter_GreetUsesTheNamerItIsGiven`.

What got harder or impossible? One line:
