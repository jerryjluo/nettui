package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jerryluo/nettui/internal/app"
	"github.com/jerryluo/nettui/internal/data/sources"
	"github.com/jerryluo/nettui/internal/tabs"
	"github.com/jerryluo/nettui/internal/tabs/arp"
	"github.com/jerryluo/nettui/internal/tabs/firewall"
	"github.com/jerryluo/nettui/internal/tabs/interfaces"
	"github.com/jerryluo/nettui/internal/tabs/processes"
	"github.com/jerryluo/nettui/internal/tabs/routes"
	"github.com/jerryluo/nettui/internal/tabs/sockets"
	"github.com/jerryluo/nettui/internal/tabs/unixsockets"
)

func main() {
	collector := sources.NewCollector()

	tabModels := []tabs.Tab{
		sockets.New(collector.DNSCache()),
		unixsockets.New(),
		processes.New(),
		interfaces.New(),
		routes.New(),
		arp.New(),
		firewall.New(),
	}

	model := app.New(tabModels, collector)

	p := tea.NewProgram(
		model,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
