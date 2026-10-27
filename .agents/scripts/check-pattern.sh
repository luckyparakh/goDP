#!/usr/bin/env bash
# Mechanical checks for one pattern folder.
#
# Usage: .agents/scripts/check-pattern.sh <id> [start|done]
#   start  expect TestBefore_* to pass, TestAfter_* to fail, no compile errors
#   done   expect every test to pass
#
# Changes nothing in the repo. The sealed solution is checked in a throwaway
# module in the system temp folder. Needs git, go and a POSIX shell.
set -u

id="${1:-}"
mode="${2:-start}"
if [ -z "$id" ] || { [ "$mode" != start ] && [ "$mode" != done ]; }; then
  echo "usage: $0 <id> [start|done]"
  exit 2
fi

root="$(git rev-parse --show-toplevel)" || exit 2
cd "$root" || exit 2
p="patterns/$id"
s="solutions/$id"
fails=0

pass() { echo "PASS  $*"; }
fail() { echo "FAIL  $*"; fails=$((fails + 1)); }

[ -d "$p" ] || { echo "FAIL  $p does not exist"; exit 1; }

echo "== files =="
for f in THEORY.md before.go after.go; do
  if [ -f "$p/$f" ]; then pass "$p/$f exists"; else fail "$p/$f missing"; fi
done
if ls "$p"/*_test.go >/dev/null 2>&1; then pass "test file exists"; else fail "no *_test.go in $p"; fi
if [ -f "$s/after.go" ]; then pass "$s/after.go exists"; else fail "$s/after.go missing"; fi

echo "== format and vet =="
out="$(gofmt -l "$p" "$s" 2>&1)"
if [ -z "$out" ]; then pass "gofmt clean"; else fail "gofmt lists: $out"; fi
if out="$(go vet "./$p/" 2>&1)"; then pass "go vet clean"; else fail "go vet: $out"; fi

echo "== test state ($mode) =="
res="$(go test -count=1 -v "./$p/" 2>&1)"
if echo "$res" | grep -qE '\[build failed\]|\[setup failed\]'; then
  fail "tests do not compile"
  echo "$res" | head -20
else
  bp=$(echo "$res" | grep -cE '^--- PASS: TestBefore' || true)
  bf=$(echo "$res" | grep -cE '^--- FAIL: TestBefore' || true)
  ap=$(echo "$res" | grep -cE '^--- PASS: TestAfter' || true)
  af=$(echo "$res" | grep -cE '^--- FAIL: TestAfter' || true)
  echo "      TestBefore pass=$bp fail=$bf | TestAfter pass=$ap fail=$af"
  if [ "$mode" = start ]; then
    if [ "$bf" -eq 0 ] && [ "$bp" -ge 1 ] && [ "$ap" -eq 0 ] && [ "$af" -ge 1 ]; then
      pass "starting state: Before passes, After fails"
    else
      fail "not the starting state (want Before all pass, After all fail)"
    fi
  else
    if [ "$bf" -eq 0 ] && [ "$af" -eq 0 ] && [ "$ap" -ge 1 ]; then pass "every test passes"; else fail "some tests fail"; fi
  fi
fi

echo "== sealed solution =="
for f in "$s"/*.go; do
  [ -f "$f" ] || continue
  if [ "$(sed -n 1p "$f")" = "//go:build ignore" ] && [ -z "$(sed -n 2p "$f")" ]; then
    pass "$f starts with //go:build ignore and a blank line"
  else
    fail "$f must start with //go:build ignore and a blank line"
  fi
done
tmp="$(mktemp -d)"
log="$(mktemp)"
trap 'rm -rf "$tmp" "$log"' EXIT
gov="$(awk '/^go /{print $2}' go.mod)"
printf 'module solcheck\n\ngo %s\n' "$gov" > "$tmp/go.mod"
cp "$p"/*.go "$tmp/"
for f in "$s"/*.go; do
  [ -f "$f" ] && sed '1,2d' "$f" > "$tmp/$(basename "$f")"
done
if (cd "$tmp" && go vet ./... && go test -count=1 -race ./...) >"$log" 2>&1 \
  || { grep -q 'requires cgo' "$log" && (cd "$tmp" && go test -count=1 ./...) >"$log" 2>&1; }; then
  pass "solution passes (before.go + tests + solution files)"
else
  fail "solution does not pass"
  head -20 "$log"
fi

echo "== THEORY.md =="
t="$p/THEORY.md"
if [ -f "$t" ]; then
  q=$(grep -cE '^[0-9]+\. ' "$t" || true)
  g=$(grep -cE '^   - Guess:' "$t" || true)
  a=$(grep -cE '^   - Answer:' "$t" || true)
  if [ "$q" -eq 3 ] && [ "$g" -eq 3 ] && [ "$a" -eq 3 ]; then
    pass "3 questions, each with Guess and Answer lines"
  else
    fail "want 3 questions with Guess and Answer (questions=$q guess=$g answer=$a)"
  fi
  for h in '## Questions first' '## The idea' '## How the code shows it' '## Wrong when' '## In the wild' '## My delete-test result'; do
    if grep -qF "$h" "$t"; then pass "section: $h"; else fail "missing section: $h"; fi
  done
  for l in '- Standard library:' '- Source:' '- Book:' '**Try:**'; do
    if grep -qF -- "$l" "$t"; then pass "has line: $l"; else fail "missing line: $l"; fi
  done
  src="$(grep -E '^- Source:' "$t" | sed -E 's/.*\[(src\/[^]]+)\].*/\1/')"
  if [ -n "$src" ] && [ -f "$(go env GOROOT)/$src" ]; then
    pass "source path exists in your Go install: $src"
  else
    fail "source path not found under GOROOT: ${src:-none}"
  fi
fi

echo "== spoiler guard =="
if grep -qE 'HOW IT WORKS|WHEN NOT TO USE' "$p/after.go" "$p/before.go" "$p/THEORY.md" 2>/dev/null; then
  fail "solution headings found outside solutions/"
else
  pass "no solution headings outside solutions/"
fi

echo
if [ "$fails" -eq 0 ]; then echo "RESULT: PASS"; else echo "RESULT: FAIL ($fails)"; exit 1; fi
