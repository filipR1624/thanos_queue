package receiver

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"sync/atomic"

	"github.com/twmb/franz-go/pkg/kgo"
)

// Receiver takes remote-write requests and hands them to Kafka.
type Receiver struct {
	writes atomic.Int64
	Client *kgo.Client
}

// Report writes back the number of successful writes to Kafka.
func (r *Receiver) Report(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Writes: %d", r.writes.Load())
}

const (
	tenantHeader  string = "THANOS-TENANT" // matches Prometheus's tenant header
	defaultTenant string = "default-tenant"
	topic         string = "remote-write" // hardcoded for now, users must ensure the topic exists
)

// Produce receives Prometheus remote-write requests and produces a record for Kafka.
// Keys are tenants to achieve separation (one topic per tenant, data does not mix).
// Returns OK on successful writes, otherwise 5xx for Prometheus to retry.
func (r *Receiver) Produce(w http.ResponseWriter, req *http.Request) {
	tenant := req.Header.Get(tenantHeader)
	if tenant == "" {
		tenant = defaultTenant
	}

	buffer, err := io.ReadAll(req.Body)
	if err != nil {
		log.Printf("error: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	record := &kgo.Record{
		Key:   []byte(tenant),
		Value: buffer,
		Topic: topic,
	}

	result := r.Client.ProduceSync(req.Context(), record)
	if err := result.FirstErr(); err != nil {
		log.Printf("error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	r.writes.Add(1)
	w.WriteHeader(http.StatusOK)
	log.Printf("write received")
}
