package processes

import "github.com/evertras/bubble-table/table"

func columns() []table.Column {
	return []table.Column{
		table.NewColumn("pid", "PID", 8).WithFiltered(true),
		table.NewColumn("name", "Name", 22).WithFiltered(true),
		table.NewColumn("user", "User", 14),
		table.NewColumn("conns", "#Connections", 14),
	}
}
