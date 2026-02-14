package sources

import (
	"time"

	"github.com/jerryluo/nettui/internal/data"
)

// ThroughputCalculator computes per-interface byte rates from IO counter deltas.
type ThroughputCalculator struct {
	prevCounters map[string]ifaceCounters
	prevTime     time.Time
}

type ifaceCounters struct {
	bytesSent uint64
	bytesRecv uint64
}

// NewThroughputCalculator creates a new calculator.
func NewThroughputCalculator() *ThroughputCalculator {
	return &ThroughputCalculator{
		prevCounters: make(map[string]ifaceCounters),
	}
}

// Calculate computes throughput for each interface by comparing current counters
// against the previous snapshot. On the first call, all rates will be zero.
func (tc *ThroughputCalculator) Calculate(interfaces []data.Interface) map[string]data.Throughput {
	now := time.Now()
	result := make(map[string]data.Throughput, len(interfaces))

	elapsed := now.Sub(tc.prevTime).Seconds()
	if elapsed <= 0 || tc.prevTime.IsZero() {
		// First call or invalid interval: store counters, return zero rates.
		for _, iface := range interfaces {
			tc.prevCounters[iface.Name] = ifaceCounters{
				bytesSent: iface.BytesSent,
				bytesRecv: iface.BytesRecv,
			}
			result[iface.Name] = data.Throughput{Interface: iface.Name}
		}
		tc.prevTime = now
		return result
	}

	for _, iface := range interfaces {
		tp := data.Throughput{Interface: iface.Name}

		if prev, ok := tc.prevCounters[iface.Name]; ok {
			if iface.BytesSent >= prev.bytesSent {
				tp.TxRate = float64(iface.BytesSent-prev.bytesSent) / elapsed
			}
			if iface.BytesRecv >= prev.bytesRecv {
				tp.RxRate = float64(iface.BytesRecv-prev.bytesRecv) / elapsed
			}
		}

		tc.prevCounters[iface.Name] = ifaceCounters{
			bytesSent: iface.BytesSent,
			bytesRecv: iface.BytesRecv,
		}
		result[iface.Name] = tp
	}

	tc.prevTime = now
	return result
}
