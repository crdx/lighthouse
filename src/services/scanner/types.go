package scanner

type networkRessage interface {
	IsnetworkRessage()
}

type arpMessage struct {
	MACAddress string
	IPAddress  string
}

func (arpMessage) IsnetworkRessage() {}

type dhcpMessage struct {
	MACAddress string
	Hostname   string
}

func (dhcpMessage) IsnetworkRessage() {}
