package sources

import (
	"fmt"

	"github.com/jerryluo/nettui/internal/data"
	psnet "github.com/shirou/gopsutil/v4/net"
)

// CollectInterfaces gathers network interface info and IO counters.
func CollectInterfaces() ([]data.Interface, []data.CollectionError) {
	var errs []data.CollectionError

	ifaces, err := psnet.Interfaces()
	if err != nil {
		return nil, []data.CollectionError{{Source: "interfaces", Error: fmt.Sprintf("Interfaces(): %v", err)}}
	}

	// Gather IO counters keyed by interface name.
	counters := make(map[string]psnet.IOCountersStat)
	ioStats, err := psnet.IOCounters(true)
	if err != nil {
		errs = append(errs, data.CollectionError{Source: "interfaces", Error: fmt.Sprintf("IOCounters(): %v", err)})
	} else {
		for _, c := range ioStats {
			counters[c.Name] = c
		}
	}

	result := make([]data.Interface, 0, len(ifaces))
	for _, iface := range ifaces {
		addrs := make([]string, len(iface.Addrs))
		for i, a := range iface.Addrs {
			addrs[i] = a.Addr
		}

		di := data.Interface{
			Name:  iface.Name,
			Index: iface.Index,
			MTU:   iface.MTU,
			Addrs: addrs,
			Up:    containsFlag(iface.Flags, "up"),
		}

		if iface.HardwareAddr != "" {
			di.HWAddr = iface.HardwareAddr
		}

		if c, ok := counters[iface.Name]; ok {
			di.BytesSent = c.BytesSent
			di.BytesRecv = c.BytesRecv
			di.PacketSent = c.PacketsSent
			di.PacketRecv = c.PacketsRecv
		}

		result = append(result, di)
	}

	return result, errs
}

func containsFlag(flags []string, target string) bool {
	for _, f := range flags {
		if f == target {
			return true
		}
	}
	return false
}
