package forwarder

import (
	"context"
	"log"

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
			log.Printf("%v; %d; %d", string(record.Key), record.Partition, record.Offset)
		}
	}
}
