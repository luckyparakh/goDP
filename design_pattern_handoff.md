# design_pattern.md — Project Context & Handoff

> **Purpose of this file.** It is the complete, self-contained record of a long design
> discussion about learning design patterns in Go. It exists so an agent (Claude Code in
> VS Code) can read it once and continue the work without re-asking anything.
>
> **This file contains no Go code on purpose.** Code is to be generated from the specs here.
>
> **If you are the agent: read §12 (Agent Instructions) first, then §0, then the rest.**

---

## 0. The one-paragraph summary

A backend services engineer wants to learn design patterns in Go **practically** — the way
a senior architect would use them in code review — not by memorising a 22-pattern catalog.
The agreed method is a **7-step "Pattern Dojo"** where every pattern must be proven by an
executable test whose output changes when the pattern is applied, and then **deleted** to see
what breaks. The curriculum is **39 tasks in 6 blocks (A–F), 40–60 hours, 12–14 weeks**.
Theory is written language-agnostic; every implementation is Go. The work is intended to be
**published publicly**. Block A task 1 (SOLID) has been fully specified as five code labs.
Nothing has been built by the learner yet.

---

## 1. Learner profile & hard constraints

| Item | Value |
|---|---|
| Role | Backend **services** engineer (HTTP/gRPC services, not data pipelines) |
| Go version | **1.26.0** |
| `go get` / third-party modules | **Allowed** |
| Editor | VS Code + Claude Code |
| Books owned | *100 Go Mistakes* (Harsanyi), *Concurrency in Go* (Cox-Buday), *Designing Data-Intensive Applications* (Kleppmann), *Learning Go* (Bodner) |
| Also owns | Refactoring Guru's *Dive Into Design Patterns* |
| Explicit scope limit | **System design is OUT of scope for now.** Learner said: "don't want to learn many things at once." DDIA stays on the shelf. |
| Learner's stated ethic | "I will not look at the solution before making a good attempt." Honour this — it shapes the folder layout. |
| Learning style signals | Wants tables not prose; got lost once when given too many files at once; wants each claim proven, dislikes hand-waving ("no bluff"); asks for first-principles reasoning |

### Preferred tone for this project
Direct, senior-engineer register. State what is verified vs unverified. Own mistakes plainly.
Do not pad. Tables over paragraphs where a table fits. Push back when the learner is about
to add ceremony.

---

## 2. How this started (so context isn't lost)

The discussion began with a concrete question: *"My codebase has a `DBStore` interface and a
`DBAdapter` struct that implements it. I understand it helps with mocking. What else does it
give me?"*

The answer listed ~11 benefits beyond mocking (technology isolation, leak prevention, error
translation, model translation, decorators, multiple implementations, DB policy in one place,
composition root, clear team contract, compile-time contract check, and mocking last). Three
throwaway demo repos were built to prove them. **Those repos are superseded by this plan and
should be considered scratch work.** The lasting output of that phase was the *method*:

> A benefit that can be demonstrated by a command whose output changes is real.
> A benefit that can only be asserted in prose is a sales pitch.

That method became the Pattern Dojo in §4.

---

## 3. Research findings (all verified during the discussion — do not re-litigate)

### 3.1 Refactoring Guru's Go catalog is unreliable as a Go teacher

Keep the book for **concepts** (Intent, Problem, Applicability, Pros/Cons — these are
language-agnostic and good). Ignore its Go code. Evidence gathered:

**(a) Their Singleton example contains a data race — reproduced, not assumed.**
Both variants on their Singleton page were run under `go run -race`. Both reported
`Found 1 data race(s)`. The bug is identical in each: the `nil` check on the package-level
pointer sits **outside** the lock / outside `once.Do`, so one goroutine reads the pointer
while another writes it. Their "another example" wraps `sync.Once` but keeps the
unsynchronized pre-check, which defeats the point.
*Fix confirmed clean under `-race`:* let `once.Do` be the only access to the variable.
On Go 1.21+ `sync.OnceValue` does it in one line.
**Teaching value: this is the model proof. Reuse this experiment in Block D.**

