package wlan

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/gopacket/pcap"
)

var WlanHandle *pcap.Handle

func InitWlan() {
	// Listar todas las interfaces de red
	devices, err := pcap.FindAllDevs()
	if err != nil {
		fmt.Printf("Error obteniendo interfaces: %v\n", err)
		return
	}

	// Mostrar interfaces con sugerencias para WLAN
	fmt.Println("Interfaces de red disponibles:")
	for i, device := range devices {
		wlanHint := ""
		if isWLANInterface(device) {
			wlanHint = " (posible WLAN)"
		}
		fmt.Printf("%2d) %-15s %s%s\n", i, device.Name, device.Description, wlanHint)
	}

	reader := bufio.NewReader(os.Stdin)

	// Seleccionar interfaz WLAN
	fmt.Print("\nSelecciona el número de la interfaz WLAN: ")
	wlanIndex := readIndex(reader, len(devices))
	selectedDevice := devices[wlanIndex]

	// Abrir interfaz seleccionada
	fmt.Printf("\nAbriendo interfaz WLAN: %s...\n", selectedDevice.Description)
	openWlanInterface(selectedDevice.Name)
}

func isWLANInterface(device pcap.Interface) bool {
	// Detección multiplataforma de interfaces inalámbricas
	name := strings.ToLower(device.Name)
	desc := strings.ToLower(device.Description)

	wlanKeywords := []string{
		"wireless", "wi-fi", "wlan", "802.11",
		"airport", "radio", "wi\\s?fi", "wifi",
	}

	for _, kw := range wlanKeywords {
		if strings.Contains(desc, kw) || strings.Contains(name, kw) {
			return true
		}
	}
	return false
}

func readIndex(reader *bufio.Reader, max int) int {
	for {
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		var index int
		if _, err := fmt.Sscanf(input, "%d", &index); err == nil {
			if index >= 0 && index < max {
				return index
			}
		}
		fmt.Print("Selección inválida. Intenta de nuevo: ")
	}
}

func openWlanInterface(interfaceName string) {
	var err error

	// Parámetros para captura básica (no promiscuo, timeout bajo)
	WlanHandle, err = pcap.OpenLive(
		interfaceName,
		1600,           // snapshot length
		false,          // promiscuous mode
		-1*time.Second, // timeout
	)

	if err != nil {
		fmt.Printf("Error abriendo interfaz: %v\n", err)

		// Mensajes específicos para errores comunes
		if strings.Contains(err.Error(), "permissions") {
			fmt.Println("\nSugerencia: Ejecuta con permisos de administrador/root")
		}
		return
	}

	// Configurar filtros básicos (opcional)
	fmt.Printf("Interfaz %s abierta exitosamente\n", interfaceName)
}
