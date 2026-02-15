package sources

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"github.com/jerryluo/nettui/internal/data"
)

// arpLineRe matches lines like:
// ? (192.168.1.1) at 68:d7:9a:6a:20:6 on en0 ifscope [ethernet]
// mdns.mcast.net (224.0.0.251) at (incomplete) on en0 ifscope permanent [ethernet]
var arpLineRe = regexp.MustCompile(
	`^(\S+)\s+\(([^)]+)\)\s+at\s+(\S+)\s+on\s+(\S+)(.*)$`,
)

// CollectARP runs `arp -a` and parses the output.
func CollectARP() ([]data.ARPEntry, []data.CollectionError) {
	out, err := exec.Command("arp", "-a").CombinedOutput()
	if err != nil {
		return nil, []data.CollectionError{{Source: "arp", Error: fmt.Sprintf("arp -a: %v: %s", err, string(out))}}
	}
	return parseARPOutput(string(out)), nil
}

func parseARPOutput(output string) []data.ARPEntry {
	var entries []data.ARPEntry
	for _, line := range strings.Split(output, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		m := arpLineRe.FindStringSubmatch(line)
		if m == nil {
			continue
		}

		hostname := m[1]
		if hostname == "?" {
			hostname = ""
		}

		entry := data.ARPEntry{
			Hostname:  hostname,
			IP:        m[2],
			MAC:       m[3],
			Interface: m[4],
		}

		// Parse trailing flags like "ifscope permanent [ethernet]"
		rest := strings.TrimSpace(m[5])
		if rest != "" {
			var flags []string
			if strings.Contains(rest, "ifscope") {
				flags = append(flags, "ifscope")
			}
			if strings.Contains(rest, "permanent") {
				flags = append(flags, "permanent")
			}
			entry.Flags = strings.Join(flags, ", ")

			// Extract type from brackets
			if idx := strings.Index(rest, "["); idx != -1 {
				if end := strings.Index(rest[idx:], "]"); end != -1 {
					entry.Type = rest[idx+1 : idx+end]
				}
			}
		}

		entries = append(entries, entry)
	}
	return entries
}