**(b) Their Factory Method example is Java in Go clothing.**
It defines an `IGun` interface with `setName`/`setPower`/`getName`/`getPower`, a `Gun` struct,
and `Ak47`/`musket` structs embedding `Gun` to simulate inheritance. Hungarian `I` prefix,
getters/setters on a plain struct, embedding-as-inheritance. The page itself concedes the
classic Factory Method cannot be implemented in Go for lack of classes and inheritance, and
shows Simple Factory instead. **That concession is the real lesson and it applies to roughly
a third of the catalog.**

**(c) The catalog is the unchanged 1994 list** — 5 creational, 7 structural, 10 behavioral,
identical to every other language page. Nothing added for channels, `context`, or errors as
values, which is where Go's actual design decisions live.

### 3.2 Standard library pattern locations (verified)

Verified by grepping a real Go source tree. **Line numbers were from go1.27.3 and shift
between versions — the learner is on 1.26, so use `go doc -src <symbol>` against their own
install rather than trusting numbers.**

| Pattern | Symbol | File |
|---|---|---|
| Consumer interface (1 method) | `http.Handler` | `src/net/http/server.go` |
| **Adapter** (its own doc calls it that) | `http.HandlerFunc` | `src/net/http/server.go` |
| **Decorator ×3** | `StripPrefix`, `TimeoutHandler`, `MaxBytesHandler` | `src/net/http/server.go` |
| Decorator (I/O) | `io.TeeReader`, `io.LimitReader`, `io.NopCloser` | `src/io/io.go` |
| Composite | `io.MultiWriter` | `src/io/multi.go` |
| Strategy | `sort.Interface` | `src/sort/sort.go` |
| Lazy init done right | `sync.OnceValue`, `sync.OnceValues` | `src/sync/oncefunc.go` |
| Proxy | `httputil.ReverseProxy` | `src/net/http/httputil/reverseproxy.go` |
| Registry / driver plugin | `database/sql/driver.Driver` | `src/database/sql/driver/driver.go` |

**Three decorators in one stdlib file (`net/http`) is the best teaching material in Go.**
Confirmed separately: `golang.org/x/sync` (errgroup, singleflight) is **not** in the stdlib.

### 3.3 Functional options — sourced history

The pattern originates in Rob Pike's post *"Self-referential functions and the design of
options"*; Dave Cheney credits it as the origin in *"Functional options for friendly APIs"*,
and later noted the gRPC project evolved it into its best form. Cheney's framing: customization
is performed not with config parameters in a struct but with functions that operate on the value
itself; the result can grow over time, is self-documenting, and never requires a nil or empty
value to satisfy the compiler.
**Uber style guide threshold:** config struct under 3 options; functional options for public
APIs with 3+ optional arguments that may grow.

### 3.4 Book chapter mapping (ToCs verified except where noted)

**100 Go Mistakes** — chapters: 1 Go: simple to learn but hard to master · 2 Code and project
organization · 3 Data types · 4 Control structures · 5 Strings · 6 Functions and methods ·
7 Error management · 8 Concurrency: Foundations · 9 Concurrency: Practice · 10 The standard
library · 11 Testing · 12 Optimizations.
Specific mistakes confirmed and directly relevant:
- **#3 "Misusing init functions"** (ch.2) → task C4
- **#5 "Interface pollution"** (ch.2) → tasks B1, B2. *The single most relevant page in any of
  the learner's books for this curriculum.*
- **#48–#54** (ch.7) — panicking, when to wrap, checking error type/value inaccurately,
  handling an error twice, not handling, defer errors → task A2
- **ch.8 §8.4 "Not understanding race problems"** — data races vs race conditions, the Go
  memory model → Block D, and it explains the RG Singleton bug exactly
- **ch.8 §8.6 "Misunderstanding Go contexts"** → task A3
- A free summary of all 100 mistakes exists at `100go.co`.

**Concurrency in Go (Cox-Buday)** — chapters: 1 An Introduction to Concurrency · 2 Modeling
Your Code: CSP · 3 Go's Concurrency Building Blocks · **4 Concurrency Patterns in Go**
(confinement, the for-select loop, pipelines, fan-out/fan-in) · **5 Concurrency at Scale**
(error propagation, timeouts and cancellation, heartbeats) · 6 Goroutines and the Go Runtime.
→ ch.4 maps to D2–D4; ch.5 maps to D5–D9.

**Learning Go (Bodner)** — ToC **not verified**. Use for Block A foundations generally.

