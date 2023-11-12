package reader

import (
	"log/slog"
	"net"
	"strings"
	"time"

	"crdx.org/db"
	"crdx.org/lighthouse/constants"
	"crdx.org/lighthouse/m"
	"crdx.org/lighthouse/m/repo/adapterR"
	"crdx.org/lighthouse/m/repo/deviceR"
	"crdx.org/lighthouse/m/repo/settingR"
	"crdx.org/lighthouse/util/netutil"
	"github.com/google/gopacket/layers"
)

func (self *Reader) handleARPPacket(packet *layers.ARP, ipNet *net.IPNet) {
	// We are interested in the sender's IP<->MAC mapping in both requests and replies.
	if packet.Operation != layers.ARPReply && packet.Operation != layers.ARPRequest {
		return
	}

	if len(packet.SourceProtAddress) != 4 || len(packet.SourceHwAddress) != 6 {
		return
	}

	ipAddress := net.IP(packet.SourceProtAddress)
	macAddress := net.HardwareAddr(packet.SourceHwAddress)

	invalid := ipAddress.IsUnspecified() ||
		ipAddress.IsLinkLocalUnicast() ||
		ipAddress.IsLinkLocalMulticast()

	if invalid {
		return
	}

	ipAddressStr := ipAddress.String()
	macAddressStr := macAddress.String()

	if self.macAddressCache.SeenWithinLast(macAddressStr, 10*time.Second) {
		return
	}

	self.handleARP(
		strings.ToUpper(macAddressStr),
		ipAddressStr,
		ipNet.IP.String(),
	)
}

func (self *Reader) handleARP(macAddress string, ipAddress, originIPAddress string) {
	adapter, adapterFound := adapterR.Upsert(macAddress, ipAddress)

	log := self.log.With(slog.Group(
		"adapter",
		"id", adapter.ID,
		"mac", adapter.MACAddress,
		"ip", adapter.IPAddress,
		"vendor", adapter.Vendor,
	))

	var device *m.Device
	hostname, hostnameFromDHCP := self.hostnameCache[macAddress]

	// If an adapter was found, then we know it must have an attached device.
	if adapterFound {
		var found bool
		if device, found = adapter.Device(); !found {
			log.Warn("existing adapter found with no associated device")
			return
		}
	} else {
		device = db.Create(&m.Device{
			State:          deviceR.StateOnline,
			StateUpdatedAt: time.Now(),
			Icon:           constants.DefaultDeviceIconClass,
		})

		adapter.Update("device_id", device.ID)

		if hostname == "" {
			if names, err := net.LookupAddr(ipAddress); err == nil && len(names) > 0 {
				hostname = netutil.UnqualifyHostname(names[0])
			}
		}
	}

	log = log.With(slog.Group("device", "id", device.ID))

	device.Update("last_seen_at", time.Now())

	if ipAddress == originIPAddress {
		db.B[m.Device]().Update("origin", false)
		device.Update("origin", true)
	}

	if device.Name == "" && hostname != "" {
		device.Update("name", hostname)
	}

	if hostname != "" {
		device.Update("hostname", hostname)
		if hostnameFromDHCP {
			device.Update("hostname_announced_at", time.Now())
		}
	}

	if !adapterFound {
		db.Create(&m.DeviceStateLog{
			DeviceID:    device.ID,
			State:       deviceR.StateOnline,
			GracePeriod: device.GracePeriod,
		})

		if settingR.Watch() {
			db.Create(&m.DeviceDiscoveryNotification{
				DeviceID: device.ID,
			})
		}

		if settingR.WatchNew() {
			device.Update("watch", true)
		}

		log.Info("new device has joined the network")
	}
}
