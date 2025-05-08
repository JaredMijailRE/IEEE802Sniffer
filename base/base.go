package base

import (
	"fmt"

	"github.com/google/gopacket/pcap"
)

var (
	Monitor string
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
	fmt.Printf("Opening device: %s\n", device.Name)
	Monitor = device.Name
}
