package isp

import "fmt"

// GOAL
//
// Welcome does the same job as WelcomeFat, but depends only on what it uses.
// Declare UserGetter here, next to the code that uses it, with exactly the
// methods Welcome calls. Behaviour and error wrapping must match WelcomeFat.
//
// You are done when every test in this folder passes:
//
//	go test ./patterns/a4-isp/

// UserGetter is the small dependency of Welcome.
type UserGetter interface {
	GetUser(id string) (User, error)
}

// Welcome greets the user with the given id.
func Welcome(g UserGetter, id string) (string, error) {
	u, err := g.GetUser(id)
	if err != nil {
		return "", fmt.Errorf("welcome %q: %w", id, err)
	}
	return "Welcome, " + u.Name + "!", nil
}
