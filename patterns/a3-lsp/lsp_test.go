package lsp

import (
	"errors"
	"testing"
)

// BEFORE: both stores satisfy Store, and the lying one hides a trap.
func TestBefore_BothStoresSatisfyTheInterface(t *testing.T) {
	for name, newStore := range map[string]func() Store{
		"good":  NewGood,
		"lying": NewLying,
	} {
		s := newStore()
		if err := s.Put("k", "v"); err != nil {
			t.Errorf("%s: Put: %v", name, err)
		}
		if v, err := s.Get("k"); err != nil || v != "v" {
			t.Errorf("%s: Get = %q, %v; want v, nil", name, v, err)
		}
	}
}

func TestBefore_LyingStoreHidesMissingKeys(t *testing.T) {
	if _, err := NewGood().Get("nope"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("good store: err = %v, want ErrNotFound", err)
	}
	if v, err := NewLying().Get("nope"); v != "" || err != nil {
		t.Fatalf("lying store: Get = %q, %v; want \"\", nil", v, err)
	}
}

// AFTER: the contract suite.
func TestAfter_GoodStoreKeepsEveryPromise(t *testing.T) {
	calls := 0
	factory := func() Store {
		calls++
		return NewGood()
	}

	if v := Contract(factory); len(v) != 0 {
		t.Fatalf("good store reported violations: %q", v)
	}
	// A Contract that never touches the store proves nothing.
	if calls == 0 {
		t.Fatal("Contract never called newStore")
	}
}

func TestAfter_LyingStoreIsCaught(t *testing.T) {
	v := Contract(NewLying)
	t.Logf("violations: %q", v)
	if len(v) < 1 {
		t.Fatal("lying store passed the contract")
	}
}

// Mutants: each breaks exactly one promise by wrapping a good store.
type putFails struct{ Store } // promise 2

func (putFails) Put(string, string) error { return errors.New("disk full") }

type forgetful struct{ Store } // promise 3: accepts the write, drops it

func (forgetful) Put(string, string) error { return nil }

type writeOnce struct{ Store } // promise 4: refuses to overwrite

func (s writeOnce) Put(k, v string) error {
	if _, err := s.Store.Get(k); err == nil {
		return errors.New("already exists")
	}
	return s.Store.Put(k, v)
}

// promise 1 again, the subtle one: it does report an error for a missing key,
// but not ErrNotFound. Callers that ask errors.Is(err, ErrNotFound) get false.
type wrongError struct{ Store }

func (s wrongError) Get(k string) (string, error) {
	v, err := s.Store.Get(k)
	if errors.Is(err, ErrNotFound) {
		return "", errors.New("record absent")
	}
	return v, err
}

type keepsFirst struct{ Store } // promise 5: overwrite is silently ignored

func (s keepsFirst) Put(k, v string) error {
	if _, err := s.Store.Get(k); err == nil {
		return nil
	}
	return s.Store.Put(k, v)
}

func TestAfter_CatchesEachBrokenPromise(t *testing.T) {
	mutants := map[string]func() Store{
		"promise 1 (missing key)":        NewLying,
		"promise 1 (wrong error)":        func() Store { return wrongError{NewGood()} },
		"promise 2 (put fails)":          func() Store { return putFails{NewGood()} },
		"promise 3 (get after put)":      func() Store { return forgetful{NewGood()} },
		"promise 4 (overwrite succeeds)": func() Store { return writeOnce{NewGood()} },
		"promise 5 (newest value wins)":  func() Store { return keepsFirst{NewGood()} },
	}
	for name, newStore := range mutants {
		t.Run(name, func(t *testing.T) {
			if v := Contract(newStore); len(v) < 1 {
				t.Fatal("broken store passed the contract")
			}
		})
	}
}
