# A5 · DIP: let the user of a thing say what it needs

## Questions first

Guess before you read on. Wrong guesses are fine. After the task, fill in the
answers and fix your guesses.

1. What must you set up to test the greeting in `before.go`, and which failure can you not test?
   - Guess:
   - Answer:
2. If the greeting code asked its question through a small interface, who decides what that question looks like?
   - Guess:
   - Answer:
3. Name one situation where adding that interface would be a waste.
   - Guess:
   - Answer:

## What the principle says

**Dependency Inversion Principle**, from Robert C. Martin: high-level code
should not depend on low-level code, and both should depend on an abstraction.
The "inversion" is the direction of that dependency. Normally the policy reaches
down and grabs a mechanism. Inverted, the mechanism fits an abstraction that the
policy side owns.

**It is not dependency injection.** Injection is the technique of handing a
dependency in, through a constructor or a parameter. Inversion is about *who
owns the interface*. You can inject a dependency and still have it backwards, if
the interface is declared by whoever provides it.

In Go, that ownership is the whole mechanism: the consumer declares the
interface in its own package, so the import arrow points from the provider to
the consumer instead of the other way round.

The usable question is: *who owns the interface, the user or the provider?*

Cross-principle view: [SOLID.md](../SOLID.md).

## The idea

**Problem.** The code that decides what to do reaches down and grabs the
real thing, here the file system. To test it you must build the real thing,
and you cannot make it fail the way you want.

**Fix.** The code asks a small question through an interface that *it* owns.
The real file reader, a plain map, or a stand-in that always fails can all
answer it. The arrow now points up: the supplier fits the need.

**Analogy (water at home).** The washing machine is plumbed straight into the
borewell. To check what it does when the water stops, you have to actually empty
the borewell. Give it a normal inlet instead, fed by the borewell, a tanker or a
bucket, and "no water" becomes a tap you turn off.

## How the code shows it

| File | What it is |
|---|---|
| `before.go` | The problem. `FileService` reads a file itself |
| `after.go` | Your job: `Greeter` with a small `Users` interface it owns. It must not import `os` |
| `dip_test.go` | `TestBefore_*` pass today. `TestAfter_*` pass when you are done |

Proof: the failure path is one line in a test, where before it needed a real
broken file.

## Wrong when

The real thing is the only source forever, and you never test the failure.
The interface then only adds a layer to read through.

## In the wild

- Standard library: `go doc net/http.FileServer`. Which side owns the interface it takes?
- Source: [src/net/http/fs.go](https://github.com/golang/go/blob/master/src/net/http/fs.go) (master may differ from your Go 1.26)
- Book: *100 Go Mistakes*, #6 "Interface on the producer side" ([summary](https://100go.co/#interface-on-the-producer-side-6)). Compare it with where `Users` lives in your code. The title is checked; that it fits this lesson is my reading.

## My delete-test result

**Try:** use `FileService` again and try to test "the user store is down". Can
you choose the error?

What got harder or impossible? One line:
