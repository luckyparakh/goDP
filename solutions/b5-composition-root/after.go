//go:build ignore

package compositionroot

// HOW IT WORKS
//
// Every constructor takes what it needs as an argument and keeps it.
// handler declares the Service it needs, service declares the Store it needs
// (B1), and store takes a file path instead of reading the environment.
//
// build in main is the composition root. It is the one function that reads
// config and calls store.New, service.New and handler.New, in that order. The
// whole object graph of the program is those few lines.
//
// The solution files are after/store/store.go, after/service/service.go,
// after/handler/handler.go and after/cmd/server/main.go.
//
// WHY IT HELPS
//
// A test hands the handler a fake service and needs no file and no
// environment. Swapping the store is a one-line change in main. Config is read
// in one place, so a missing setting fails at startup with one clear error,
// not deep inside a request.
//
// build takes getenv as a parameter, so main_test.go can test the real wiring
// without touching the process environment.
//
// WHEN NOT TO USE
//
// A tiny program with one implementation of everything and no tests that need
// fakes gains little. And keep main boring: reading config and calling
// constructors in order. Once branches and business rules appear in main, move
// them into a package. You do not need a dependency-injection framework for
// what ten lines of main can do.
