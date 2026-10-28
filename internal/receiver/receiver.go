package receiver

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"sync/atomic"

	"github.com/twmb/franz-go/pkg/kgo"
)

// todo
type Receiver struct {
	writes atomic.Int64
	Client *kgo.Client
}

func (r *Receiver) Report(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "req: %d", r.writes.Load())
}

const (
	tenantHeader  string = "THANOS-TENANT"
	defaultTenant string = "default-tenant"
	topic         string = "remote-write"
)

// todo
func (r *Receiver) Produce(w http.ResponseWriter, req *http.Request) {
	tenant := req.Header.Get(tenantHeader)
	if tenant == "" {
		tenant = defaultTenant
	}

	buffer, err := io.ReadAll(req.Body)
	if err != nil {
		log.Printf("error: %v\n", err)
		w.WriteHeader(http.StatusUnprocessableEntity)
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
}
