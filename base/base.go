package base

import (
	"fmt"

	"github.com/google/gopacket/pcap"
)

var (
	Monitor *pcap.Handle
	Devices []pcap.Interface
)

func Get_divices() {
	var err error
	Devices, err = pcap.FindAllDevs()
	if err != nil {
		fmt.Printf("Error finding devices: %v", err)
	}
}

func Set_monitor(index uint8) {
	device := Devices[index]
	var err error
	fmt.Printf("Opening device: %s\n", device.Name)
	Monitor, err = pcap.OpenLive(device.Name, 262144, true, pcap.BlockForever)
	if err != nil {
		fmt.Printf("Error opening device: %v", err)
	}
}
