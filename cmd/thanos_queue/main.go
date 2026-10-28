package main

import (
	"log"
	"net/http"

	"github.com/filipR1624/thanos_queue/internal/receiver"
	"github.com/twmb/franz-go/pkg/kgo"
)

func main() {
	seeds := []string{"redpanda-service:9092"} // hardcoded for now

	cl, err := kgo.NewClient(
		kgo.SeedBrokers(seeds...),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer cl.Close()

	r := receiver.Receiver{
		Client: cl,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/write", r.Produce)
	mux.HandleFunc("/", r.Report)

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
