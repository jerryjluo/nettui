package routes

import "github.com/evertras/bubble-table/table"

func columns() []table.Column {
	return []table.Column{
		table.NewColumn("dest", "Destination", 18),
		table.NewColumn("gateway", "Gateway", 18),
		table.NewColumn("netmask", "Netmask", 18),
		table.NewColumn("iface", "Interface", 12),
		table.NewColumn("flags", "Flags", 10),
	}
}
