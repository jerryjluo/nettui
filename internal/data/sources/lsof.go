package sources

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/jerryluo/nettui/internal/data"
)

// LsofResult holds parsed lsof data for enriching sockets and collecting unix sockets.
type LsofResult struct {
	// PIDProcess maps PID to process name from lsof.
	PIDProcess map[int32]string
	// SocketPIDs maps "localAddr:localPort" to PID for inet sockets.
	SocketPIDs map[string]int32
	// UnixSockets are the parsed unix domain sockets.
	UnixSockets []data.UnixSocket
}

// CollectLsof runs lsof to gather inet socket-to-PID mappings and unix sockets.
func CollectLsof() (*LsofResult, []data.CollectionError) {
	var errs []data.CollectionError
	result := &LsofResult{
		PIDProcess: make(map[int32]string),
		SocketPIDs: make(map[string]int32),
	}

	// Collect inet sockets.
	inetOut, err := exec.Command("lsof", "-i", "-P", "-n", "-F", "pcfn").Output()
	if err != nil {
		errs = append(errs, data.CollectionError{Source: "lsof-inet", Error: fmt.Sprintf("lsof -i: %v", err)})
	} else {
		parseInetLsof(string(inetOut), result)
	}

	// Collect unix sockets.
	unixOut, err := exec.Command("lsof", "-U", "-F", "pcfn").Output()
	if err != nil {
		errs = append(errs, data.CollectionError{Source: "lsof-unix", Error: fmt.Sprintf("lsof -U: %v", err)})
	} else {
		result.UnixSockets = parseUnixLsof(string(unixOut))
	}

	return result, errs
}

func parseInetLsof(output string, result *LsofResult) {
	var currentPID int32
	var currentCmd string

	for _, line := range strings.Split(output, "\n") {
		if len(line) == 0 {
			continue
		}
		field := line[0]
		value := line[1:]

		switch field {
		case 'p':
			pid, err := strconv.ParseInt(value, 10, 32)
			if err == nil {
				currentPID = int32(pid)
			}
		case 'c':
			currentCmd = value
			if currentPID > 0 {
				result.PIDProcess[currentPID] = currentCmd
			}
		case 'n':
			// value is like "192.168.1.5:443" or "*:80" or "127.0.0.1:8080->10.0.0.1:443"
			// We want the local part (before ->).
			local := value
			if idx := strings.Index(local, "->"); idx >= 0 {
				local = local[:idx]
			}
			if currentPID > 0 && strings.Contains(local, ":") {
				result.SocketPIDs[local] = currentPID
			}
		}
	}
}

func parseUnixLsof(output string) []data.UnixSocket {
	var sockets []data.UnixSocket
	var currentPID int32
	var currentCmd string
	var currentFD string

	for _, line := range strings.Split(output, "\n") {
		if len(line) == 0 {
			continue
		}
		field := line[0]
		value := line[1:]

		switch field {
		case 'p':
			pid, err := strconv.ParseInt(value, 10, 32)
			if err == nil {
				currentPID = int32(pid)
			}
		case 'c':
			currentCmd = value
		case 'f':
			currentFD = value
		case 'n':
			if currentPID > 0 {
				sockets = append(sockets, data.UnixSocket{
					Path:    value,
					PID:     currentPID,
					Process: currentCmd,
					FD:      currentFD,
					Type:    "stream", // lsof -F doesn't easily distinguish; default to stream
				})
			}
		}
	}

	return sockets
}

// EnrichSockets fills in Process names on sockets using lsof data.
func EnrichSockets(sockets []data.Socket, lsofResult *LsofResult) {
	if lsofResult == nil {
		return
	}
	for i := range sockets {
		s := &sockets[i]
		// Try to match by PID first.
		if s.PID > 0 {
			if name, ok := lsofResult.PIDProcess[s.PID]; ok {
				s.Process = name
			}
		}
		// Try to fill PID from lsof socket mapping if missing.
		if s.PID == 0 {
			key := fmt.Sprintf("%s:%d", s.LocalAddr, s.LocalPort)
			if pid, ok := lsofResult.SocketPIDs[key]; ok {
				s.PID = pid
				if name, found := lsofResult.PIDProcess[pid]; found {
					s.Process = name
				}
			}
			// Also try wildcard match.
			wildKey := fmt.Sprintf("*:%d", s.LocalPort)
			if pid, ok := lsofResult.SocketPIDs[wildKey]; ok && s.PID == 0 {
				s.PID = pid
				if name, found := lsofResult.PIDProcess[pid]; found {
					s.Process = name
				}
			}
		}
	}
}
