//go:build ignore

// Command server serves order totals. build is its composition root.
package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/luckyparakh/goDP/patterns/b5-composition-root/after/handler"
	"github.com/luckyparakh/goDP/patterns/b5-composition-root/after/service"
	"github.com/luckyparakh/goDP/patterns/b5-composition-root/after/store"
)

// build reads config through getenv and wires store, service and handler
// together. It is the only place that calls their New functions.
func build(getenv func(string) string) (http.Handler, error) {
	st, err := store.New(getenv("ORDERS_FILE"))
	if err != nil {
		return nil, err
	}
	return handler.New(service.New(st)), nil
}

func main() {
	h, err := build(os.Getenv)
	if err != nil {
		log.Fatal(err)
	}
	srv := &http.Server{Addr: ":8080", Handler: h, ReadHeaderTimeout: 5 * time.Second}
	log.Fatal(srv.ListenAndServe())
}
