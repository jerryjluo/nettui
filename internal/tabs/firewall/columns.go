package firewall

import "github.com/evertras/bubble-table/table"

func columns() []table.Column {
	return []table.Column{
		table.NewColumn("rule", "Rule#", 7),
		table.NewColumn("action", "Action", 8).WithFiltered(true),
		table.NewColumn("dir", "Direction", 11).WithFiltered(true),
		table.NewColumn("proto", "Proto", 8),
		table.NewFlexColumn("src", "Src", 1).WithFiltered(true),
		table.NewFlexColumn("dst", "Dst", 1).WithFiltered(true),
		table.NewColumn("packets", "Packets", 10),
		table.NewColumn("bytes", "Bytes", 10),
	}
}
