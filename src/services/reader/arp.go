package reader

import (
	"net"

	"github.com/google/gopacket/layers"
)

func (self *Reader) handleARPPacket(arpPacket *layers.ARP, originIPAddress string) {
	// We are interested in the sender's IP<->MAC mapping in both requests and replies.
	if arpPacket.Operation != layers.ARPReply && arpPacket.Operation != layers.ARPRequest {
		return
	}

	if len(arpPacket.SourceProtAddress) != 4 || len(arpPacket.SourceHwAddress) != 6 {
		return
	}

	ipAddress := net.IP(arpPacket.SourceProtAddress)
	macAddress := net.HardwareAddr(arpPacket.SourceHwAddress)

	invalid := ipAddress.IsUnspecified() ||
		ipAddress.IsLinkLocalUnicast() ||
		ipAddress.IsLinkLocalMulticast()

	if invalid {
		return
	}

	self.handleIncoming(
		macAddress.String(),
		ipAddress.String(),
		ipAddress.String() == originIPAddress,
	)
}
