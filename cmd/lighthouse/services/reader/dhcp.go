package reader

import (
	"strings"

	"crdx.org/lighthouse/db"
	"github.com/google/gopacket/layers"
)

func (self *Reader) handleDHCPPacket(dhcpPacket *layers.DHCPv4) {
	for _, option := range dhcpPacket.Options {
		if option.Type == layers.DHCPOptHostname {
			self.handleDHCP(
				strings.ToUpper(dhcpPacket.ClientHWAddr.String()),
				string(option.Data),
			)
		}
	}
}

func (self *Reader) handleDHCP(macAddress string, hostname string) {
	self.hostnameCache[macAddress] = hostname

	if adapter, found := db.FindAdapterByMACAddress(macAddress); found {
		if device, found := db.FindDevice(adapter.DeviceID); found {
			device.UpdateHostname(hostname)
			device.UpdateHostnameAnnouncedAt(db.Now())
		}
	}
}
