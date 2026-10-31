package errs

import "fmt"

// GOAL
//
// Three layers call each other: Handle calls Load, and Load calls Get. Make
// every error keep its identity through both wraps.
//
//   - Get returns ErrNotFound for an unknown id, and a *ValidationError
//     (Field "id") for an empty id. It does not wrap anything.
//   - Load wraps whatever Get returns, adding context, so errors.Is and
//     errors.As still work on the result.
//   - Handle does the same with whatever Load returns.
//
// The known order is "42", which returns "order 42" with a nil error.
//
// You are done when every test in this folder passes:
//
//	go test ./patterns/a6-errors/

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
