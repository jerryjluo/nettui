package sources

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/jerryluo/nettui/internal/data"
)

var (
	// Matches lines like: @0 pass in on en0 proto tcp from any to any port = 22
	ruleLineRe = regexp.MustCompile(`^@(\d+)\s+(.*)$`)
	// Matches evaluation counters like: [ Evaluations: 1234  Packets: 5678  Bytes: 91011  States: 12 ]
	evalLineRe = regexp.MustCompile(`\[\s*Evaluations:\s*\d+\s+Packets:\s*(\d+)\s+Bytes:\s*(\d+)`)
)

// CollectFirewall parses pfctl -vsr output to collect firewall rules.
// Requires root access; returns an error if not root.
func CollectFirewall(isRoot bool) ([]data.FirewallRule, []data.CollectionError) {
	if !isRoot {
		return nil, []data.CollectionError{{Source: "firewall", Error: "pfctl requires root access"}}
	}

	out, err := exec.Command("pfctl", "-vsr").CombinedOutput()
	if err != nil {
		return nil, []data.CollectionError{{Source: "firewall", Error: fmt.Sprintf("pfctl -vsr: %v: %s", err, string(out))}}
	}

	return parsePfctlOutput(string(out)), nil
}

func parsePfctlOutput(output string) []data.FirewallRule {
	var rules []data.FirewallRule
	lines := strings.Split(output, "\n")

	for i := 0; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}

		m := ruleLineRe.FindStringSubmatch(line)
		if m == nil {
			continue
		}

		ruleNum, _ := strconv.Atoi(m[1])
		ruleText := m[2]

		rule := data.FirewallRule{
			RuleNum: ruleNum,
			RawRule: ruleText,
		}

		// Parse action, direction, proto, src, dst from the rule text.
		parseRuleFields(ruleText, &rule)

		// Look ahead for evaluation counters.
		if i+1 < len(lines) {
			em := evalLineRe.FindStringSubmatch(lines[i+1])
			if em != nil {
				rule.Packets, _ = strconv.ParseUint(em[1], 10, 64)
				rule.Bytes, _ = strconv.ParseUint(em[2], 10, 64)
				i++ // Skip the counter line.
			}
		}

		rules = append(rules, rule)
	}

	return rules
}

func parseRuleFields(ruleText string, rule *data.FirewallRule) {
	fields := strings.Fields(ruleText)
	if len(fields) == 0 {
		return
	}

	// Action is first word: pass, block, match, etc.
	rule.Action = fields[0]

	for i, f := range fields {
		switch f {
		case "in", "out":
			if rule.Direction == "" {
				rule.Direction = f
			}
		case "proto":
			if i+1 < len(fields) {
				rule.Proto = fields[i+1]
			}
		case "from":
			if i+1 < len(fields) {
				rule.Src = collectAddrSpec(fields, i+1)
			}
		case "to":
			if i+1 < len(fields) {
				rule.Dst = collectAddrSpec(fields, i+1)
			}
		}
	}
}

// collectAddrSpec extracts address spec following "from" or "to".
// It takes the next token, and if followed by "port", appends that.
func collectAddrSpec(fields []string, start int) string {
	if start >= len(fields) {
		return ""
	}
	addr := fields[start]
	if start+2 < len(fields) && fields[start+1] == "port" {
		addr += " port " + fields[start+2]
	}
	return addr
}
