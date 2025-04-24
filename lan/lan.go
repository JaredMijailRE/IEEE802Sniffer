package lan

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/google/gopacket/pcap"
)

var LanHandle *pcap.Handle

func Set_monitor_Lan() {
	// 1. Listar interfaces
	devices, err := pcap.FindAllDevs()
	if err != nil {
		panic(err)
	}

	fmt.Println("Interfaces de red disponibles:")
	for i, device := range devices {
		fmt.Printf("%d) %s - %s\n", i, device.Name, device.Description)
	}

	reader := bufio.NewReader(os.Stdin)

	// 2. Seleccionar interfaz LAN
	fmt.Print("Selecciona el número de la interfaz LAN: ")
	lanIndex := readIndex(reader, len(devices))
	lanDevice := devices[lanIndex]

	// 3. Abrir interfaces
	fmt.Println("\nAbriendo LAN en modo promiscuo...")
	openInterface(lanDevice.Name)

}

func readIndex(reader *bufio.Reader, max int) int {
	for {
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		var index int
		_, err := fmt.Sscanf(input, "%d", &index)
		if err == nil && index >= 0 && index < max {
			return index
		}
		fmt.Print("Índice inválido. Intenta de nuevo: ")
	}
}

func openInterface(name string) {
	var err error

	LanHandle, err = pcap.OpenLive(name, 65535, true, pcap.BlockForever)
	if err != nil {
		fmt.Println("❌ Error abriendo la interfaz:", err)
		return
	}
	fmt.Println("Interfaz Lan abierta con éxito:", name)
}