**DDIA** — deliberately unused. It is a systems book; systems are out of scope.

### 3.5 Publishing / licensing (verified)

| Source | What is permitted |
|---|---|
| Refactoring Guru **text** | Free use **with a hyperlink to the specific page** it came from; you may cite any text as long as it is not a substantial part of the cited article. |
| Refactoring Guru **code** | **Nothing.** Their Go examples repo states the code is copyright-protected and licensing questions go to their support email. Every line in this project must be original. |
| Go standard library | BSD-3-Clause. Link to it; do not vendor snippets. |
| The learner's books | Cite as further reading. Never reproduce text or code. |

**Suggested licensing for the published repo:** MIT or Apache-2.0 for `*.go`, CC-BY-4.0 for
the Markdown. Keep every `THEORY.md` free of Go so the theory stands alone in any language.

**Honest positioning note:** dozens of "design patterns in Go" repos already exist and most are
Java transliterations. The differentiator here is **not** the pattern list. It is that every
pattern carries an **executable proof** and a **recorded deletion test**. Lead the README with
that, not with a list of 22 patterns.

---

## 4. The pedagogy — the 7-step Pattern Dojo

> A pattern you cannot prove with a failing test is a pattern you have not learned.

One task = one sitting = one package = one proof.

| Step | What the learner does | Output file |
|---|---|---|
| **1 PAIN** | Run the naive version. Feel what hurts. **Do not write the pattern yet.** | — |
| **2 FORCE** | Finish the sentence *"What varies here is ___"*. Cannot finish it → the pattern does not apply. | `NOTES.md` |
| **3 SHAPE** | Intent in their own words + a 3-box sketch of who points at whom. **No Go in this step.** | `THEORY.md` |
| **4 GO FORM** | Smallest idiomatic Go carrying that shape. Often a `func` type, not an interface. | `attempt/` |
| **5 PROOF** | The exact command that must pass. Paste **real** output, not expected output. | `PROOF.md` |
| **6 WILD** | Read one real stdlib file. Write one line on how it differs from what they wrote. | `THEORY.md` |
| **7 DON'T** | Write "wrong when ___", then run the **deletion test**. | `DONT.md` |

**Step 3 must come before step 4.** That ordering is what makes the output publishable as
language-agnostic material rather than a Go tutorial.

### What counts as proof

| Proof type | Use for | Example |
|---|---|---|
| A test that could not be written before | Adapter, DI, Strategy | Cannot unit-test the service without a seam |
| A compile error | Contracts | `var _ Iface = (*T)(nil)` breaks on drift |
| `go test -race` | Any concurrency task | Reproduces real data races (see §3.1a) |
| A call count | Decorator, Proxy, cache, singleflight | 3 calls → 1 DB hit |
| `go test -bench -benchmem` | Flyweight, `sync.Pool`, buffering | allocs/op drops |
| `git diff --stat` | Facade, OCP, composition root | Adding a backend touches 1 file, not 12 |

### The deletion test (the anti-ceremony device)

After the pattern works, **delete it** and observe what breaks. If nothing breaks, the
abstraction was decoration. Whatever broke **is** the justification — record that sentence in
`DONT.md`. After ~20 tasks the learner owns a review checklist written in their own words,
which is worth more than any book.

### The review drill (twice a week, 20 minutes)

Take a real Go file — their own, a dependency's, or AI-generated — and answer in writing:
1. What varies here, and is that variation isolated or smeared across the file?
2. Which pattern is **missing**?
3. Which pattern is **over-applied** (an abstraction with exactly one implementation and one caller)?

**Question 3 matters more than question 2 for the first six months.** The predictable failure
mode after learning patterns is over-application, not under-application.

---

## 5. Where SOLID fits (resolved — first-principles version)

SOLID is **not** a separate course. It is the *why* underneath Block B, budgeted at **one task
(A1), not a week**. The learner said they had read SOLID several times and it was still hazy;
the diagnosis was that the principles are stated as rules with **no units**, so they cannot be
checked. The fix was to derive them from one question:

> **When requirement X changes, how much code must I open, re-read, re-test and risk?**

Call that the **blast radius**. All five principles are tactics for shrinking it. A second
question falls out of the first, because shrinking blast radius means swapping parts:

