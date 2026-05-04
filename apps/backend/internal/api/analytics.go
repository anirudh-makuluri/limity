package api

import (
	"context"
	"log"
	"sync"
	"time"
)

type RequestEvent struct {
	Timestamp      time.Time
	Method         string
	Route          string
	Path           string
	StatusCode     int
	DurationMs     int64
	ClientIP       string
	UserAgent      string
	OwnerUserID    string
	APIKey         string
}

type AnalyticsStore interface {
	InsertRequestEvents(ctx context.Context, events []RequestEvent) error
	ResolveOwnerUserIDByAPIKey(ctx context.Context, apiKey string) (string, error)
}

type AsyncAnalytics struct {
	store         AnalyticsStore
	events        chan RequestEvent
	batchSize     int
	flushInterval time.Duration

	stopCh chan struct{}
	wg     sync.WaitGroup

	ownerByAPIKey map[string]string
}

func NewAsyncAnalytics(store AnalyticsStore, queueSize, batchSize int, flushInterval time.Duration) *AsyncAnalytics {
	if queueSize <= 0 {
		queueSize = 10000
	}
	if batchSize <= 0 {
		batchSize = 200
	}
	if flushInterval <= 0 {
		flushInterval = time.Second
	}

	a := &AsyncAnalytics{
		store:         store,
		events:        make(chan RequestEvent, queueSize),
		batchSize:     batchSize,
		flushInterval: flushInterval,
		stopCh:        make(chan struct{}),
		ownerByAPIKey: make(map[string]string),
	}
	a.wg.Add(1)
	go a.run()
	return a
}

func (a *AsyncAnalytics) Enqueue(event RequestEvent) bool {
	select {
	case a.events <- event:
		return true
	default:
		return false
	}
}

func (a *AsyncAnalytics) Close() {
	close(a.stopCh)
	a.wg.Wait()
}

func (a *AsyncAnalytics) run() {
	defer a.wg.Done()
	ticker := time.NewTicker(a.flushInterval)
	defer ticker.Stop()

	batch := make([]RequestEvent, 0, a.batchSize)
	flush := func() {
		if len(batch) == 0 {
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		a.enrichOwnerUserIDs(ctx, batch)
		err := a.store.InsertRequestEvents(ctx, batch)
		cancel()
		if err != nil {
			log.Printf("analytics flush failed: %v", err)
		}
		batch = batch[:0]
	}

	for {
		select {
		case <-a.stopCh:
			flush()
			return
		case <-ticker.C:
			flush()
		case event := <-a.events:
			batch = append(batch, event)
			if len(batch) >= a.batchSize {
				flush()
			}
		}
	}
}

func (a *AsyncAnalytics) enrichOwnerUserIDs(ctx context.Context, batch []RequestEvent) {
	for i := range batch {
		if batch[i].OwnerUserID != "" || batch[i].APIKey == "" {
			continue
		}

		if cached, ok := a.ownerByAPIKey[batch[i].APIKey]; ok {
			batch[i].OwnerUserID = cached
			continue
		}

		ownerUserID, err := a.store.ResolveOwnerUserIDByAPIKey(ctx, batch[i].APIKey)
		if err != nil || ownerUserID == "" {
			continue
		}
		a.ownerByAPIKey[batch[i].APIKey] = ownerUserID
		batch[i].OwnerUserID = ownerUserID
	}
}
