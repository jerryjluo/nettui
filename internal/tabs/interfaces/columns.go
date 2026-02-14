package interfaces

import "github.com/evertras/bubble-table/table"

func columns() []table.Column {
	return []table.Column{
		table.NewColumn("name", "Name", 10).WithFiltered(true),
		table.NewColumn("addrs", "IPs", 22).WithFiltered(true),
		table.NewColumn("mac", "MAC", 19),
		table.NewColumn("mtu", "MTU", 7),
		table.NewColumn("status", "Status", 8),
		table.NewColumn("tx_bytes", "TX", 10),
		table.NewColumn("rx_bytes", "RX", 10),
		table.NewColumn("tx_rate", "TX Rate", 12),
		table.NewColumn("rx_rate", "RX Rate", 12),
	}
}
