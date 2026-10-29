package forwarder

import (
	"bytes"
	"context"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"
)

const (
	Group        string = "thanos-forwarder"
	endpoint     string = "http://receive-service:19291/api/v1/receive"
	tenantHeader string = "Thanos-Tenant"
)

type Forwarder struct {
	Client *kgo.Client
}

var headers = map[string]string{
	"Content-Type":                      "application/x-protobuf",
	"Content-Encoding":                  "snappy",
	"X-Prometheus-Remote-Write-Version": "0.1.0",
}

var timeouts = []time.Duration{
	time.Second,
	time.Second * 5,
	time.Second * 10,
	time.Second * 30,
}

// Run continually fetches records from Kafka and sends them to the metric backend.
func (f *Forwarder) Run(ctx context.Context) {
	hc := http.Client{Timeout: time.Second * 5}

	for {
		batch := f.Client.PollFetches(ctx)
		if err := batch.Errors(); len(err) > 0 {
			log.Print(err) // fine for iteration 1
		}

		itr := batch.RecordIter()
		for !itr.Done() {
			record := itr.Next()
			attempt := 0
			for {
				req, err := buildRequest(ctx, record)
				if err != nil {
					log.Fatal(err) // fine for now, for future find way to propagate error to main
				}
				status, err := sendRecord(&hc, req)

				if (status >= 100 && status < 200) || (status >= 300 && status < 400) { // stop the program
					log.Fatalf("status: %d", status)
				} else if status >= 200 && status < 300 {
					log.Print("record sent to thanos")
					break
				} else if (status >= 400 && status < 500) && status != 429 { // log and move on
					log.Printf("dropping record, status: %d", status)
					break
				}

				time.Sleep(timeouts[attempt])
				if attempt < len(timeouts)-1 {
					attempt++
				}
				log.Printf("status: %d; attempt: %d; error: %v", status, attempt, err)
			}
		}
		if err := f.Client.CommitUncommittedOffsets(ctx); err != nil {
			log.Fatal(err)
		} // after a batch, commit offset
	}
}

func buildRequest(ctx context.Context, record *kgo.Record) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(record.Value))
	if err != nil {
		return nil, err
	}
	// build custom header that matches Prometheus remote write 1.0 specification
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	req.Header.Set(tenantHeader, string(record.Key))

	return req, nil
}

// getResponse is a helper function that sends metrics to an endpoint over a POST request.
// It returns the status code, 0 means error.
func sendRecord(client *http.Client, req *http.Request) (int, error) {
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	_, err = io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	return resp.StatusCode, nil
}
