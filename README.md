# nettui

A terminal UI for exploring network state on macOS. View sockets, processes, interfaces, routes, Unix domain sockets, and firewall rules in an interactive, tabbed interface with vim-style navigation.

## Features

- **6 data tabs** — Sockets (TCP/UDP), Unix Sockets, Processes, Interfaces, Routes, Firewall Rules
- **Cross-reference navigation** — Jump from a socket to its owning process, from a process to its sockets, or between connected local sockets
- **Search & filter** — Filter any table by typing `/` and entering a query
- **Protocol filtering** — Filter the sockets tab by TCP/UDP, IPv4/IPv6
- **Column sorting** — Sort any column ascending or descending
- **Detail side panel** — Press `p` to open a panel with full details for the selected row
- **Async DNS resolution** — Reverse-resolve remote addresses with a cached, concurrent resolver
- **Clipboard yank** — Copy specific fields to the clipboard with chord shortcuts
- **Throughput rates** — Live bytes/sec and packets/sec on the interfaces tab

## Requirements

- **macOS** (uses BSD routing APIs, `lsof`, and `pfctl`)
- **Go 1.25+**
- **Root privileges** recommended — required for PID mapping on sockets, Unix socket enumeration, and firewall rules

## Installation

```bash
go install github.com/jerryluo/nettui@latest
```

Or build from source:

```bash
git clone https://github.com/jerryluo/nettui.git
cd nettui
go build -o nettui .
```

## Usage

```bash
# Basic (some tabs will show warnings without root)
./nettui

# Full functionality
sudo ./nettui
```

### Keybindings

| Key | Action |
|-----|--------|
| `h`/`l` or `Tab`/`Shift+Tab` | Switch tabs |
| `1`–`6` | Jump to tab |
| `j`/`k` or `Up`/`Down` | Navigate rows |
| `d`/`u` | Page down / up |
| `/` | Search / filter |
| `Esc` | Clear filter or close panel |
| `p` | Toggle detail side panel |
| `r` | Refresh data |
| `D` | Toggle DNS resolution |
| `?` | Help screen |
| `q` / `Ctrl+C` | Quit |

**Chord shortcuts** (press leader key, then second key):

| Chord | Action |
|-------|--------|
| `g` + `s` | Go to sockets for selected process |
| `g` + `u` | Go to Unix sockets for selected process |
| `g` + `p` | Go to process for selected socket |
| `g` + `r` | Go to remote peer socket (localhost connections) |
| `f` + `t/u/4/6/c` | Filter by TCP / UDP / IPv4 / IPv6 / clear |
| `s` + column key | Sort by column |
| `y` + field key | Yank (copy) field to clipboard |

## Architecture

nettui is built with the [Bubble Tea](https://github.com/charmbracelet/bubbletea) TUI framework, following its Elm-inspired model-update-view pattern.

### Project structure

```
main.go                     Entry point — wires tabs and collector, starts Bubble Tea
internal/
  app/
    app.go                  Root model — manages tabs, panel, global key handling
    keys.go                 Keybinding definitions
  data/
    types.go                Core data types (Socket, Process, Interface, Route, etc.)
    store.go                Thread-safe data store with cross-reference indices
    sources/
      collector.go          Collection orchestrator — calls all sources, enriches data
      connections.go        TCP/UDP sockets via gopsutil
      processes.go          Process list via gopsutil
      interfaces.go         Network interfaces + IO counters via gopsutil
      routes.go             BSD routing table via golang.org/x/net/route
      lsof.go               PID-to-socket mapping and Unix sockets via lsof
      firewall.go           pf firewall rules via pfctl
      dns.go                Async reverse DNS with TTL cache
      throughput.go         Per-interface bytes/sec rate calculation
  tabs/
    tab.go                  Tab interface — all tabs implement this contract
    sort.go                 Generic column sorting (numeric + string)
    sockets/                TCP/UDP sockets tab
    unixsockets/            Unix domain sockets tab
    processes/              Process list tab
    interfaces/             Network interfaces tab
    routes/                 Routing table tab
    firewall/               Firewall rules tab
  ui/
    layout.go               Terminal layout calculation
    tabbar.go               Tab bar renderer
    statusbar.go            Status bar with hints and chord state
    sidepanel.go            Detail side panel renderer
  model/
    tabid.go                Tab identifier constants
  util/
    format.go               Address, byte size, and duration formatting
    clipboard.go            Clipboard integration (pbcopy)
```

### Data flow

1. **Collect** — `collector.Collect()` gathers data from system sources (gopsutil for connections/processes/interfaces, BSD route API, `lsof` for PID mapping, `pfctl` for firewall rules)
2. **Store** — Results are written to a `Store` that builds cross-reference indices (sockets by PID, processes by PID, routes by interface)
3. **Update tabs** — Each tab receives the updated store via `SetData()`, rebuilds its table rows, and reapplies any active sort or filter
4. **Render** — Bubble Tea calls `View()` on the root model, which composites the tab bar, active tab table, status bar, and optional side panel

### Key dependencies

| Library | Purpose |
|---------|---------|
| [charmbracelet/bubbletea](https://github.com/charmbracelet/bubbletea) | TUI framework |
| [charmbracelet/lipgloss](https://github.com/charmbracelet/lipgloss) | Terminal styling |
| [evertras/bubble-table](https://github.com/evertras/bubble-table) | Interactive table widget with filtering |
| [shirou/gopsutil](https://github.com/shirou/gopsutil) | Cross-platform system info (connections, processes, interfaces) |
| [golang.org/x/net](https://pkg.go.dev/golang.org/x/net) | BSD routing table access |

## License

MIT
