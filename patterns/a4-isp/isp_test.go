package isp

import (
	"errors"
	"reflect"
	"testing"
)

var errBoom = errors.New("boom")

// fakeFat is what testing against FatStore costs: ten methods, one of which
// does anything.
type fakeFat struct{ users map[string]User }

func (f fakeFat) GetUser(id string) (User, error) {
	u, ok := f.users[id]
	if !ok {
		return User{}, errBoom
	}
	return u, nil
}
func (fakeFat) CreateUser(User) error              { panic("never called") }
func (fakeFat) UpdateUser(User) error              { panic("never called") }
func (fakeFat) DeleteUser(string) error            { panic("never called") }
func (fakeFat) ListUsers() ([]User, error)         { panic("never called") }
func (fakeFat) GetOrder(string) (Order, error)     { panic("never called") }
func (fakeFat) CreateOrder(Order) error            { panic("never called") }
func (fakeFat) DeleteOrder(string) error           { panic("never called") }
func (fakeFat) ListOrders(string) ([]Order, error) { panic("never called") }
func (fakeFat) Ping() error                        { panic("never called") }

// fakeGetter is what testing against the segregated interface costs.
type fakeGetter struct{ users map[string]User }

func (f fakeGetter) GetUser(id string) (User, error) {
	u, ok := f.users[id]
	if !ok {
		return User{}, errBoom
	}
	return u, nil
}

// Segregating costs the implementer nothing: the fat fake satisfies both.
var (
	_ FatStore   = fakeFat{}
	_ UserGetter = fakeFat{}
	_ UserGetter = fakeGetter{}
)

var users = map[string]User{"u1": {ID: "u1", Name: "Ada"}}

func TestBefore_FatStoreHasTenMethods(t *testing.T) {
	if n := reflect.TypeFor[FatStore]().NumMethod(); n != 10 {
		t.Fatalf("FatStore has %d methods, want 10", n)
	}
}

func TestBefore_WelcomeFat(t *testing.T) {
	got, err := WelcomeFat(fakeFat{users}, "u1")
	if err != nil || got != "Welcome, Ada!" {
		t.Fatalf("WelcomeFat = %q, %v", got, err)
	}
	if _, err := WelcomeFat(fakeFat{users}, "nope"); !errors.Is(err, errBoom) {
		t.Fatalf("err = %v, want errBoom", err)
	}
}

// AFTER: the number to beat is 10.
func TestAfter_InterfaceHasOneMethod(t *testing.T) {
	if n := reflect.TypeFor[UserGetter]().NumMethod(); n != 1 {
		t.Fatalf("UserGetter has %d methods, want 1", n)
	}
}

func TestAfter_Welcome(t *testing.T) {
	got, err := Welcome(fakeGetter{users}, "u1")
	if err != nil || got != "Welcome, Ada!" {
		t.Fatalf("Welcome = %q, %v", got, err)
	}
}

func TestAfter_WelcomeWrapsError(t *testing.T) {
	_, err := Welcome(fakeGetter{users}, "nope")
	if !errors.Is(err, errBoom) {
		t.Fatalf("err = %v, want it to wrap errBoom", err)
	}
}
