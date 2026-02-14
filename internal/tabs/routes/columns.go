package routes

import "github.com/evertras/bubble-table/table"

func columns() []table.Column {
	return []table.Column{
		table.NewColumn("dest", "Destination", 18).WithFiltered(true),
		table.NewColumn("gateway", "Gateway", 18).WithFiltered(true),
		table.NewColumn("netmask", "Netmask", 18),
		table.NewColumn("iface", "Interface", 12).WithFiltered(true),
		table.NewColumn("flags", "Flags", 10),
	}
}
