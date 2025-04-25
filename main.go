package main

import (
	"github.com/JaredMijailRE/IEEE802Sniffer/lan"
	"github.com/JaredMijailRE/IEEE802Sniffer/wlan"
)

func main() {
	wlan.InitWlan()
	wlan.CapturaWLAN()

	if wlan.WlanHandle == nil {
		defer lan.LanHandle.Close()
	}
	if lan.LanHandle == nil {
		defer wlan.WlanHandle.Close()
	}
}
