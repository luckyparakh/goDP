//go:build ignore

package errs

import "fmt"

// HOW IT WORKS
//
// Get creates the error and wraps nothing: a sentinel for a known condition,
// and a typed error when the caller needs data (which field). Load and Handle
// add context with %w. The %w verb keeps the original error inside the new
// one, so errors.Is and errors.As walk down the chain and still find it.
//
// WHY IT HELPS
//
// A caller can ask "is this a not-found?" or "which field was invalid?" no
// matter how many layers added their own words, and nobody has to search
// message text. Each layer says where it was, so the final message reads like
// a path: handle: load "nope": not found.
//
// WHEN NOT TO USE
//
// %w makes the inner error part of your API, because callers can now depend
// on it. If you want to hide an implementation detail, wrap with %v on
// purpose. Also, wrap only when you add something useful. Wrapping at every
// line makes messages long, and logging an error and also returning it
// reports it twice.

// Get looks the order up. It is the bottom layer.
func Get(id string) (string, error) {
	if id == "" {
		return "", &ValidationError{Field: "id", Reason: "must not be empty"}
	}
	v, ok := orders[id]
	if !ok {
		return "", ErrNotFound
	}
	return v, nil
}

// Load calls Get and adds context to any error.
func Load(id string) (string, error) {
	v, err := Get(id)
	if err != nil {
		return "", fmt.Errorf("load %q: %w", id, err)
	}
	return v, nil
}

// Handle calls Load and adds context to any error. It is the top layer.
func Handle(id string) (string, error) {
	v, err := Load(id)
	if err != nil {
		return "", fmt.Errorf("handle: %w", err)
	}
	return v, nil
}