> **Can I replace a part without the rest of the system noticing?**

| | Name | The force it fights | Question to ask the code |
|---|---|---|---|
| **S** | Single Responsibility | A module answers to two masters, so either can break it | "How many *different* requirements would make me edit this?" |
| **O** | Open/Closed | New cases force edits to working, tested code | "To add a case, do I *add* or do I *modify*?" |
| **L** | Liskov Substitution | A replacement type-checks but misbehaves | "Does every implementation keep the same promises?" |
| **I** | Interface Segregation | A consumer is coupled to methods it never calls | "Do I depend on more than I use?" |
| **D** | Dependency Inversion | Policy reaches down and grabs a mechanism | "Who owns the interface — the user or the provider?" |

### Weighting in Go (this is the honest part)

- **I and D are load-bearing.** Implicit interface satisfaction lets the *consumer* declare the
  interface, so segregation and inversion are nearly free. Same idea as the proverb
  *the bigger the interface, the weaker the abstraction*.
- **L stops being a type-system concern and becomes a testing concern.** With no subclassing,
  substitutability can only be checked by running one suite against every implementation.
  Hence the L lab is a contract test.
- **S is a heuristic, not a rule.** Nobody agrees what "a responsibility" is. The usable version
  is the question in the table.
- **O is mostly a *result*.** Get I and D right and open/closed tends to fall out. Chasing it
  directly produces speculative abstraction.

### The counter-force (must stay in the material)

Every principle costs indirection. Indirection you do not need is **worse** than the duplication
it replaced, because it hides control flow. The deletion test is the discipline that keeps this
honest.

---

## 6. Repository layout & conventions

```
go-pattern-dojo/
  README.md             lead with "executable proof + deletion test", not the pattern list
  BLUEPRINT.md          curriculum + the 7-step method  (can be generated from §4 and §7)
  PROGRESS.md           the learner's running log — see §10
  design_pattern.md     THIS FILE (project context / handoff)
  go.mod                module pattern-dojo, go 1.26
  patterns/<id>-<name>/
    TASK.md             the 7-step card for this task
    THEORY.md           language-agnostic. NO GO CODE.
    attempt/            the learner writes here
    solution/           reference — sealed until they attempt
    PROOF.md            real terminal output + the measurements
    DONT.md             "wrong when ___" + deletion-test result
```

### Rules for `attempt/` vs `solution/` (the learner asked for this split explicitly)

- Both sides have **identical directory structure and identical test files**.
- In `attempt/`, the "before" code and the tests are **complete and runnable**. Only the
  "after" implementations are stubbed with `// TODO` bodies returning zero values.
- **Stubs must compile.** The tests must *fail*, not fail to compile. A compile error is a bad
  starting line; a red test is a good one.
- Therefore the correct starting state is: every `TestBefore_*` passes, every `TestAfter_*` fails.
- `go test ./...` from the repo root will report overall `FAIL` while `attempt/` is incomplete.
  That is expected. Scope runs with `go test ./patterns/<id>/attempt/...`.

### Commands used throughout

```bash
go test ./patterns/<id>/attempt/...     # the working loop
go test -race ./...                      # mandatory for every Block D task
gofmt -l . && go vet ./...               # must be clean before any task is "done"
go doc -src <symbol>                     # for the step-6 stdlib hunt
git diff --stat                          # the proof for OCP / facade / composition-root tasks
```

---

## 7. The curriculum — 39 tasks, 6 blocks

**Sizes:** **S** = 30–45 min · **M** = 60–90 min · **L** = 2–3 h, split over two sittings.
**Total: 40–60 hours, 12–14 weeks at ~4 h/week.**
If a task runs past 2× its size, the task was specified badly — stop and re-scope it rather
than letting the learner grind.

### Block A — Foundations (week 1, ~6 h)

| ID | Task | Size | Proof |
|---|---|---|---|
| A1 | **SOLID by measurement** — five code labs (fully specified in §8) | L | Five numbers; see §8 |
| A2 | Error taxonomy: sentinel · typed · wrapped | M | `errors.Is` / `errors.As` resolve through 2 wrap layers |
| A3 | Context: propagate, timeout, cancel across 3 call levels | M | `ctx.Err()` is `DeadlineExceeded` |
| A4 | **Embedding is not inheritance** | S | A test that breaks the learner's OO intuition |
| A5 | Same job as a `func` type and as a 1-method interface | S | Benchmark + a written "which would I ship?" |
| A6 | A type whose zero value is usable | S | `var c Counter` works with no constructor |

