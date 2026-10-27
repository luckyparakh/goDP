package ctxflow

import "context"

// GOAL
//
// Three levels call each other: Handle calls Service, and Service calls Repo.
// Repo does the slow work: totalSteps steps, each taking stepDelay. Pass the
// context down through every level so that a timeout or a cancel set at the
// top reaches the bottom and stops the work.
//
//   - Repo checks the context before and during each step, and stops as soon
//     as it is done. It returns how many steps finished, and the context's error.
//   - Service and Handle pass ctx down and return what they get back, adding
//     context to an error without hiding it.
//   - If the context is already cancelled when Repo starts, no work is done.
//   - If every step finishes, the error is nil and the count is totalSteps.
//
// You are done when every test in this folder passes:
//
//	go test ./patterns/a7-context/

// Handle is the top level. It calls Service.
func Handle(ctx context.Context) (int, error) {
	// TODO
	return 0, nil
}

// Service is the middle level. It calls Repo.
func Service(ctx context.Context) (int, error) {
	// TODO
	return 0, nil
}

// Repo is the bottom level. It does the slow work.
func Repo(ctx context.Context) (int, error) {
	// TODO
	return 0, nil
}
