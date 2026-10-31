# B2 · Adapter: keep the database's words out of your service

## Questions first

Guess before you read on. Wrong guesses are fine. After the task, fill in the
answers and fix your guesses.

1. In `before/`, how does a caller find out an order is missing, and what must it import to do that?
   - Guess:
   - Answer:
2. The adapter gets a driver error it has never seen. Should callers be able to reach it with `errors.As`?
   - Guess:
   - Answer:
3. Why should `ErrNotFound` be declared in `orders`, and not in the adapter?
   - Guess:
   - Answer:

## What the principle says

**Adapter**, from *Design Patterns* by Gamma, Helm, Johnson and Vlissides
(1994): convert the interface of a class into another interface that clients
expect. The textbook case wraps a third-party class so the client sees only its
own interface.

**What changes in Go:** the adapter is a small struct that holds the foreign
thing, here a `*sql.DB`, and has the one method the consumer declared. It
satisfies that interface implicitly, and lives in its own package that imports
the consumer (B1).

The book says little about errors, but in a service that is where the leak
usually is. A driver's errors are its own vocabulary. If they pass through, every
caller learns that vocabulary, and the database is no longer swappable. So the
adapter translates errors as well as calls: it maps the ones it knows into the
service's words, and hides the rest (A6's "hide it on purpose with `%v`").

The usable question is: *can a caller of my service tell which database I use?*

Cross-principle view: [SOLID.md](../SOLID.md), the "D" row.

## The idea

**Problem.** The service runs SQL itself and passes the database's errors up.
To spot a missing order, every caller must import `database/sql`, and any
caller can reach into the driver's error type and depend on its codes.

**Fix.** The service declares what it needs and the errors it understands. An
adapter holds the database, answers the service's question, and translates
every error at the boundary.

**Analogy (the mechanic).** The garage hands you the scanner printout, "P0300",
and now you are searching codes on your phone to decide whether to drive home.
Next month the garage buys a different scanner and every code you learned means
nothing. A good mechanic reads the printout and tells you one of two things in
your words, "nothing wrong, drive home" or "leave it here", and keeps the
printout in the file.

## How the code shows it

| File | What it is |
|---|---|
| `before.go` | The problem, in comments |
| `before/orders` | The service runs SQL and leaks its errors. Do not edit |
| `fakedb/` | A stand-in driver, standard library only. Do not edit |
| `after.go` | The GOAL |
| `after/orders`, `after/sqlstore` | Your job |
| `adapter_test.go` | `TestBefore_*` pass today. `TestAfter_*` pass when you are done |

Proof: `go list -deps` of `after/orders` has no `database/sql`, and the driver's
error cannot be reached from the caller, though its text is still in the
message.

## Wrong when

A small tool over one database you will never swap, where nothing above checks
which error it got. The adapter is then a layer to read through.

## FAQ

**Should every database error be hidden?**

No. Hide what is specific to the driver: its error types and codes. Let
`context.Canceled` and `context.DeadlineExceeded` through with `%w`, because
they are the caller's own signal that it gave up or ran out of time (A7), not a
storage detail. The GOAL asks for exactly that.

## In the wild

- Standard library: `go doc database/sql.ErrNoRows`. A caller must import `database/sql` to check for it. Which of your packages should ever need to?
- Source: [src/database/sql/sql.go](https://github.com/golang/go/blob/master/src/database/sql/sql.go) (master may differ from your Go 1.26)
- Book: *100 Go Mistakes*, #49 "Ignoring when to wrap an error" ([summary](https://100go.co/#ignoring-when-to-wrap-an-error-49)). The title is checked; that it fits this lesson is my reading.

## My delete-test result

**Try:** in `sqlstore`, wrap the unknown driver error with `%w` instead of
`%v`, and run `TestAfter_DriverErrorIsHidden`.

What got harder or impossible? One line:
