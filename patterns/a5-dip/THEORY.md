# A5 · DIP: let the user of a thing say what it needs

## Questions first

Guess before you read on. Wrong guesses are fine. After the task, fill in the
answers and fix your guesses.

1. What must you set up to test the greeting in `before.go`, and which failure can you not test?
   - Guess: A folder and any error other than OS one
   - Answer: You need a real temp folder with a real file in it, just to check a string join. The failure you cannot test is any failure the OS does not hand you: "store is down", a timeout, or a sentinel error you want to assert with errors.Is
2. If the greeting code asked its question through a small interface, who decides what that question looks like?
   - Guess: Greeting code decides, ie consumer
   - Answer: The consumer. The code that needs the answer declares the interface, so the shape of the question is set by the side being served
3. Name one situation where adding that interface would be a waste.
   - Guess: When we need only one way to answer (not swapable) and fix set of errors
   - Answer: One source forever and no failure path worth faking. Then the interface is a layer you read through for nothing

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

**Analogy 1 (water at home).** The washing machine is plumbed straight into the
borewell. To check what it does when the water stops, you have to actually empty
the borewell. Give it a normal inlet instead, fed by the borewell, a tanker or a
bucket, and "no water" becomes a tap you turn off.

**Analogy 2 (eating out).** The menu is written by the kitchen, so your order
bends to their list. Write the slip yourself, one line in your terms, and any
kitchen can take it. The menu is an interface too. The difference is whose hand
wrote it.

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

## FAQ

**Why is it called "inversion"? In Go the consumer never imports the provider
and the provider never imports the consumer. Nothing looks inverted.**

The import graph is not the thing being inverted. There are two arrows between
the greeting code and the file reader:

| Arrow | `before.go` | `after.go` |
|---|---|---|
| Who calls whom at run time | Greeter → file read | Greeter → file read, unchanged |
| Who must exist to compile | Greeter → `os` | gone |

The call still goes down. The compile-time dependency does not. That is the
inversion: the source arrow stopped following the call arrow.

The word comes from Java and C++. There the provider writes `implements Users`,
so it must name the consumer's package, and the arrow really does turn around
and point at the consumer. Go's implicit satisfaction deletes the arrow instead
of reversing it. Same ownership, no import.

**Does the Go import graph show anything at all, then?**

Yes, once the two sides are in different packages. Same code both times, only
the `type Users interface` line moves. Listed with
`go list -f '{{.ImportPath}} -> {{.Imports}}' ./...`:

```
# interface declared in the consumer
consumer -> []
provider -> [os path/filepath strings]

# interface declared in the provider
consumer -> [provider]
provider -> [os path/filepath strings]
```

Both compile. Both hand the dependency in through a constructor. Only the first
is inverted. In the second the greeting code depends on `os` again, one hop
away.

So "the consumer never imports the provider" is not a fact about Go. It is a
fact about putting the interface on the consumer side. That placement is the
whole choice.

**What is the one-line test?**

Change `Name(id string)` to `Name(ctx context.Context, id string)`. Who has to
follow? If every provider must change, the consumer owns the contract and the
dependency is inverted. If instead the consumer must rewrite itself to match a
signature the provider published, it is not.

**Why can't this folder show it?**

Everything here is one package, so there is no import arrow to look at. The
"must not import `os`" rule stands in for the first half of DIP: high-level code
must not depend on low-level code. The second half, ownership, you can only see
by noticing that `type Users interface` sits in `after.go` next to `Greeter`,
not next to `FileService`.

## In the wild

- Standard library: `go doc net/http.FileServer`. Which side owns the interface it takes?
- Source: [src/net/http/fs.go](https://github.com/golang/go/blob/master/src/net/http/fs.go) (master may differ from your Go 1.26)
- Book: *100 Go Mistakes*, #6 "Interface on the producer side" ([summary](https://100go.co/#interface-on-the-producer-side-6)). Compare it with where `Users` lives in your code. The title is checked; that it fits this lesson is my reading.

## My delete-test result

**Try:** use `FileService` again and try to test "the user store is down". Can
you choose the error?

What got harder or impossible? One line:
