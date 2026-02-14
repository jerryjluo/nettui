package util

import "fmt"

// FormatBytes formats bytes into a human-readable string.
func FormatBytes(b uint64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
		TB = GB * 1024
	)
	switch {
	case b >= TB:
		return fmt.Sprintf("%.1f TB", float64(b)/float64(TB))
	case b >= GB:
		return fmt.Sprintf("%.1f GB", float64(b)/float64(GB))
	case b >= MB:
		return fmt.Sprintf("%.1f MB", float64(b)/float64(MB))
	case b >= KB:
		return fmt.Sprintf("%.1f KB", float64(b)/float64(KB))
	default:
		return fmt.Sprintf("%d B", b)
	}
}

// FormatRate formats bytes/sec into a human-readable throughput string.
func FormatRate(bytesPerSec float64) string {
	if bytesPerSec < 0 {
		bytesPerSec = 0
	}
	const (
		KB = 1024.0
		MB = KB * 1024
		GB = MB * 1024
	)
	switch {
	case bytesPerSec >= GB:
		return fmt.Sprintf("%.1f GB/s", bytesPerSec/GB)
	case bytesPerSec >= MB:
		return fmt.Sprintf("%.1f MB/s", bytesPerSec/MB)
	case bytesPerSec >= KB:
		return fmt.Sprintf("%.1f KB/s", bytesPerSec/KB)
	default:
		return fmt.Sprintf("%.0f B/s", bytesPerSec)
	}
}

// FormatPort returns the port as a string, or "*" if 0.
func FormatPort(port uint32) string {
	if port == 0 {
		return "*"
	}
	return fmt.Sprintf("%d", port)
}

// FormatAddr returns the address or "*" if empty.
func FormatAddr(addr string) string {
	if addr == "" || addr == "0.0.0.0" || addr == "::" {
		return "*"
	}
	return addr
}

// FormatAddrPort formats an address and port into addr:port.
func FormatAddrPort(addr string, port uint32) string {
	return fmt.Sprintf("%s:%s", FormatAddr(addr), FormatPort(port))
}

// FormatPID returns the PID as a string, or "--" if unavailable.
func FormatPID(pid int32) string {
	if pid <= 0 {
		return "--"
	}
	return fmt.Sprintf("%d", pid)
}

// FormatProcess returns the process name or "--" if unavailable.
func FormatProcess(name string) string {
	if name == "" {
		return "--"
	}
	return name
}