A4 is the task that kills the most Java habits. Do not let it be skipped.

### Block B — The daily five (weeks 2–3, ~10 h)

These five cover most real design decisions in Go services. Grind until reflexive.

| ID | Task | Size | Proof |
|---|---|---|---|
| B1 | **ISP:** split a 12-method store by what each consumer uses | M | Hand-written fake drops from 12 methods to 2 |
| B2 | **DIP:** move the interface into the consumer package | M | `go list -deps` shows the dependency arrow reversed |
| B3 | **Adapter + error translation** (driver errors → domain errors) | M | Service compiles with zero `database/sql` imports |
| B4 | **LSP as a test:** one contract suite, two implementations | M | Break the fake deliberately, watch the suite catch it |
| B5 | **Decorator:** middleware chain (log, auth, timing) | M | Execution-order test; then read the 3 stdlib decorators |
| B6 | **Strategy:** replace a growing `switch` with a func value | S | Add a 4th strategy without touching the caller |
| B7 | **Composition root:** wire handler → service → store in `main` | M | `grep` proves no package constructs its own dependencies |
| B8 | **Null object:** no-op logger / metrics recorder | S | Zero nil checks; benchmark shows the no-op is ~free |

### Block C — Construction (week 4, ~4 h)

| ID | Task | Size | Proof |
|---|---|---|---|
| C1 | `NewX` with validation — make invalid state unconstructable | S | Bad config returns an error, never a half-built struct |
| C2 | **Functional options** | M | Adding a 5th option leaves the *old* call site compiling unchanged |
| C3 | Options vs config struct — written decision for 3 real cases, no code | S | Written justification using the Uber threshold (§3.3) |
| C4 | **Registry** — `map[string]Constructor` + `init()` | M | Add a driver in a new file without editing the registry |

### Block D — Concurrency (weeks 5–7, ~16 h) — **the critical path for services**

| ID | Task | Size | Proof |
|---|---|---|---|
| D1 | **Goroutine-leak harness** — build once, reuse in D2–D9 | M | `runtime.NumGoroutine()` before/after with a settle period |
| D2 | Bounded worker pool | M | Atomic gauge: max in-flight never exceeds N |
| D3 | Pipeline + fan-out/fan-in | M | `-race` clean; every stage exits on cancel |
| D4 | `errgroup.WithContext` | M | First error cancels siblings (assert via per-task cancel flags) |
| D5 | **Graceful shutdown** of an HTTP server | L | In-flight request completes; new requests rejected |
| D6 | singleflight | M | 50 concurrent identical calls → exactly 1 downstream call |
| D7 | Rate limiter (token bucket) | M | Rate not exceeded over a measured window |
| D8 | **Circuit breaker** — closed / open / half-open | L | Table-driven transition test |
| D9 | Retry + backoff + jitter + **idempotency key** | L | No double-write after a retried ambiguous failure |

**D1 is non-negotiable and must come first.** Without a leak test, every later proof in this
block is a guess. Also: re-run the §3.1a race experiment here on the learner's own first draft
before "fixing" anything — finding your own race is the lesson.
`errgroup` and `singleflight` require `golang.org/x/sync` (not stdlib — confirmed).

### Block E — Structure (weeks 8–9, ~7 h)

| ID | Task | Size | Proof |
|---|---|---|---|
| E1 | Template method as a func parameter (`sort.Slice` shape) | S | Two algorithms share one skeleton |
| E2 | Chain of responsibility with early exit (auth rejects → chain stops) | M | Downstream handler's call count is 0 |
| E3 | **Facade + `internal/`** | M | Importing `internal/` from outside fails to compile |
| E4 | **Proxy:** caching proxy over a slow dependency | M | 3 calls → 1 downstream; then read `httputil.ReverseProxy` |
| E5 | **Explicit state machine** — `map[state]map[event]state`, not polymorphic state objects | M | Illegal transition rejected; all transitions table-tested |
| E6 | Observer: callbacks vs channels | M | **Leak test** on an unsubscribed observer |

