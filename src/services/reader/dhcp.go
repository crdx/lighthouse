package reader

import (
	"strings"
	"time"

	"crdx.org/db"
	"crdx.org/lighthouse/m"
	"crdx.org/lighthouse/m/repo/adapterR"
	"github.com/google/gopacket/layers"
)

func (self *Reader) handleDHCPPacket(packet *layers.DHCPv4) {
	for _, option := range packet.Options {
		if option.Type == layers.DHCPOptHostname {
			self.handleDHCP(
				strings.ToUpper(packet.ClientHWAddr.String()),
				string(option.Data),
			)
		}
	}
}

func (self *Reader) handleDHCP(macAddress string, hostname string) {
	self.hostnameCache[macAddress] = hostname

	if adapter, found := adapterR.FindByMACAddress(macAddress); found {
		if device, found := db.First[m.Device](adapter.DeviceID); found {
			device.Update("hostname", hostname)
			device.Update("hostname_announced_at", time.Now())
		}
	}
}
