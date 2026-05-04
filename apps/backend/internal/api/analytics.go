package api

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net"
	"net/http"
	"net/netip"
	"strings"
	"sync"
	"time"
)

type RequestEvent struct {
	Timestamp   time.Time
	Method      string
	Route       string
	Path        string
	StatusCode  int
	DurationMs  int64
	ClientIP    string
	Country     string
	UserAgent   string
	OwnerUserID string
	APIKey      string
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
	countryByIP   map[string]string
	geoHTTPClient *http.Client
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
		countryByIP:   make(map[string]string),
		geoHTTPClient: &http.Client{Timeout: 2 * time.Second},
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
		a.enrichCountries(ctx, batch)
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

func (a *AsyncAnalytics) enrichCountries(ctx context.Context, batch []RequestEvent) {
	for i := range batch {
		if strings.TrimSpace(batch[i].Country) != "" {
			continue
		}
		clientIP := strings.TrimSpace(batch[i].ClientIP)
		if clientIP == "" {
			continue
		}

		if cached, ok := a.countryByIP[clientIP]; ok {
			batch[i].Country = cached
			continue
		}

		country, ok := a.lookupCountryByIP(ctx, clientIP)
		if !ok {
			continue
		}
		a.countryByIP[clientIP] = country
		batch[i].Country = country
	}
}

func (a *AsyncAnalytics) lookupCountryByIP(ctx context.Context, ip string) (string, bool) {
	addr, ok := parseIP(ip)
	if !ok || !addr.IsGlobalUnicast() || addr.IsPrivate() || addr.IsLoopback() {
		return "", false
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://ipapi.co/"+addr.String()+"/country/", nil)
	if err != nil {
		return "", false
	}
	resp, err := a.geoHTTPClient.Do(req)
	if err == nil && resp != nil && resp.Body != nil {
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			body, readErr := io.ReadAll(resp.Body)
			if readErr == nil {
				countryCode := strings.ToUpper(strings.TrimSpace(string(body)))
				if len(countryCode) == 2 {
					return countryCode, true
				}
			}
		}
	}

	req, err = http.NewRequestWithContext(ctx, http.MethodGet, "https://ipwho.is/"+addr.String(), nil)
	if err != nil {
		return "", false
	}
	resp, err = a.geoHTTPClient.Do(req)
	if err != nil {
		return "", false
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", false
	}

	var payload struct {
		CountryCode string `json:"country_code"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return "", false
	}
	countryCode := strings.ToUpper(strings.TrimSpace(payload.CountryCode))
	if len(countryCode) != 2 {
		return "", false
	}
	return countryCode, true
}

func parseIP(value string) (netip.Addr, bool) {
	value = strings.TrimSpace(value)
	if value == "" {
		return netip.Addr{}, false
	}
	if addr, err := netip.ParseAddr(value); err == nil {
		return addr, true
	}

	host, _, err := net.SplitHostPort(value)
	if err != nil {
		return netip.Addr{}, false
	}
	addr, err := netip.ParseAddr(host)
	if err != nil {
		return netip.Addr{}, false
	}
	return addr, true
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