**Cut on purpose:** Composite and Iterator. Worth *reading* (`io.MultiWriter`, range-over-func
in Go 1.23+) but not worth a build task for a services engineer. Re-add only if asked.

### Block F — Pro & architecture (weeks 10–14, ~12 h)

| ID | Task | Size | Proof |
|---|---|---|---|
| F1 | Generics: typed cache or store | M | Compile-time type safety, no `any` leakage |
| F2 | Fakes over mocks — delete a mocking library from a test | M | Contract suite reused for the fake |
| F3 | Table-driven tests + golden files | M | One table, 12 cases, readable failure output |
| F4 | **Hexagonal restructure of one real service of theirs** (capstone) | L | Domain package imports zero infrastructure |
| F5 | Package boundary audit script | M | CI fails if the domain package imports `database/sql` |
| F6 | Review drill on a real OSS Go repo | M | Written: one missing pattern, one over-applied |

### Architecture frameworks — decided ordering

1. **Hexagonal / Ports & Adapters** — yes, as F4. It *is* DIP + ISP at package scale. After
   B1, B2, B3 and B7 it is a half-hour restructuring exercise, not a new concept. **Teaching it
   first is the classic mistake:** you get the folder names (`domain/`, `ports/`, `adapters/`)
   without the judgment, and end up with ceremony.
2. **Clean Architecture** — same idea, more layers, more naming ceremony. Read a one-page
   summary; skip unless the learner's team already uses the vocabulary.
3. **DDD** — tactical bits (entity, value object, aggregate, repository) are cheap and useful.
   Strategic bits (bounded contexts, context mapping) only pay off across multiple teams.
   Optional Block G: G1 tactical DDD vocabulary applied to their own domain; G2 one
   bounded-context sketch. **Only after F4 goes smoothly.**
4. **CQRS / event sourcing** — do not touch until there is a real read/write asymmetry.

### Patterns deliberately demoted to "know the name, rarely write it in Go"

Abstract Factory, Prototype, Bridge, Mediator, Memento, Visitor (use a type switch),
Flyweight (use `sync.Pool`), Singleton (use DI; if unavoidable, `sync.OnceValue`),
Command (use a closure), Interpreter, Builder (use struct literals or functional options),
Template Method as inheritance (pass a func instead).

### The reviewer's lens — symptom → pattern

This table is the practical payoff. It should end up in the published README.

| Symptom in the code | Likely missing pattern |
|---|---|
| `sql.ErrNoRows` or driver error codes outside the DB layer | Adapter + error translation |
| A growing `switch` on a type or mode string | Strategy, or a registry map |
| Same logging/auth/timing copy-pasted across 8 handlers | Decorator / middleware |
| A package creating its own DB or HTTP client | DI + composition root |
| Constructor with 7 parameters, 5 optional | Functional options |
| Nil checks on a logger everywhere | Null object |
| N identical concurrent requests hitting the DB | singleflight or cache |
| Unbounded `go func()` inside a loop | Worker pool / bounded parallelism |
| Retry wrapped around a write with no key | Idempotency key (fix that *before* retrying) |
| **An interface with one implementation and one caller** | **Nothing. Delete the interface.** |

---

## 8. Task A1 — full specification (ready to build)

**Title:** SOLID by measurement · **Size:** L (~2.5 h, split: S+O+L, then I+D)
**Package path:** `patterns/a1-solid/` with five sub-packages: `srp ocp lsp isp dip`,
mirrored under both `attempt/` and `solution/`.

**Design intent:** the learner has read SOLID repeatedly and it stays hazy. Reading does not
work because the principles have no units. Every lab therefore produces **a number**.

### The measurable form (this table belongs in `THEORY.md` and `PROOF.md`)

| | What is measured | Violation | Fixed |
|---|---|---|---|
| S | I/O needed to test arithmetic | needs a temp dir | zero |
| O | files modified to add a case | 1 modified | 0 modified, 1 added |
| L | contract violations reported | ≥ 1 | 0 |
| I | methods a fake must implement | 10 | 1 |
| D | can the error path be simulated? | only by breaking a real file | trivially |

> If the learner cannot state the number, they have not found the violation.

### Lab specs

**srp — Single Responsibility**
- *Before:* an invoice type whose single method computes a total, renders a receipt string, and
  writes it to a file. Three reasons to change (tax rules, wording, storage) welded together.
