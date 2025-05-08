package view

import (
	"log"
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

	if b.Monitor == "" {
		log.Println("Monitor not initialized")
		c.Close()
		return
	}

	monitor, err := pcap.OpenLive(b.Monitor, 262144, true, pcap.BlockForever)
	if err != nil {
		log.Printf("Error creating the monitor %s: %v", b.Monitor, err)
		c.Close()
		return
	}
	defer monitor.Close()

	log.Printf("Sniffing on interface %s with link type: %s", b.Monitor, monitor.LinkType().String())

	packetSource := gopacket.NewPacketSource(monitor, monitor.LinkType())

	for packet := range packetSource.Packets() {
		info := b.FrameInfo{
			Timestamp:  packet.Metadata().Timestamp,
			PayloadLen: len(packet.Data()),
		}

		// Layer presence flags
		// hasRadiotap := false // Removed as it's not used elsewhere for now
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
			// hasRadiotap = true // Not strictly needed for current logic if only used for this log block
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

		// Diagnostic log if WiFi interface shows Ethernet but not Dot11
		isLikelyWifiInterface := strings.Contains(b.Monitor, "wlan") || strings.Contains(b.Monitor, "wlp") || strings.Contains(b.Monitor, "ath") || strings.Contains(b.Monitor, "wifi")
		if isLikelyWifiInterface && hasEthernet && !hasDot11 {
			log.Printf("DEBUG: WiFi interface (%s) packet seen as Ethernet without Dot11. LinkType: %s", b.Monitor, monitor.LinkType().String())
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
			log.Println("Error enviando JSON:", err)
			return
		}
	}
}
