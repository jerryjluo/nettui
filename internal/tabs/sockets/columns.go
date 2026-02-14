package sockets

import "github.com/evertras/bubble-table/table"

func columns() []table.Column {
	return []table.Column{
		table.NewColumn("proto", "Proto", 7),
		table.NewColumn("local", "Local Address", 24).WithFiltered(true),
		table.NewColumn("remote", "Remote Address", 24).WithFiltered(true),
		table.NewColumn("state", "State", 14),
		table.NewColumn("pid", "PID", 8),
		table.NewColumn("process", "Process", 18).WithFiltered(true),
	}
}