- *After:* three separate things — a pure `Total()`, a pure `Render()`, and a `Save()`.
- *Proof:* the "before" test must call `t.TempDir()` **just to test arithmetic**; the "after"
  test for the same arithmetic has zero I/O, no cleanup, no error path.

**ocp — Open/Closed**
- *Before:* a pricing function with a `switch` over tiers (`regular`, `member`) plus an
  `ErrUnknownTier` sentinel. Adding a tier means editing tested code.
- *After:* an extension point as a **plain func type** (not an interface) plus a registry map,
  a `Register` function, and an `init()` that registers the existing tiers.
- *The exercise:* add a `student` tier (25% off) **in a brand-new file** that registers itself.
- *Proof:* `git diff --stat` shows **1 file added, 0 files modified**.

**lsp — Liskov Substitution**
- *Setup:* a 2-method `Store` interface (`Get`, `Put`) and an `ErrNotFound` sentinel that is
  **part of the contract**, not an implementation detail.
- *Two implementations:* one that keeps the promise; one "lying" store that returns
  `("", nil)` for a missing key. **Both satisfy the interface and both compile.**
- *The exercise:* write a `Contract(newStore func() Store) []string` that returns **one string
  per broken promise** (5 promises: missing→ErrNotFound; Put succeeds; Get returns what Put
  stored; Put overwrites; Get returns the overwritten value).
- *Design note:* returning a violation slice — rather than calling `t.Error` — lets the test
  assert that the lying store **fails**, without shipping a deliberately-red test.
- *Proof:* good store → 0 violations; lying store → ≥ 1. Include `var _ Store = ...` lines for
  both to make the point that **the compiler accepts both**.
- *This is the best lab in A1.* It shows that in Go, LSP is a testing discipline, not a type
  property.

**isp — Interface Segregation**
- *Before:* a 10-method `FatStore` interface; a `Welcome` function that calls exactly one of them.
- *After:* a 1-method interface declared next to the **consumer**, and `Welcome` written against it.
- *Proof:* in the test, the fake for the fat interface needs 10 methods (9 of them
  `panic("never called")`); the fake for the segregated one needs 1. Add `var _` lines showing
  the fat fake satisfies **both** — segregating costs the implementer nothing.

**dip — Dependency Inversion**
- *Before:* a service that reaches down and reads files itself (hardcoded mechanism). The arrow
  points down: business logic → infrastructure.
- *After:* the **consumer package owns** a small `Users` interface; the service takes it via
  constructor. The arrow now points up.
- *Proof:* the "before" test needs `t.TempDir()` + a real file to test string concatenation;
  the "after" test uses an in-memory map. **Additional proof:** the error path becomes trivially
  reachable, where before it could only be triggered by breaking a real file.

### Supporting files for A1

- `TASK.md` — the 7-step card, with a per-lab table of which file to edit, the exact proof
  commands, and a DONE checklist.
- `THEORY.md` — §5 of this document, written out language-agnostic, **no Go**.
- `PROOF.md` — a template: command, a blank fenced block for real output, and the measurable-form
  table with an empty "Mine" column.
- `DONT.md` — a 5-row table: principle · "wrong when…" · "deletion test: what broke?", left blank
  for the learner.

---

## 9. Mistakes made during this discussion (recorded so they are not repeated)

Honesty here is part of the method. Five corrections were made after the learner pushed back:

1. **Said concurrency was the critical path for services, then scheduled it at weeks 7–9**
   behind a six-week structural block. Fixed: concurrency moved up to Block D; the structural
   block moved back to E and was trimmed.
2. **The first task-card format silently dropped 3 of the 7 Dojo steps** (name the force, the
   language-agnostic shape, the deletion test). Dropping the shape step also broke the
   publish-as-language-agnostic goal. Fixed: the card now carries all seven.
3. **"45 minutes per task" was a lazy uniform estimate.** A 3-state circuit breaker with a
   table-driven transition test is not 45 minutes. Fixed: S/M/L sizing and an honest
   40–60 hour total (originally understated as ~30).
4. **A1 was first specified as a written exercise.** The learner correctly asked why SOLID
   wasn't being demonstrated in code given it kept going hazy. Fixed: five code labs with
   countable proofs (§8).
