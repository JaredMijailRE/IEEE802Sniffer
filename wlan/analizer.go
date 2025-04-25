package wlan

import (
	"fmt"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

// CapturaWLAN inicia la captura de tramas IEEE 802.11
func CapturaWLAN() {
	if WlanHandle == nil {
		fmt.Println("❌ WlanHandle no inicializado. Ejecuta InitWlan() primero.")
		return
	}

	// Crear un packet source desde el handler
	packetSource := gopacket.NewPacketSource(WlanHandle, WlanHandle.LinkType())

	fmt.Println("📡 Iniciando captura de tramas IEEE 802.11 (modo monitor)...")
	for packet := range packetSource.Packets() {
		// Obtener la capa Dot11
		dot11Layer := packet.Layer(layers.LayerTypeDot11)
		if dot11Layer == nil {
			continue
		}
		frame, _ := dot11Layer.(*layers.Dot11)

		// Mostrar información básica de la trama
		fmt.Println("──────────── Nueva Trama IEEE 802.11 ────────────")
		fmt.Println("  ➤ Dirección 1:", frame.Address1)
		fmt.Println("  ➤ Dirección 2:", frame.Address2)
		fmt.Println("  ➤ Dirección 3:", frame.Address3)

		// Tipo de la trama
		switch frame.Type {
		case layers.Dot11TypeMgmt:
			fmt.Println("  🔹 Tipo: Gestión")
		case layers.Dot11TypeCtrl:
			fmt.Println("  🔹 Tipo: Control")
		case layers.Dot11TypeData:
			fmt.Println("  🔹 Tipo: Datos")
		default:
			fmt.Println("  🔹 Tipo: Desconocido")
		}

		// Evaluar seguridad (bit de protección en Flags)
		if frame.Flags&0x40 != 0 {
			fmt.Println("  🔐 Seguridad: Activada (WEP/WPA/WPA2)")
		} else {
			fmt.Println("  🔓 Seguridad: No protegida")
		}

		// Evaluar Calidad de Servicio (QoS) basándose en el subtipo

		fmt.Println("─────────────────────────────────────────────")
	}
}
