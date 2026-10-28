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

const Group string = "thanos-forwarder"

type Forwarder struct {
	Client *kgo.Client
}

func (f *Forwarder) Run(ctx context.Context) {
	for {
		fetches := f.Client.PollFetches(ctx)
		if err := fetches.Errors(); len(err) > 0 {
			log.Print(err) // fine for iteration 1
		}

		itr := fetches.RecordIter()
		for !itr.Done() {
			record := itr.Next()

			req, err := http.NewRequestWithContext(ctx, http.MethodPost, "http://receive-service:19291/api/v1/receive", bytes.NewReader(record.Value))
			if err != nil {
				log.Print(err)
			}
			// build custom header that matches Prometheus remote write 1.0 requests
			h := http.Header{}
			h.Add("Content-Type", "application/x-protobuf")
			h.Add("Content-Encoding", "snappy")
			h.Add("X-Prometheus-Remote-Write-Version", "0.1.0")
			h.Add("THANOS-TENANT", string(record.Key))
			req.Header = h

			client := http.Client{Timeout: time.Second * 5}
			resp, err := client.Do(req)
			if err != nil {
				log.Print("failed to send request: ", err)
			}
			defer resp.Body.Close()
			_, err = io.ReadAll(resp.Body)
			if err != nil {
				log.Print("failed to read response: ", err)
			}
			log.Print("YOU MADE IT")
		}
	}
}
