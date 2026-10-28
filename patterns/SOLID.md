# SOLID: the view across all five

Optional background. Each of [A1](a1-srp/THEORY.md) to [A5](a5-dip/THEORY.md)
names and explains its own principle, so you do not need this page to do the
tasks. Read it when you want the part that only makes sense across all five,
either before A1 or after A5.

## One question instead of five rules

SOLID is usually taught as five rules with no units, so nobody can check whether
a piece of code obeys them. Start from one measurable question instead:

> **When requirement X changes, how much code must I open, re-read, re-test and risk?**

Call that the **blast radius**. All five principles are tactics for shrinking it.

Shrinking blast radius means being able to swap parts, so a second question falls
out of the first:

> **Can I replace a part without the rest of the system noticing?**

## The five principles as questions to ask code

| | Principle | The force it fights | Question to ask the code | Task |
|---|---|---|---|---|
| **S** | Single Responsibility | A module answers to two masters, so either can break it | "How many *different* requirements would make me edit this?" | [A1](a1-srp/THEORY.md) |
| **O** | Open/Closed | New cases force edits to working, tested code | "To add a case, do I *add* or do I *modify*?" | [A2](a2-ocp/THEORY.md) |
| **L** | Liskov Substitution | A replacement is accepted by the compiler but misbehaves | "Does every implementation keep the same promises?" | [A3](a3-lsp/THEORY.md) |
| **I** | Interface Segregation | A consumer is coupled to operations it never uses | "Do I depend on more than I use?" | [A4](a4-isp/THEORY.md) |
| **D** | Dependency Inversion | Policy reaches down and grabs a mechanism | "Who owns the interface: the user or the provider?" | [A5](a5-dip/THEORY.md) |

## How much each principle matters, depending on the language

This is a judgement, not a law. It depends on how a language lets types satisfy
interfaces.

- **I and D carry the most weight** where a type satisfies an interface
  implicitly (Go is one such language). The *consumer* can declare the small
  interface it needs, so segregation and inversion cost almost nothing. This is
  the same idea as the proverb "the bigger the interface, the weaker the
  abstraction".
- **L stops being a type-system concern and becomes a testing concern** where
  there is no subclassing. A compiler can check that signatures match; it cannot
  check that behaviour matches. The only way to check substitutability is to run
  one suite of promises against every implementation.
- **S is a heuristic, not a rule.** Nobody agrees what "a responsibility" is. The
  usable version is the question in the table above.
- **O is mostly a result.** If I and D are right, open/closed tends to fall out.
  Chasing it directly produces speculative abstraction.

## The counter-force

Every principle costs indirection. Indirection you do not need is **worse** than
the duplication it replaced, because it hides control flow. The delete test at
the end of each task is the discipline that keeps this honest: remove the
abstraction and name the specific thing that stops working. If you cannot, it
was decoration.
