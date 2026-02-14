package util

import (
	"os/exec"
	"strings"
)

// CopyToClipboard copies text to the macOS clipboard via pbcopy.
func CopyToClipboard(text string) error {
	cmd := exec.Command("pbcopy")
	cmd.Stdin = strings.NewReader(text)
	return cmd.Run()
}
