package data

import "sync"

// Store holds the latest snapshot plus cross-reference indices.
type Store struct {
	mu sync.RWMutex

	Interfaces  []Interface
	Routes      []Route
	Sockets     []Socket
	UnixSockets []UnixSocket
	Processes   []Process
	Firewall    []FirewallRule
	Throughputs map[string]Throughput
	Errors      []CollectionError
	IsRoot      bool

	// Cross-reference indices
	SocketsByPID  map[int32][]Socket
	ProcessByPID  map[int32]*Process
	RoutesByIface map[string][]Route
	IfaceByName   map[string]*Interface
}

// NewStore creates an empty Store.
func NewStore() *Store {
	return &Store{
		Throughputs:   make(map[string]Throughput),
		SocketsByPID:  make(map[int32][]Socket),
		ProcessByPID:  make(map[int32]*Process),
		RoutesByIface: make(map[string][]Route),
		IfaceByName:   make(map[string]*Interface),
	}
}

// Update replaces the store with a new collection result and rebuilds indices.
func (s *Store) Update(result CollectionResult) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Interfaces = result.Interfaces
	s.Routes = result.Routes
	s.Sockets = result.Sockets
	s.UnixSockets = result.UnixSockets
	s.Processes = result.Processes
	s.Firewall = result.Firewall
	s.Throughputs = result.Throughputs
	s.Errors = result.Errors
	s.IsRoot = result.IsRoot

	s.rebuildIndices()
}

func (s *Store) rebuildIndices() {
	s.SocketsByPID = make(map[int32][]Socket, len(s.Processes))
	for _, sock := range s.Sockets {
		if sock.PID > 0 {
			s.SocketsByPID[sock.PID] = append(s.SocketsByPID[sock.PID], sock)
		}
	}

	s.ProcessByPID = make(map[int32]*Process, len(s.Processes))
	for i := range s.Processes {
		s.ProcessByPID[s.Processes[i].PID] = &s.Processes[i]
	}

	s.RoutesByIface = make(map[string][]Route)
	for _, r := range s.Routes {
		if r.Interface != "" {
			s.RoutesByIface[r.Interface] = append(s.RoutesByIface[r.Interface], r)
		}
	}

	s.IfaceByName = make(map[string]*Interface, len(s.Interfaces))
	for i := range s.Interfaces {
		s.IfaceByName[s.Interfaces[i].Name] = &s.Interfaces[i]
	}
}

// Snapshot returns a read-locked copy of the current store data.
func (s *Store) Snapshot() *Store {
	s.mu.RLock()
	defer s.mu.RUnlock()

	snap := &Store{
		Interfaces:    make([]Interface, len(s.Interfaces)),
		Routes:        make([]Route, len(s.Routes)),
		Sockets:       make([]Socket, len(s.Sockets)),
		UnixSockets:   make([]UnixSocket, len(s.UnixSockets)),
		Processes:     make([]Process, len(s.Processes)),
		Firewall:      make([]FirewallRule, len(s.Firewall)),
		Throughputs:   make(map[string]Throughput, len(s.Throughputs)),
		Errors:        make([]CollectionError, len(s.Errors)),
		IsRoot:        s.IsRoot,
		SocketsByPID:  s.SocketsByPID,
		ProcessByPID:  s.ProcessByPID,
		RoutesByIface: s.RoutesByIface,
		IfaceByName:   s.IfaceByName,
	}
	copy(snap.Interfaces, s.Interfaces)
	copy(snap.Routes, s.Routes)
	copy(snap.Sockets, s.Sockets)
	copy(snap.UnixSockets, s.UnixSockets)
	copy(snap.Processes, s.Processes)
	copy(snap.Firewall, s.Firewall)
	copy(snap.Errors, s.Errors)
	for k, v := range s.Throughputs {
		snap.Throughputs[k] = v
	}
	return snap
}
