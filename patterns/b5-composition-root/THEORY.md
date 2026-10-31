# B5 · Composition root: build everything in main

## Questions first

Guess before you read on. Wrong guesses are fine. After the task, fill in the
answers and fix your guesses.

1. In `before/`, what must a test set up before it can call the handler, and why should a handler test have to care?
   - Guess:
   - Answer:
2. If every constructor takes its dependencies as arguments, who is left to call those constructors?
   - Guess:
   - Answer:
3. `main` will get longer. When is that fine, and when is it a smell?
   - Guess:
   - Answer:

## What the principle says

**Composition Root**, named by Mark Seemann in *Dependency Injection in .NET*:
a single place in an application, as close to its entry point as possible,
where the parts are put together. Only applications have one. A library takes
what it needs and leaves the building to whoever uses it. (This is a
paraphrase from secondary sources; the original wording was not checked.)

It is the other half of dependency injection. Once every part accepts its
dependencies instead of building them, something still has to build them, and
the rule is that it happens in one place.

The textbook violation is a constructor that quietly builds its own
dependency, or reads config, deep inside the program.

**What changes in Go:** you do not need a container or a framework. `main`
calls the constructors in order and passes each result to the next. Because
interfaces are satisfied implicitly, the packages never need to know about
each other's concrete types.

The usable question is: *who calls `New`? If it is not `main` or a test, why
not?*

Cross-principle view: [SOLID.md](../SOLID.md), the "D" row.

## The idea

**Problem.** Each package builds what it needs, so the program's wiring is
spread over several packages, config is read deep inside, and a test of the top
layer drags in the bottom one.

**Fix.** Each part takes what it needs as an argument. `main` reads the config
and calls every constructor, in one place you can read top to bottom.

**Analogy (the house).** Each worker on your house buys their own material from
whichever shop they like, so nobody can tell you what went into the walls, and
checking the electrician's work means chasing their shop too. Switch to a
labour-only contract: you buy everything and hand each worker what they need.
Now one bill says what went into the house, and you can try the electrician on
a spare board without involving any shop.

## How the code shows it

| File | What it is |
|---|---|
| `before.go` | The problem, in comments |
| `before/handler`, `before/service`, `before/store` | Each builds its own dependency. Do not edit |
| `after.go` | The GOAL, with a `grep` you can run |
| `after/handler`, `after/service`, `after/store` | Your job. They take arguments but ignore them |
| `after/cmd/server/main.go` | Your job: `build`, the composition root |
| `root_test.go`, `after/cmd/server/main_test.go` | `TestBefore_*` pass today. `TestAfter_*` pass when you are done |

Proof: the `grep` in `after.go` prints nothing, and the handler and service
each work with a fake, with no file and no environment variable.

## Wrong when

A tiny program with one implementation of everything and no tests that need
fakes. Building in place is then simpler. And a dependency-injection framework
is not needed for what ten lines of `main` can do.

## FAQ

**Is a package-level `var db = openDB()` the same problem?**

Yes, and worse. It runs as soon as the package is imported, before `main` can
read any config, and a test cannot replace it without changing global state.
A2's registry is the narrow exception: it is filled once in `init()` and only
read after that.

## In the wild

- Standard library: `go doc net/http.Server`. Which field does the server fill in for itself if you leave it nil, and what does that tell you about where the choice should be made?
- Source: [src/net/http/server.go](https://github.com/golang/go/blob/master/src/net/http/server.go) (master may differ from your Go 1.26)
- Book: none verified yet.

## My delete-test result

**Try:** change `service.New` back so it ignores `st` and calls
`store.New(os.Getenv("ORDERS_FILE"))` itself, and run
`TestAfter_ServiceTakesAFakeStore`.

What got harder or impossible? One line:
