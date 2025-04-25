package lan

import (
	"fmt"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

// CapturaLAN inicia la captura de tramas IEEE 802.3 desde la interfaz seleccionada
func CapturaLAN() {
	if LanHandle == nil {
		fmt.Println("❌ LanHandle no está inicializado. Asegúrate de llamar a Set_monitor_Lan primero.")
		return
	}

	// Crear un packet source para leer paquetes
	packetSource := gopacket.NewPacketSource(LanHandle, LanHandle.LinkType())

	fmt.Println("🔍 Iniciando captura de tramas IEEE 802.3 (Ethernet)...")
	for packet := range packetSource.Packets() {
		ethernetLayer := packet.Layer(layers.LayerTypeEthernet)
		if ethernetLayer != nil {
			ethernetPacket, _ := ethernetLayer.(*layers.Ethernet)

			fmt.Println("📦 Nueva trama capturada:")
			fmt.Println("  ➤ Dirección MAC Origen:     ", ethernetPacket.SrcMAC)
			fmt.Println("  ➤ Dirección MAC Destino:    ", ethernetPacket.DstMAC)
			fmt.Printf("  ➤ Tipo de protocolo (Ethertype): 0x%04x\n", uint16(ethernetPacket.EthernetType))

			// Podrías inspeccionar más a fondo según el tipo (por ejemplo IPv4, ARP, etc.)
			fmt.Println("─────────────────────────────────────────────")
		}
	}
}
