//go:build ignore

package isp

import "fmt"

// HOW IT WORKS
//
// Welcome declares the small interface it needs, with the one method it
// actually calls. Any type that has that method fits, including the big
// FatStore, so the people who provide the real thing change nothing.
//
// WHY IT HELPS
//
// A test stand-in needs one method instead of ten, and none of them is
// fake code that "never gets called". The interface also documents what
// Welcome really needs.
//
// WHEN NOT TO USE
//
// Do not split an interface just to have small ones. If one function is the
// only user and it really needs everything, the big interface is fine. Split
// when different callers use different parts, or when a stand-in is painful.

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
