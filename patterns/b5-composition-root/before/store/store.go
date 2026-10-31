// Package store reads order totals from a JSON file. In this version it also
// decides which file, by reading an environment variable itself.
package store

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

// File holds order totals loaded from a JSON object of id to cents.
type File struct{ cents map[string]int64 }

// New loads the file named by the ORDERS_FILE environment variable.
func New() (*File, error) {
	path := os.Getenv("ORDERS_FILE")
	if path == "" {
		return nil, errors.New("store: ORDERS_FILE is not set")
	}
	return load(path)
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
