package sources

import (
	"fmt"

	"github.com/jerryluo/nettui/internal/data"
	psnet "github.com/shirou/gopsutil/v4/net"
)

// CollectConnections gathers all TCP and UDP connections.
func CollectConnections() ([]data.Socket, []data.CollectionError) {
	conns, err := psnet.Connections("all")
	if err != nil {
		return nil, []data.CollectionError{{Source: "connections", Error: fmt.Sprintf("Connections(): %v", err)}}
	}

	sockets := make([]data.Socket, 0, len(conns))
	for _, c := range conns {
		proto := protoName(c.Type, c.Family)
		if proto == "" {
			continue
		}
		sockets = append(sockets, data.Socket{
			Proto:      proto,
			LocalAddr:  c.Laddr.IP,
			LocalPort:  c.Laddr.Port,
			RemoteAddr: c.Raddr.IP,
			RemotePort: c.Raddr.Port,
			State:      c.Status,
			PID:        c.Pid,
		})
	}

	return sockets, nil
}

func protoName(connType uint32, family uint32) string {
	// gopsutil: type 1=SOCK_STREAM (TCP), type 2=SOCK_DGRAM (UDP)
	// family 2=AF_INET, 10/30=AF_INET6
	switch connType {
	case 1: // TCP
		if family == 2 {
			return "tcp"
		}
		return "tcp6"
	case 2: // UDP
		if family == 2 {
			return "udp"
		}
		return "udp6"
	default:
		return ""
	}
}
