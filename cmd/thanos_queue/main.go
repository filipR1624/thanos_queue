package main

import (
	"context"
	"log"
	"net/http"

	"github.com/filipR1624/thanos_queue/internal/forwarder"
	"github.com/filipR1624/thanos_queue/internal/receiver"
	"github.com/twmb/franz-go/pkg/kgo"
)

func main() {
	seeds := []string{"redpanda-service:9092"} // hardcoded for now

	// creating receiver client
	rec, err := kgo.NewClient(kgo.SeedBrokers(seeds...))
	if err != nil {
		log.Fatal(err)
	}
	defer rec.Close()
	r := receiver.Receiver{
		Client: rec,
	}

	// creating forwarder client
	fwd, err := kgo.NewClient(
		kgo.SeedBrokers(seeds...),
		kgo.ConsumerGroup(forwarder.Group),
		kgo.ConsumeTopics("remote-write"),
		kgo.DisableAutoCommit(),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer fwd.Close()
	f := forwarder.Forwarder{
		Client: fwd,
	}
	go f.Run(context.Background())

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/write", r.Produce)
	mux.HandleFunc("/", r.Report)

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
