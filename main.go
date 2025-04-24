package main

import (
	"github.com/JaredMijailRE/IEEE802Sniffer/lan"
)

func main() {
	lan.Set_monitor_Lan()

	defer lan.LanHandle.Close()
}
