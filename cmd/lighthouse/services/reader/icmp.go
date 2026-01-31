package reader

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

func (self *Reader) handleICMPPacket(icmpPacket *layers.ICMPv4, packet gopacket.Packet, originIPAddress string) {
	if icmpPacket.TypeCode != layers.CreateICMPv4TypeCode(layers.ICMPv4TypeEchoReply, 0) {
		return
	}

	ethernetLayer := packet.Layer(layers.LayerTypeEthernet)
	ipLayer := packet.Layer(layers.LayerTypeIPv4)

	if ethernetLayer == nil || ipLayer == nil {
		return
	}

	ethernet, ok := ethernetLayer.(*layers.Ethernet)
	if !ok {
		return
	}

	ip, ok := ipLayer.(*layers.IPv4)
	if !ok {
		return
	}

	macAddress := ethernet.SrcMAC.String()
	ipAddress := ip.SrcIP.String()

	self.logger.Info("received ICMP reply", "ip", ipAddress)

	self.handleIncoming(
		macAddress,
		ipAddress,
		ipAddress == originIPAddress,
	)
}
