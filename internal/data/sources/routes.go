package sources

import (
	"fmt"
	"net"
	"syscall"

	"github.com/jerryluo/nettui/internal/data"
	"golang.org/x/net/route"
)

// CollectRoutes reads the Darwin routing table via route.FetchRIB.
func CollectRoutes() ([]data.Route, []data.CollectionError) {
	rib, err := route.FetchRIB(syscall.AF_UNSPEC, route.RIBTypeRoute, 0)
	if err != nil {
		return nil, []data.CollectionError{{Source: "routes", Error: fmt.Sprintf("FetchRIB(): %v", err)}}
	}

	msgs, err := route.ParseRIB(route.RIBTypeRoute, rib)
	if err != nil {
		return nil, []data.CollectionError{{Source: "routes", Error: fmt.Sprintf("ParseRIB(): %v", err)}}
	}

	// Look up interface names by index.
	ifaces, _ := net.Interfaces()
	ifaceNames := make(map[int]string, len(ifaces))
	for _, iface := range ifaces {
		ifaceNames[iface.Index] = iface.Name
	}

	var routes []data.Route
	for _, msg := range msgs {
		rm, ok := msg.(*route.RouteMessage)
		if !ok {
			continue
		}

		r := data.Route{
			Flags: fmt.Sprintf("0x%x", rm.Flags),
		}

		if rm.Index > 0 {
			r.Interface = ifaceNames[rm.Index]
		}

		addrs := rm.Addrs
		if len(addrs) > 0 && addrs[0] != nil {
			r.Destination = formatAddr(addrs[0])
		}
		if len(addrs) > 1 && addrs[1] != nil {
			r.Gateway = formatAddr(addrs[1])
		}
		if len(addrs) > 2 && addrs[2] != nil {
			r.Netmask = formatAddr(addrs[2])
		}

		routes = append(routes, r)
	}

	return routes, nil
}

func formatAddr(a route.Addr) string {
	switch v := a.(type) {
	case *route.Inet4Addr:
		return net.IP(v.IP[:]).String()
	case *route.Inet6Addr:
		return net.IP(v.IP[:]).String()
	case *route.LinkAddr:
		if v.Name != "" {
			return v.Name
		}
		return fmt.Sprintf("link#%d", v.Index)
	default:
		return ""
	}
}
