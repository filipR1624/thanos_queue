package receiver

import (
	"fmt"
	"net/http"
	"sync"
)

// Handler tracks number of requests to a endpoint and offers synchronization through a mutex.
type Receiver struct {
	counter int
	m       sync.Mutex
}

// Receive is a closure that tracks how many requests have been recorded.
func (t *Receiver) Receive() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		t.m.Lock()
		t.counter++
		t.m.Unlock()
	}
	return http.HandlerFunc(fn)
}

// Write is a closure that reports the number of requests that have occured.
func (t *Receiver) Report() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		t.m.Lock()
		c := t.counter
		t.m.Unlock()
		fmt.Fprintf(w, "req: %d", c)

	}
	return http.HandlerFunc(fn)
}
