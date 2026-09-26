package outbox

import (
	"log"
	"time"
)

type Relay struct {
	db           *DB
	broker       Broker
	pollInterval time.Duration
	stopChan     chan struct{}
}

func NewRelay(db *DB, broker Broker, pollInterval time.Duration) *Relay {
	return &Relay{
		db:           db,
		broker:       broker,
		pollInterval: pollInterval,
		stopChan:     make(chan struct{}),
	}
}

func (r *Relay) Start() {
	go func() {
		ticker := time.NewTicker(r.pollInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				r.PollAndDispatch()
			case <-r.stopChan:
				return
			}
		}
	}()
}

func (r *Relay) Stop() {
	close(r.stopChan)
}

func (r *Relay) PollAndDispatch() int {
	pending := r.db.GetPendingOutbox()
	dispatched := 0
	for _, msg := range pending {
		err := r.broker.Publish(msg)
		if err == nil {
			err = r.db.MarkOutboxProcessed(msg.ID)
			if err != nil {
				log.Printf("failed to mark outbox msg %s as processed: %v\n", msg.ID, err)
			} else {
				dispatched++
			}
		} else {
			log.Printf("failed to publish outbox msg %s: %v\n", msg.ID, err)
		}
	}
	return dispatched
}
