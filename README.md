# goDP: design patterns in Go, proven by tests

Learn each pattern by running a test that fails before the fix and passes after
it. Then delete the pattern and write down what broke. If nothing broke, it was
decoration.

- **Start here:** [HOWTO.md](HOWTO.md). It is the same six steps for every pattern.
- **The task list:** [BLUEPRINT.md](BLUEPRINT.md).
- **Your log:** [PROGRESS.md](PROGRESS.md).
- **AI helpers (rules, review, QA):** [.agents/README.md](.agents/README.md).

Each pattern is one folder in `patterns/`: a short `THEORY.md`, `before.go`
(the problem), `after.go` (you fix it) and a test. The finished code is in
`solutions/`. Open it only after you have tried.

## Status

Block A (A1 to A10) and Block B (B1 to B6) are ready. Everything else is listed
in the blueprint and not written yet.

```bash
go test ./patterns/a1-srp/     # start here
```

Requires Go 1.26 or newer.

## Licence

Code is Apache-2.0 (see [LICENSE](LICENSE)). A separate licence for the
Markdown has not been chosen yet.
