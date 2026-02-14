package sources

import (
	"context"
	"net"
	"sync"
	"time"
)

const (
	dnsCacheTTL    = 5 * time.Minute
	dnsTimeout     = 500 * time.Millisecond
	dnsMaxWorkers  = 10
)

type dnsEntry struct {
	hostname  string
	expiresAt time.Time
}

// DNSCache provides cached asynchronous reverse DNS lookups.
type DNSCache struct {
	mu      sync.RWMutex
	cache   map[string]dnsEntry
	pending sync.Map // tracks IPs currently being resolved

	sem chan struct{} // concurrency limiter
}

// NewDNSCache creates a new DNS cache with bounded concurrency.
func NewDNSCache() *DNSCache {
	return &DNSCache{
		cache: make(map[string]dnsEntry),
		sem:   make(chan struct{}, dnsMaxWorkers),
	}
}

// Lookup returns the cached hostname for an IP, or the IP itself if not cached.
// It triggers an async resolution if the IP is not in the cache or has expired.
func (d *DNSCache) Lookup(ip string) string {
	d.mu.RLock()
	entry, ok := d.cache[ip]
	d.mu.RUnlock()

	if ok && time.Now().Before(entry.expiresAt) {
		return entry.hostname
	}

	// Trigger async resolution if not already pending.
	if _, loaded := d.pending.LoadOrStore(ip, struct{}{}); !loaded {
		go d.resolve(ip)
	}

	// Return cached value if we have one (even if expired), otherwise the raw IP.
	if ok {
		return entry.hostname
	}
	return ip
}

func (d *DNSCache) resolve(ip string) {
	defer d.pending.Delete(ip)

	d.sem <- struct{}{}
	defer func() { <-d.sem }()

	ctx, cancel := context.WithTimeout(context.Background(), dnsTimeout)
	defer cancel()

	resolver := &net.Resolver{}
	names, err := resolver.LookupAddr(ctx, ip)

	hostname := ip
	if err == nil && len(names) > 0 {
		// Remove trailing dot from FQDN.
		h := names[0]
		if len(h) > 0 && h[len(h)-1] == '.' {
			h = h[:len(h)-1]
		}
		hostname = h
	}

	d.mu.Lock()
	d.cache[ip] = dnsEntry{
		hostname:  hostname,
		expiresAt: time.Now().Add(dnsCacheTTL),
	}
	d.mu.Unlock()
}

// ResolveAll triggers async lookups for a batch of IPs.
func (d *DNSCache) ResolveAll(ips []string) {
	for _, ip := range ips {
		if ip == "" || ip == "*" || ip == "0.0.0.0" || ip == "::" {
			continue
		}
		d.Lookup(ip) // triggers async resolution as a side effect
	}
}
