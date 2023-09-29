package scanner

type networkMessage interface {
	IsNetworkMessage()
}

type arpMessage struct {
	MACAddress string
	IPAddress  string
}

func (arpMessage) IsNetworkMessage() {}

type dhcpMessage struct {
	MACAddress string
	Hostname   string
}

func (dhcpMessage) IsNetworkMessage() {}
