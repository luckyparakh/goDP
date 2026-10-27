package dip

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

// BEFORE: to test string concatenation we must create a real file.
func TestBefore_GreetingNeedsARealFile(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "u1.txt"), []byte("Ada\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	got, err := NewFileService(dir).Greeting("u1")
	if err != nil || got != "Hello, Ada!" {
		t.Fatalf("Greeting = %q, %v", got, err)
	}
}

// BEFORE: the only error we can provoke is the one the OS gives us, and only
// by arranging the filesystem. We cannot ask for "storage is down".
func TestBefore_ErrorPathNeedsAMissingFile(t *testing.T) {
	if _, err := NewFileService(t.TempDir()).Greeting("nope"); err == nil {
		t.Fatal("expected an error for a missing file")
	}
}

var errBoom = errors.New("boom")

// memUsers is the whole test double: a map. No disk, no cleanup.
type memUsers map[string]string

func (m memUsers) Name(id string) (string, error) {
	n, ok := m[id]
	if !ok {
		return "", errBoom
	}
	return n, nil
}

var _ Users = memUsers(nil)

func TestAfter_Greeting(t *testing.T) {
	got, err := NewGreeter(memUsers{"u1": "Ada"}).Greeting("u1")
	if err != nil || got != "Hello, Ada!" {
		t.Fatalf("Greeting = %q, %v", got, err)
	}
}

// AFTER: the error path is one line to reach, and we choose the error.
func TestAfter_StorageErrorSurfaces(t *testing.T) {
	_, err := NewGreeter(memUsers{}).Greeting("u1")
	if !errors.Is(err, errBoom) {
		t.Fatalf("err = %v, want it to wrap errBoom", err)
	}
}
