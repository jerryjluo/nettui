package firewall

import "github.com/evertras/bubble-table/table"

func columns() []table.Column {
	return []table.Column{
		table.NewColumn("rule", "Rule#", 7),
		table.NewColumn("action", "Action", 8),
		table.NewColumn("dir", "Direction", 11),
		table.NewColumn("proto", "Proto", 8),
		table.NewColumn("src", "Src", 22),
		table.NewColumn("dst", "Dst", 22),
		table.NewColumn("packets", "Packets", 10),
		table.NewColumn("bytes", "Bytes", 10),
	}
}
