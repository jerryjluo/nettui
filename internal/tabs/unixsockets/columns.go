package unixsockets

import "github.com/evertras/bubble-table/table"

func columns() []table.Column {
	return []table.Column{
		table.NewColumn("path", "Path", 35).WithFiltered(true),
		table.NewColumn("type", "Type", 10),
		table.NewColumn("pid", "PID", 8),
		table.NewColumn("process", "Process", 18).WithFiltered(true),
		table.NewColumn("fd", "FD", 8),
	}
}
