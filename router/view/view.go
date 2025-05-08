package view

import (
	"log"

	b "github.com/JaredMijailRE/IEEE802Sniffer/base"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

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

func wsAnalizer(c *websocket.Conn) {

	if b.Monitor == "" {
		log.Println("Monitor not initialized")
		return
	}

	monitor, err := pcap.OpenLive(b.Monitor, 262144, true, pcap.BlockForever)
	if err != nil {
		log.Printf("error creating the monitor %s", err)
	}
	defer monitor.Close()

	packetSource := gopacket.NewPacketSource(monitor, monitor.LinkType())

	for packet := range packetSource.Packets() {
		info := b.FrameInfo{
			Timestamp:  packet.Metadata().Timestamp,
			PayloadLen: len(packet.Data()),
		}

		if rt := packet.Layer(layers.LayerTypeRadioTap); rt != nil {
			r := rt.(*layers.RadioTap)
			info.Radiotap = &b.RadiotapInfo{
				ChannelFreq:   uint16(r.ChannelFrequency),
				ChannelFlags:  uint16(r.ChannelFlags),
				DataRate:      uint8(r.Rate),
				Flags:         uint8(r.Flags),
				DBMAntennaSig: r.DBMAntennaSignal,
				Antenna:       uint8(r.Antenna),
			}
		}

		// 802.11
		if dot11 := packet.Layer(layers.LayerTypeDot11); dot11 != nil {
			d := dot11.(*layers.Dot11)
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
		}

		// Ethernet (802.3)
		if eth := packet.Layer(layers.LayerTypeEthernet); eth != nil {
			e := eth.(*layers.Ethernet)
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
		}

		// LLC
		if llc := packet.Layer(layers.LayerTypeLLC); llc != nil {
			l := llc.(*layers.LLC)
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
