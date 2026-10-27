# One job per piece

**TL;DR** — A function that totals a bill, writes a receipt and saves it to disk
has three different reasons to change, and you cannot ask it for the total
without also asking it to write a file. Split it by reason to change, and the
arithmetic becomes testable with no folder, no cleanup and no error path.

## The everyday version

Picture a single recipe card that also carries your shopping list and the
washing-up rota. Change the shopping list and you reprint the whole card,
recipe included. Want to check how much salt the recipe asks for? You still end
up reading the washing-up rota to get there. Three different people keep asking
you to change one card.

## How you'd fix it in real life

Use three cards: the recipe, the shopping list, the rota. Now the shopping list
changes without touching the recipe, and checking the salt means reading one
short card.

> Where this analogy breaks: reprinting a card is slow but safe. Editing a
> shared function is fast and can silently break every other caller, which a
> recipe card can't do.

## The same pain, in code

Here is the version we start from. `Process` does three jobs in a row:

- it adds up the lines and applies tax, so **tax rules** live here
- it builds the receipt text, so **wording** lives here
- it writes that text to a file, so **storage** lives here

```go
func (inv Invoice) Process(dir string) (int64, error) {
	var subtotal int64
	for _, it := range inv.Items {
		subtotal += it.UnitCents * it.Qty
	}
	total := subtotal + subtotal*inv.TaxPercent/100

	var b strings.Builder
	for _, it := range inv.Items {
		fmt.Fprintf(&b, "%s x%d %s\n", it.Name, it.Qty, money(it.UnitCents*it.Qty))
	}
	fmt.Fprintf(&b, "TOTAL %s\n", money(total))

	if err := os.WriteFile(filepath.Join(dir, "receipt.txt"), []byte(b.String()), 0o644); err != nil {
		return 0, fmt.Errorf("write receipt: %w", err)
	}
	return total, nil
}
```

Look at the signature: `Process(dir string)`. The arithmetic is in there, but
you cannot reach it without handing over a folder to write into. So a test for
a sum looks like this:

```go
func TestBefore_TotalNeedsADirectory(t *testing.T) {
	dir := t.TempDir()

	got, err := sample.Process(dir)
	...
}
```

That is the cost, and it is countable: **one temporary folder to check one
sum**. The knock-on effect is worse. Point it at a folder that doesn't exist
and a plain piece of arithmetic turns into an error path:

```go
got, err := sample.Process(missing)
// err is non-nil, got is 0 — the sum is gone
```

(This `Process` chose to return `0` alongside the error. It could have returned
the total as well. Even then, the only way to reach the total is to go through
the save.)

## The fix

Cut it into three pieces, one per reason to change. `Total` only does
arithmetic. `Render` only builds text, and it takes the total as an argument
rather than recalculating it, so the wording never depends on the tax rule.
`Save` is the only piece allowed to touch the disk.

```go
func (inv Invoice) Total() int64 {
	var subtotal int64
	for _, it := range inv.Items {
		subtotal += it.UnitCents * it.Qty
	}
	return subtotal + subtotal*inv.TaxPercent/100
}

func Render(inv Invoice, total int64) string {
	var b strings.Builder
	for _, it := range inv.Items {
		fmt.Fprintf(&b, "%s x%d %s\n", it.Name, it.Qty, money(it.UnitCents*it.Qty))
	}
	fmt.Fprintf(&b, "TOTAL %s\n", money(total))
	return b.String()
}

func Save(path, receipt string) error {
	if err := os.WriteFile(path, []byte(receipt), 0o644); err != nil {
		return fmt.Errorf("save receipt: %w", err)
	}
	return nil
}
```

A tax change now edits `Total` and nothing else. And the sum is testable on its
own:

```go
func TestAfter_Total(t *testing.T) {
	if got := sample.Total(); got != wantTotal {
		t.Fatalf("Total() = %d, want %d", got, wantTotal)
	}
}
```

No folder. No cleanup. No error to handle.

## Proof

Starting from the stubbed version, `go test` reports:

```
--- FAIL: TestAfter_Total (0.00s)
    srp_test.go:52: Total() = 0, want 2750
--- FAIL: TestAfter_Render (0.00s)
    srp_test.go:58: Render() = "", want "Widget x2 20.00\nGadget x1 5.00\nTOTAL 27.50\n"
--- FAIL: TestAfter_Save (0.00s)
    srp_test.go:70: read back: open ...\TestAfter_Save2833447461\001\r.txt: The system cannot find the file specified.
--- FAIL: TestAfter_SaveReportsError (0.00s)
    srp_test.go:81: Save into a missing directory returned nil error
--- FAIL: TestAfter_PartsMatchBefore (0.01s)
    srp_test.go:99: Total() = 0, Process total = 2750
FAIL
```

After the split, the same command:

```
--- PASS: TestBefore_TotalNeedsADirectory (0.01s)
--- PASS: TestBefore_BrokenFolderTurnsASumIntoAnError (0.00s)
--- PASS: TestAfter_Total (0.00s)
--- PASS: TestAfter_Render (0.00s)
--- PASS: TestAfter_Save (0.01s)
--- PASS: TestAfter_SaveReportsError (0.00s)
--- PASS: TestAfter_PartsMatchBefore (0.01s)
ok  	github.com/luckyparakh/goDP/patterns/a1-srp	0.191s
```

Note `TestAfter_PartsMatchBefore`. The split has to produce the same total and
the same bytes as the original, otherwise it isn't a refactor.

Both runs are on Go 1.26, Windows.

## What I got wrong

> TODO-yours. From PROGRESS.md: *"Permission of as os.ModeAppend but it should
> be 0o644 and it did't use %w for error wrap"*.
>
> Write 3 or 4 sentences: what you expected `os.ModeAppend` to do, what actually
> happened on your machine, and how you found out. Worth mentioning that the
> tests still passed, because each one writes to a fresh temp folder — which is
> exactly the kind of bug a green test suite hides.

## The delete test

Put `Process` back and try to check one sum without creating a folder.

> "I can't get sum w/o folder and a bad folder will give error."

That sentence is the justification for the split. If nothing had got harder,
the three pieces would have been decoration and the original was fine.

## Do and don't

| Do | Don't |
|---|---|
| Split when you can name two different reasons to change | Split because a function looks long |
| Keep the pure parts pure: no files, no printing | Let the formatting step recalculate the total |
| Pass the result in, like `Render(inv, total)` | Chain the pieces so one can only run after another |
| Let the caller decide what to do when saving fails | Throw away a correct result because a later step failed |

## When not to use this

If the three parts always change together and are never used apart, one
function is easier to read than three. The split also moves work to the caller:
someone now has to call `Total`, `Render` and `Save` in order and decide what
happens when one of them fails. Get that wrong and you have rebuilt the original
problem with extra steps.

> TODO-yours. Close with a line or two in your own words: would you reach for
> this split in your own service code, and where would you stop?

## Further reading

- `log/slog.Handler` — a logger has two jobs, making a record and deciding how
  it looks and where it goes. `slog` hands the second to a `Handler`, so the
  same `logger.Info(...)` call produces text or JSON depending only on which
  handler you built it with. Run `go doc log/slog.Handler`.
- Source: [src/log/slog/handler.go](https://github.com/golang/go/blob/master/src/log/slog/handler.go)
- The code, tests and task card for this post:
  [`patterns/a1-srp`](../patterns/a1-srp)
