# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build & Run

```bash
go build -o nettui .          # build
go vet ./...                  # lint
sudo ./nettui                 # run with full functionality (root needed for PID mapping, unix sockets, firewall)
./nettui                      # run without root (some tabs show warnings)
```

There are no tests in this project.

## Architecture

nettui is a macOS-only terminal network monitor built with [Bubble Tea](https://github.com/charmbracelet/bubbletea) (Elm architecture: Model → Update → View).

### Data flow

`collector.Collect()` → `Store.Update()` (rebuilds cross-ref indices) → each tab receives store via `SetData()` → Bubble Tea calls `View()`

### Key architectural patterns

**Tab interface** (`internal/tabs/tab.go`): Every tab implements `tabs.Tab`, which extends `tea.Model` with methods for data binding (`SetData`, `SetSize`), cross-referencing (`CrossRef`, `NavigateTo`), filtering (`IsFiltering`, `HasActiveFilter`, `ClearFilter`), sorting (`ApplySort`, `SortHint`, `SortLabel`), and clipboard (`YankField`, `YankHint`). Adding a new tab means implementing this interface.

**Each tab has 3 files**: `model.go` (state + Update logic), `columns.go` (bubble-table column definitions), `detail.go` (side panel content rendering).

**Chord keybindings** (`internal/app/app.go`): Leader keys (`g`, `f`, `s`, `y`) set `pendingChord` and show a `chordHint` in the status bar. The second keypress dispatches via `handleChordSecondKey` → `handleGotoChord`/`handleFilterChord`/`handleSortChord`/`handleYankChord`. Chords auto-clear after 2 seconds.

**Cross-reference navigation**: `CrossRefMsg` carries a target tab + filter key/value. The root model switches tabs and calls `NavigateTo()` on the target. The `Store` maintains index maps (`SocketsByPID`, `ProcessByPID`, `RoutesByIface`, `IfaceByName`) for fast lookups.

**Tab ordering**: Tabs are created in `main.go` as a slice. The index position must match `model.TabID` constants in `internal/model/tabid.go`. When the root model references `m.tabs[m.activeTab]`, it indexes into this slice directly.

### Column filtering

bubble-table's `.WithFiltered(true)` on a column definition makes it searchable via `/`. To make a column searchable, add this flag in the tab's `columns.go`.

### Data sources

All system data collection lives in `internal/data/sources/`. The collector calls gopsutil for connections/processes/interfaces, `lsof` for PID-to-socket mapping and unix socket enumeration, BSD `route` API for routing table, and `pfctl` for firewall rules. The DNS resolver (`dns.go`) runs async with a TTL cache.
