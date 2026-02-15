package sockets

import "github.com/evertras/bubble-table/table"

func columns() []table.Column {
	return []table.Column{
		table.NewColumn("proto", "Proto", 7),
		table.NewFlexColumn("local", "Local Address", 1).WithFiltered(true),
		table.NewFlexColumn("remote", "Remote Address", 1).WithFiltered(true),
		table.NewColumn("state", "State", 14),
		table.NewColumn("pid", "PID", 8).WithFiltered(true),
		table.NewColumn("process", "Process", 18).WithFiltered(true),
	}
}
