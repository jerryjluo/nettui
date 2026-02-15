package data

import (
	"net"
	"time"
)

// Interface represents a network interface with IO stats.
type Interface struct {
	Name       string
	Index      int
	MTU        int
	Flags      net.Flags
	HWAddr     string
	Addrs      []string
	Up         bool
	BytesSent  uint64
	BytesRecv  uint64
	PacketSent uint64
	PacketRecv uint64
	TxRate     float64 // bytes/sec
	RxRate     float64 // bytes/sec
}

// Route represents a routing table entry.
type Route struct {
	Destination string
	Gateway     string
	Netmask     string
	Interface   string
	Flags       string
}

// Socket represents a TCP or UDP connection.
type Socket struct {
	Proto      string // tcp, tcp6, udp, udp6
	LocalAddr  string
	LocalPort  uint32
	RemoteAddr string
	RemotePort uint32
	State      string
	PID        int32
	Process    string
}

// UnixSocket represents a Unix domain socket.
type UnixSocket struct {
	Path    string
	Type    string // stream, dgram
	State   string // LISTEN, CONNECTED, etc.
	PID     int32
	Process string
	FD      string
}

// Process represents a process with network activity.
type Process struct {
	PID          int32
	Name         string
	Command      string
	User         string
	NumConns     int
	NumUnixSocks int
	Connections  []Socket
}

// FirewallRule represents a pf firewall rule.
type FirewallRule struct {
	RuleNum    int
	Action     string // pass, block
	Direction  string // in, out
	Proto      string
	Src        string
	Dst        string
	Packets    uint64
	Bytes      uint64
	RawRule    string
}

// Throughput holds per-interface throughput data.
type Throughput struct {
	Interface string
	TxRate    float64 // bytes/sec
	RxRate    float64 // bytes/sec
}

// CollectionResult holds the result of a single data collection cycle.
type CollectionResult struct {
	Interfaces  []Interface
	Routes      []Route
	Sockets     []Socket
	UnixSockets []UnixSocket
	Processes   []Process
	Firewall    []FirewallRule
	Throughputs map[string]Throughput
	Errors      []CollectionError
	Timestamp   time.Time
	IsRoot      bool
}

// CollectionError records a non-fatal error during collection.
type CollectionError struct {
	Source  string
	Error  string
}
