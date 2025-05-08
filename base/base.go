package base

import (
	"fmt"

	"github.com/google/gopacket/pcap"
)

var (
	Monitor string
	Devices []pcap.Interface
)

// consults and saves all pcap diveces available
func Get_divices() {
	var err error
	Devices, err = pcap.FindAllDevs()
	if err != nil {
		fmt.Printf("Error finding devices: %v", err)
	}
}

// stores the name of the divece to be store
func Set_monitor(index uint8) {
	device := Devices[index]
	fmt.Printf("Opening device: %s\n", device.Name)
	Monitor = device.Name
}
