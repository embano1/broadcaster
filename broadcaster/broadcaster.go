package broadcaster

import (
	"context"
	"fmt"
	"log"
	"sync"
)

type Broadcaster struct {
	inCh  chan interface{}
	mu    sync.RWMutex
	outCh []*chan interface{}
}

func New(in chan interface{}, out ...*chan interface{}) (*Broadcaster, error) {
	if out == nil {
		return nil, fmt.Errorf("at least one output channel must be specified")
	}

	return &Broadcaster{
		inCh:  in,
		outCh: out,
	}, nil
}

func (b *Broadcaster) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case v := <-b.inCh:
			b.mu.RLock()
			for i, c := range b.outCh {
				select {
				case *c <- v:
				default:
					log.Printf("channel %d blocked, skipping", i)
				}

			}
			b.mu.RUnlock()
		}
	}
}

func (b *Broadcaster) Add(c *chan interface{}) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.outCh = append(b.outCh, c)
}
