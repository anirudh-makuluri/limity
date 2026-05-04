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
	metrics       *Metrics
	events        chan RequestEvent
	batchSize     int
	flushInterval time.Duration
	flushTimeout  time.Duration

	stopCh chan struct{}
	wg     sync.WaitGroup

	ownerByAPIKey map[string]string
}

func NewAsyncAnalytics(store AnalyticsStore, metrics *Metrics, queueSize, batchSize int, flushInterval, flushTimeout time.Duration) *AsyncAnalytics {
	if queueSize <= 0 {
		queueSize = 10000
	}
	if batchSize <= 0 {
		batchSize = 200
	}
	if flushInterval <= 0 {
		flushInterval = time.Second
	}
	if flushTimeout <= 0 {
		flushTimeout = 15 * time.Second
	}

	a := &AsyncAnalytics{
		store:         store,
		metrics:       metrics,
		events:        make(chan RequestEvent, queueSize),
		batchSize:     batchSize,
		flushInterval: flushInterval,
		flushTimeout:  flushTimeout,
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
		ctx, cancel := context.WithTimeout(context.Background(), a.flushTimeout)
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
			if a.metrics != nil {
				a.metrics.ownerLookupSuccessTotal.Inc()
			}
			continue
		}

		ownerUserID, err := a.store.ResolveOwnerUserIDByAPIKey(ctx, batch[i].APIKey)
		if err != nil {
			log.Printf("analytics owner lookup error for /check event: %v", err)
			if a.metrics != nil {
				a.metrics.ownerLookupErrorTotal.Inc()
			}
			continue
		}
		if ownerUserID == "" {
			log.Printf("analytics owner lookup miss for /check event key prefix=%s", safeKeyPrefix(batch[i].APIKey))
			if a.metrics != nil {
				a.metrics.ownerLookupMissTotal.Inc()
			}
			continue
		}
		a.ownerByAPIKey[batch[i].APIKey] = ownerUserID
		batch[i].OwnerUserID = ownerUserID
		if a.metrics != nil {
			a.metrics.ownerLookupSuccessTotal.Inc()
		}
	}
}

func safeKeyPrefix(apiKey string) string {
	if len(apiKey) <= 12 {
		return apiKey
	}
	return apiKey[:12]
}
