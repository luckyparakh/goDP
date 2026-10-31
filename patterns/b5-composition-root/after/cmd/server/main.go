// Command server serves order totals. Your work: see the GOAL in
// ../../../after.go.
package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"time"
)

// build reads config through getenv and wires store, service and handler
// together. It is the only place in this task allowed to call their New
// functions, and the only place allowed to read the environment.
func build(getenv func(string) string) (http.Handler, error) {
	// TODO
	return nil, errors.New("TODO: wire store, service and handler here")
}

func main() {
	h, err := build(os.Getenv)
	if err != nil {
		log.Fatal(err)
	}
	srv := &http.Server{Addr: ":8080", Handler: h, ReadHeaderTimeout: 5 * time.Second}
	log.Fatal(srv.ListenAndServe())
}
