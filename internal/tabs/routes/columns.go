package routes

import "github.com/evertras/bubble-table/table"

func columns() []table.Column {
	return []table.Column{
		table.NewFlexColumn("dest", "Destination", 1).WithFiltered(true),
		table.NewFlexColumn("gateway", "Gateway", 1).WithFiltered(true),
		table.NewFlexColumn("netmask", "Netmask", 1),
		table.NewColumn("iface", "Interface", 12).WithFiltered(true),
		table.NewColumn("flags", "Flags", 10),
	}
}
