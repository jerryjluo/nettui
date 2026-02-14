package sources

import (
	"fmt"
	"time"

	"github.com/jerryluo/nettui/internal/data"
	"github.com/jerryluo/nettui/internal/util"
)

// Collector orchestrates data collection from all sources.
type Collector struct {
	isRoot     bool
	throughput *ThroughputCalculator
	dns        *DNSCache
}

// NewCollector creates a new Collector.
func NewCollector() *Collector {
	return &Collector{
		isRoot:     util.IsRoot(),
		throughput: NewThroughputCalculator(),
		dns:        NewDNSCache(),
	}
}

// DNSCache returns the collector's DNS cache for use by the UI.
func (c *Collector) DNSCache() *DNSCache {
	return c.dns
}

// Collect gathers data from all sources and returns a CollectionResult.
func (c *Collector) Collect() data.CollectionResult {
	result := data.CollectionResult{
		Timestamp: time.Now(),
		IsRoot:    c.isRoot,
	}

	// Interfaces + IO counters.
	ifaces, errs := CollectInterfaces()
	result.Interfaces = ifaces
	result.Errors = append(result.Errors, errs...)

	// Calculate throughput from interface counters.
	throughputs := c.throughput.Calculate(result.Interfaces)
	result.Throughputs = throughputs

	// Apply throughput rates back to interfaces.
	for i := range result.Interfaces {
		if tp, ok := throughputs[result.Interfaces[i].Name]; ok {
			result.Interfaces[i].TxRate = tp.TxRate
			result.Interfaces[i].RxRate = tp.RxRate
		}
	}

	// Routes.
	routes, errs := c.collectRoutes()
	result.Routes = routes
	result.Errors = append(result.Errors, errs...)

	// Connections.
	sockets, errs := CollectConnections()
	result.Sockets = sockets
	result.Errors = append(result.Errors, errs...)

	// Lsof: enrich sockets with PID/process info and get unix sockets.
	lsofResult, errs := CollectLsof()
	result.Errors = append(result.Errors, errs...)
	if lsofResult != nil {
		EnrichSockets(result.Sockets, lsofResult)
		result.UnixSockets = lsofResult.UnixSockets
	}

	// Processes.
	procs, errs := CollectProcesses()
	result.Errors = append(result.Errors, errs...)

	// Enrich processes with connection counts from sockets.
	pidConns := make(map[int32]int)
	for _, s := range result.Sockets {
		if s.PID > 0 {
			pidConns[s.PID]++
		}
	}
	pidUnix := make(map[int32]int)
	for _, u := range result.UnixSockets {
		if u.PID > 0 {
			pidUnix[u.PID]++
		}
	}
	for i := range procs {
		procs[i].NumConns = pidConns[procs[i].PID]
		procs[i].NumUnixSocks = pidUnix[procs[i].PID]
	}
	result.Processes = procs

	// Firewall (requires root).
	fwRules, errs := CollectFirewall(c.isRoot)
	result.Firewall = fwRules
	result.Errors = append(result.Errors, errs...)

	// Trigger async DNS resolution for unique remote addresses.
	c.triggerDNS(result.Sockets)

	return result
}

func (c *Collector) collectRoutes() ([]data.Route, []data.CollectionError) {
	defer func() {
		if r := recover(); r != nil {
			_ = fmt.Sprintf("routes panic: %v", r)
		}
	}()
	return CollectRoutes()
}

func (c *Collector) triggerDNS(sockets []data.Socket) {
	ips := make([]string, 0, len(sockets)*2)
	for _, s := range sockets {
		ips = append(ips, s.RemoteAddr)
	}
	c.dns.ResolveAll(ips)
}
