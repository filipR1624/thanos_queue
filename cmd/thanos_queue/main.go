package main

import (
	"log"
	"net/http"

	"github.com/filipR1624/thanos_queue/internal/receiver"
)

func main() {
	t := receiver.Receiver{}
	mux := http.NewServeMux()
	mux.Handle("/api/v1/write", t.Receive())
	mux.Handle("/", t.Report())

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
