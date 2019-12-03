package broadcaster

import (
	"context"
	"log"
	"sync"

	"github.com/google/uuid"
)

// Broadcaster notifies zero or multiple subscribers on every message received
// on the configured input channel
type Broadcaster interface {
	Run(ctx context.Context)
	Subscribe(c *chan interface{}) string
	Unsubscribe(id string)
}

type broadcaster struct {
	inCh        chan interface{}
	mu          sync.RWMutex
	subscribers map[string]*chan interface{}
}

// New returns an initialized Broadcaster with the input channel set to "input"
func New(input chan interface{}) Broadcaster {
	return &broadcaster{
		inCh:        input,
		subscribers: make(map[string]*chan interface{}),
	}
}

// Run starts the bradocaster and blocks
func (b *broadcaster) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case v := <-b.inCh:
			b.mu.RLock()
			for sub, c := range b.subscribers {
				select {
				case *c <- v:
				default:
					log.Printf("subscriber %s blocked, skipping", sub)
				}
			}
			b.mu.RUnlock()
		}
	}
}

// Subscribe adds a subscriber channel to the broadcaster and returns an ID used
// to unsubscribe from the broadcaster. Before closing the given channel
// subscribers must call Unsubscribe() to avoid panics
func (b *broadcaster) Subscribe(c *chan interface{}) string {
	b.mu.Lock()
	defer b.mu.Unlock()
	u := uuid.New().String()
	b.subscribers[u] = c
	return u
}

// Unsubscribe removes the given subscriber ID from the broadcaster
func (b *broadcaster) Unsubscribe(id string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if _, ok := b.subscribers[id]; ok {
		delete(b.subscribers, id)
	}
}
