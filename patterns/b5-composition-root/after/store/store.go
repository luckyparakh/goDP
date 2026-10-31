// Package store reads order totals from a JSON file. Your work: see the GOAL
// in ../../after.go.
package store

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
)

// File holds order totals loaded from a JSON object of id to cents.
type File struct{ cents map[string]int64 }

// New loads totals from the JSON file at path.
func New(path string) (*File, error) {
	// TODO: use path. This still reads the environment itself, the way
	// before/ does.
	return load(os.Getenv("ORDERS_FILE"))
}

func load(path string) (*File, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}
	f := &File{}
	if err := json.Unmarshal(b, &f.cents); err != nil {
		return nil, fmt.Errorf("store %s: %w", path, err)
	}
	return f, nil
}

// Cents returns the total for id.
func (f *File) Cents(_ context.Context, id string) (int64, error) {
	c, ok := f.cents[id]
	if !ok {
		return 0, fmt.Errorf("order %q not found", id)
	}
	return c, nil
}
