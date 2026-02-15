package arp

import "github.com/evertras/bubble-table/table"

func columns() []table.Column {
	return []table.Column{
		table.NewColumn("ip", "IP", 18).WithFiltered(true),
		table.NewColumn("mac", "MAC", 19).WithFiltered(true),
		table.NewColumn("iface", "Interface", 12).WithFiltered(true),
		table.NewFlexColumn("flags", "Flags", 1),
	}
}
