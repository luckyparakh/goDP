package srp

import (
	"os"
	"path/filepath"
	"testing"
)

var sample = Invoice{
	Items: []Item{
		{Name: "Widget", UnitCents: 1000, Qty: 2},
		{Name: "Gadget", UnitCents: 500, Qty: 1},
	},
	TaxPercent: 10,
}

const (
	wantTotal   int64 = 2750
	wantReceipt       = "Widget x2 20.00\nGadget x1 5.00\nTOTAL 27.50\n"
)

// BEFORE: to check arithmetic we are forced to hand over a directory.
func TestBefore_TotalNeedsADirectory(t *testing.T) {
	dir := t.TempDir()

	got, err := sample.Process(dir)
	if err != nil {
		t.Fatalf("Process: %v", err)
	}
	if got != wantTotal {
		t.Fatalf("total = %d, want %d", got, wantTotal)
	}
}

// BEFORE: the sum cannot be asked for on its own. A missing folder turns a
// plain sum into an error path, and this Process returns 0 for the sum.
func TestBefore_BrokenFolderTurnsASumIntoAnError(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "no-such-folder")

	got, err := sample.Process(missing)
	if err == nil {
		t.Fatal("expected an error when the folder does not exist")
	}
	if got != 0 {
		t.Fatalf("total = %d, want 0 (the sum is lost along with the save)", got)
	}
}

// AFTER: the same arithmetic, no directory, no cleanup, no error path.
func TestAfter_Total(t *testing.T) {
	if got := sample.Total(); got != wantTotal {
		t.Fatalf("Total() = %d, want %d", got, wantTotal)
	}
}

func TestAfter_Render(t *testing.T) {
	if got := Render(sample, wantTotal); got != wantReceipt {
		t.Fatalf("Render() = %q, want %q", got, wantReceipt)
	}
}

func TestAfter_Save(t *testing.T) {
	path := filepath.Join(t.TempDir(), "r.txt")

	if err := Save(path, "hello"); err != nil {
		t.Fatalf("Save: %v", err)
	}
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read back: %v", err)
	}
	if string(b) != "hello" {
		t.Fatalf("file = %q, want %q", b, "hello")
	}
}

func TestAfter_SaveReportsError(t *testing.T) {
	path := filepath.Join(t.TempDir(), "no-such-dir", "r.txt")

	if err := Save(path, "hello"); err == nil {
		t.Fatal("Save into a missing directory returned nil error")
	}
}

// The split must not change behaviour: same total, same bytes on disk.
func TestAfter_PartsMatchBefore(t *testing.T) {
	dir := t.TempDir()
	oldTotal, err := sample.Process(dir)
	if err != nil {
		t.Fatalf("Process: %v", err)
	}
	oldReceipt, err := os.ReadFile(filepath.Join(dir, "receipt.txt"))
	if err != nil {
		t.Fatalf("read receipt: %v", err)
	}

	newTotal := sample.Total()
	if newTotal != oldTotal {
		t.Errorf("Total() = %d, Process total = %d", newTotal, oldTotal)
	}
	if got := Render(sample, newTotal); got != string(oldReceipt) {
		t.Errorf("Render() = %q, Process wrote %q", got, oldReceipt)
	}
}
