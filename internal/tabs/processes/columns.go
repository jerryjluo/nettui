package processes

import "github.com/evertras/bubble-table/table"

func columns() []table.Column {
	return []table.Column{
		table.NewColumn("pid", "PID", 8).WithFiltered(true),
		table.NewFlexColumn("name", "Name", 1).WithFiltered(true),
		table.NewFlexColumn("command", "Command", 2).WithFiltered(true),
		table.NewFlexColumn("user", "User", 1),
		table.NewColumn("conns", "#Sockets", 10),
		table.NewColumn("unix_socks", "#Unix", 8),
	}
}