5. **Several blocks were delivered as prose instead of tables** after the learner had asked for
   tables. Fixed throughout.

Also worth recording: an earlier phase produced three scratch demo repos, and the learner got
lost in the file count. **Lesson: ship one task at a time, not a batch.**

---

## 10. `PROGRESS.md` — purpose and required shape

**Claude has no memory between conversations.** This file *is* the memory. The learner should
paste it as the first message of any new chat to resume.

Required sections:
1. **Context** — goal, Go version, method, and the learner's own "I won't peek" rule.
2. **Status table** — one row per block: Block · Tasks · Done · Notes.
3. **Task log** — one row per task: ID · Date · Size · Time taken · Proof passed? ·
   **What I got wrong** · My one-line "don't".
   *The "what I got wrong" column is the valuable one.* It is what makes the review drill sharp.
   It should not be left blank.
4. **Open questions for Claude** — so a new chat starts productive.
5. **Review drill log** — Date · File/repo reviewed · Missing pattern · Over-applied pattern.

The learner should also consider enabling memory in their Claude settings, but keep this file
regardless — it is also a published artifact.

---

## 11. Open items (unresolved at the time of writing)

1. **OCP proof rigour.** The `git diff --stat` proof for the `ocp` lab only works if the learner
   commits before starting. An open question was whether to add a `make verify` target that
   snapshots `after.go` and checks it was untouched automatically, versus leaving it manual.
   **The learner has not answered this.** Default to manual unless they say otherwise.
2. **Module path.** Currently `pattern-dojo`. Before publishing:
   `go mod edit -module github.com/<user>/go-pattern-dojo` and fix imports.
3. **Blocks A2 onward are specified at the one-line level only** (task + size + proof in §7).
   Each needs expansion into a full task card + before/after spec, the way A1 was in §8.
4. **Third-party library claims in §3.2 are unverified** for `grpc-go`, `chi`, `zap`,
   `sony/gobreaker`, `cenkalti/backoff`, `google/wire`. They were not downloaded. Verify APIs
   before writing any task that depends on them.

---

## 12. Agent instructions (read first)

You are continuing as the learner's mentor and pair. Operating rules:

### Do
- **Ship one task at a time.** Generate the full task folder for *one* ID, wait for the
  learner's attempt, review it, then move on. Do not batch tasks — they got lost once already.
- **Run everything you generate.** `gofmt -l .`, `go vet ./...`, `go test`. For Block D also
  `go test -race`. Paste the **real** output into `PROOF.md`. Never write expected output as if
  it were observed.
- **Keep `attempt/` stubs compiling** so tests fail rather than error (§6).
- **Honour the sealed-solution rule.** Do not reveal or hint at the `solution/` contents until
  the learner says they have attempted. If asked for a hint, give a question, not an answer.
- **Keep `THEORY.md` free of Go.** It is the publishable, language-agnostic layer.
- **State verified vs unverified.** If you could not run or fetch something, say so explicitly.
- **Push back on ceremony.** If the learner's solution contains an interface with one
  implementation and one caller, tell them to delete it. That is the job.
- Use tables. Keep prose tight.

### Don't
- Do not copy any code from Refactoring Guru (§3.5). Everything must be original.
- Do not reintroduce system design, DDIA, CQRS or strategic DDD. Out of scope by agreement.
- Do not re-derive the decisions in §3, §5 or §7 — they were settled with evidence.
- Do not pad task counts or invent pattern coverage. 39 tasks is the plan.

### Suggested next action
A1 is fully specified in §8 and can be generated immediately. Generate
`patterns/a1-solid/` (TASK.md, THEORY.md, PROOF.md, DONT.md, and the five labs under both
`attempt/` and `solution/`), verify the starting state is *before-tests-pass / after-tests-fail*,
and hand it over. Do not generate A2 until A1 is reviewed.

### Definition of done for any task
- [ ] `gofmt -l .` clean and `go vet ./...` passes
- [ ] The stated proof command produces the stated result, pasted real into `PROOF.md`
- [ ] `THEORY.md` has the learner's own words, no Go
- [ ] `DONT.md` has a "wrong when" line and a deletion-test result
- [ ] `PROGRESS.md` row filled in, including "what I got wrong"