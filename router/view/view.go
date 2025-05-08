package view

import (
	// For exec.Command output
	"fmt"     // For error formatting
	"log"     // For OS-specific logic
	"os/exec" // For running airmon-ng
	"regexp"  // For parsing airmon-ng output
	"runtime" // For GOOS
	"strings"

	b "github.com/JaredMijailRE/IEEE802Sniffer/base"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

// setup the view for the websocket
func SetupView(app *fiber.App) {
	// middleware para websockets
	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws/analizer/:type", websocket.New(wsAnalizer))
}

// analizer the packets and send them to the websocket
func wsAnalizer(c *websocket.Conn) {
	log.Println("New WebSocket connection for /ws/analizer")

	if b.Monitor == "" {
		log.Println("Monitor interface not initialized in config")
		_ = c.WriteMessage(websocket.TextMessage, []byte("Error: Monitor interface not configured on the server."))
		c.Close()
		return
	}

	actualInterfaceToSniff := b.Monitor
	var airmonManagedInterface string = ""

	if runtime.GOOS == "linux" {
		log.Printf("Linux detected. Attempting to manage interface %s with airmon-ng.", b.Monitor)

		// Step 1: Run airmon-ng check kill
		log.Println("Running: sudo airmon-ng check kill")
		cmdCheckKill := exec.Command("sudo", "airmon-ng", "check", "kill")
		checkKillOutput, errCheckKill := cmdCheckKill.CombinedOutput()
		if errCheckKill != nil {
			log.Printf("Warning: 'sudo airmon-ng check kill' failed: %v. Output: %s", errCheckKill, string(checkKillOutput))
		} else {
			log.Printf("'sudo airmon-ng check kill' output: %s", string(checkKillOutput))
		}

		// Step 2: Start monitor mode with airmon-ng
		log.Printf("Running: sudo airmon-ng start %s", b.Monitor)
		cmdStart := exec.Command("sudo", "airmon-ng", "start", b.Monitor)
		startOutput, errStart := cmdStart.CombinedOutput()

		if errStart != nil {
			log.Printf("Error running 'sudo airmon-ng start %s': %v. Output: %s", b.Monitor, errStart, string(startOutput))
			log.Println("Falling back to using original interface name. It might already be in monitor mode or pcap.OpenLive might succeed.")
		} else { // airmon-ng start command was successful
			log.Printf("'sudo airmon-ng start %s' command successful. Output: %s", b.Monitor, string(startOutput))
			parsedNewInterface := ""
			re := regexp.MustCompile(`(?i)(?:on|to|enabled)\s+((?:\[phy\d+\])?\w+mon)\b`)
			matches := re.FindStringSubmatch(string(startOutput))

			if len(matches) > 1 {
				parsedNewInterface = matches[1]
				log.Printf("Potential monitor interface name from airmon-ng output: %s", parsedNewInterface)

				// Remove [phyX] prefix if present (e.g. [phy0]wlan0mon -> wlan0mon)
				rePhy := regexp.MustCompile(`^\[phy\d+\]`)
				actualInterfaceToSniff = rePhy.ReplaceAllString(parsedNewInterface, "")
				airmonManagedInterface = actualInterfaceToSniff
				log.Printf("Using interface %s for sniffing (derived from %s)", actualInterfaceToSniff, parsedNewInterface)
			} else {
				log.Printf("Could not parse a new '...mon' interface name from 'airmon-ng start' output. Output was: %s", string(startOutput))
				log.Println("Assuming airmon-ng enabled monitor mode on the original interface name or the name is otherwise unchanged by 'airmon-ng start'.")
				actualInterfaceToSniff = b.Monitor
				airmonManagedInterface = b.Monitor
			}
		}
	}

	if airmonManagedInterface != "" && runtime.GOOS == "linux" {
		defer func(ifaceToStop string) {
			log.Printf("Deferred action: Attempting to stop monitor mode on interface '%s' using 'sudo airmon-ng stop %s'", ifaceToStop, ifaceToStop)
			cmdStop := exec.Command("sudo", "airmon-ng", "stop", ifaceToStop)
			stopOutput, errStop := cmdStop.CombinedOutput()
			if errStop != nil {
				log.Printf("Error running 'sudo airmon-ng stop %s': %v. Output: %s", ifaceToStop, errStop, string(stopOutput))
			} else {
				log.Printf("'sudo airmon-ng stop %s' successful. Output: %s", ifaceToStop, string(stopOutput))
			}
		}(airmonManagedInterface)
	}

	log.Printf("Attempting to open live capture on interface: %s", actualInterfaceToSniff)
	// Add a read timeout to the pcap handle
	handle, err := pcap.OpenLive(actualInterfaceToSniff, 262144, true, pcap.BlockForever) // Keep BlockForever for now, but we'll use SetReadTimeout later if needed with gopacket's PacketSource
	if err != nil {
		errMsg := fmt.Sprintf("Error creating pcap handle on interface %s: %v", actualInterfaceToSniff, err)
		log.Println(errMsg)
		_ = c.WriteMessage(websocket.TextMessage, []byte("Error: "+errMsg))
		c.Close()
		return
	}
	defer handle.Close()

	// It seems gopacket.NewPacketSource with pcap.BlockForever doesn't directly benefit from handle.SetReadTimeout.
	// Instead, the blocking nature is handled by the Packets() channel.
	// We will add logging to see if packets are being received from the source.

	log.Printf("Successfully opened interface %s for sniffing. Link type: %s", actualInterfaceToSniff, handle.LinkType().String())

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	log.Println("Packet source created. Waiting for packets...")
	packetsReceived := 0
	loopIterations := 0

	// Create a ticker to log if we are stuck waiting for packets
	// ticker := time.NewTicker(5 * time.Second) // Removed time import for now
	// defer ticker.Stop()

	for packet := range packetSource.Packets() { // This is a blocking call until a packet arrives or the source closes
		loopIterations++
		packetsReceived++
		if packetsReceived%10 == 0 { // Log every 10 packets received
			log.Printf("Processed %d packets from %s (loop iterations: %d)...", packetsReceived, actualInterfaceToSniff, loopIterations)
		}
		// If packet is nil, the source is closing or has an error.
		if packet == nil {
			log.Printf("Packet source on %s returned a nil packet. Closing capture.", actualInterfaceToSniff)
			break
		}

		info := b.FrameInfo{
			Timestamp:  packet.Metadata().Timestamp,
			PayloadLen: len(packet.Data()),
		}

		hasDot11 := false
		hasEthernet := false

		if rtLayer := packet.Layer(layers.LayerTypeRadioTap); rtLayer != nil {
			r := rtLayer.(*layers.RadioTap)
			info.Radiotap = &b.RadiotapInfo{
				ChannelFreq:   uint16(r.ChannelFrequency),
				ChannelFlags:  uint16(r.ChannelFlags),
				DataRate:      uint8(r.Rate),
				Flags:         uint8(r.Flags),
				DBMAntennaSig: r.DBMAntennaSignal,
				Antenna:       uint8(r.Antenna),
			}
		}

		if dot11Layer := packet.Layer(layers.LayerTypeDot11); dot11Layer != nil {
			d := dot11Layer.(*layers.Dot11)
			d11 := b.Dot11Info{
				Type:     d.Type.String(),
				Subtype:  d.LayerType().String(),
				Addr1:    d.Address1.String(),
				Addr2:    d.Address2.String(),
				Addr3:    d.Address3.String(),
				Sequence: d.SequenceNumber,
				Fragment: d.FragmentNumber,
			}
			if d.Address4 != nil {
				d11.Addr4 = d.Address4.String()
			}
			info.Dot11 = &d11
			hasDot11 = true
		}

		if ethLayer := packet.Layer(layers.LayerTypeEthernet); ethLayer != nil {
			e := ethLayer.(*layers.Ethernet)
			el := b.EthernetInfo{
				SrcMAC:       e.SrcMAC.String(),
				DstMAC:       e.DstMAC.String(),
				EthernetType: e.EthernetType.String(),
			}
			if vlan := packet.Layer(layers.LayerTypeDot1Q); vlan != nil {
				v := vlan.(*layers.Dot1Q)
				el.VLAN = &v.VLANIdentifier
			}
			info.Ethernet = &el
			hasEthernet = true
		}

		isLikelyWifiInterface := strings.Contains(actualInterfaceToSniff, "wlan") || strings.Contains(actualInterfaceToSniff, "wlp") || strings.Contains(actualInterfaceToSniff, "ath") || strings.Contains(actualInterfaceToSniff, "wifi") || strings.HasSuffix(actualInterfaceToSniff, "mon")
		if isLikelyWifiInterface && hasEthernet && !hasDot11 {
			log.Printf("DEBUG: WiFi-like interface (%s) packet seen as Ethernet without Dot11. Pcap LinkType: %s", actualInterfaceToSniff, handle.LinkType().String())
			var detectedLayerTypes []string
			for _, layer := range packet.Layers() {
				detectedLayerTypes = append(detectedLayerTypes, layer.LayerType().String())
			}
			log.Printf("DEBUG: Detected layers: %v", detectedLayerTypes)
			if ethLayerConv, ok := packet.Layer(layers.LayerTypeEthernet).(*layers.Ethernet); ok {
				log.Printf("DEBUG: Ethernet Type: %s, Payload Length: %d", ethLayerConv.EthernetType.String(), len(ethLayerConv.Payload)) // Corrected .Payload() to .Payload
			}
		}

		// LLC
		if llcLayer := packet.Layer(layers.LayerTypeLLC); llcLayer != nil {
			l := llcLayer.(*layers.LLC)
			ll := b.LLCInfo{
				DSAP:    l.DSAP,
				SSAP:    l.SSAP,
				Control: uint8(l.Control),
			}
			info.LLC = &ll
		}

		// Enviar JSON
		if err := c.WriteJSON(info); err != nil {
			log.Println("Error enviando JSON via WebSocket:", err)
			return // Exit goroutine for this connection
		}
	}
	log.Printf("Packet source channel closed for %s. Total packets received: %d, Loop iterations: %d", actualInterfaceToSniff, packetsReceived, loopIterations)
}
