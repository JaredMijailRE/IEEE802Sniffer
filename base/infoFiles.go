package base

import "time"

type FrameInfo struct {
	Timestamp  time.Time     `json:"timestamp"`
	PayloadLen int           `json:"payload_length"`
	Radiotap   *RadiotapInfo `json:"radiotap,omitempty"`
	Dot11      *Dot11Info    `json:"dot11,omitempty"`
	Ethernet   *EthernetInfo `json:"ethernet,omitempty"`
	LLC        *LLCInfo      `json:"llc,omitempty"`
}

type RadiotapInfo struct {
	ChannelFreq   uint16 `json:"channel_freq"`
	ChannelFlags  uint16 `json:"channel_flags"`
	DataRate      uint8  `json:"data_rate"`
	Flags         uint8  `json:"flags"`
	DBMAntennaSig int8   `json:"dbm_antenna_signal"`
	Antenna       uint8  `json:"antenna"`
}

type Dot11Info struct {
	Version   uint8  `json:"version"`
	Type      string `json:"type"` // Management / Control / Data
	Subtype   string `json:"subtype"`
	ToDS      bool   `json:"to_ds"`
	FromDS    bool   `json:"from_ds"`
	MoreFrag  bool   `json:"more_frag"`
	Retry     bool   `json:"retry"`
	PwrMgmt   bool   `json:"pwr_mgmt"`
	MoreData  bool   `json:"more_data"`
	Protected bool   `json:"protected"`
	Order     bool   `json:"order"`
	Addr1     string `json:"addr1"`
	Addr2     string `json:"addr2"`
	Addr3     string `json:"addr3"`
	Addr4     string `json:"addr4,omitempty"`
	Sequence  uint16 `json:"sequence_control"`
	Fragment  uint16 `json:"fragment_number"`
	QoS       *struct {
		TID       uint8  `json:"tid"`
		EOSP      bool   `json:"eosp"`
		AckPolicy uint8  `json:"ack_policy"`
		TxOP      uint16 `json:"txop_limit"`
	} `json:"qos,omitempty"`
}

type EthernetInfo struct {
	SrcMAC       string  `json:"src_mac"`
	DstMAC       string  `json:"dst_mac"`
	EthernetType string  `json:"eth_type"`
	VLAN         *uint16 `json:"vlan_id,omitempty"`
}

type LLCInfo struct {
	DSAP         uint8  `json:"dsap"`
	SSAP         uint8  `json:"ssap"`
	Control      uint8  `json:"control"`
	Organization []byte `json:"org_code,omitempty"`
	EtherType    uint16 `json:"ethertype,omitempty"`
}
