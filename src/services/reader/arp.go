package reader

import (
	"net"
	"slices"
	"strings"
	"time"

	"crdx.org/db"
	"crdx.org/lighthouse/constants"
	"crdx.org/lighthouse/m"
	"crdx.org/lighthouse/m/repo/deviceR"
	"crdx.org/lighthouse/m/repo/mappingR"
	"crdx.org/lighthouse/m/repo/settingR"
	"crdx.org/lighthouse/pkg/util/netutil"
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
	macAddressStr := strings.ToUpper(macAddress.String())

	if self.macAddressCache.SeenWithinLast(macAddressStr, 10*time.Second) {
		return
	}

	if repeaters, ok := netutil.ParseMACList(settingR.SourceMACAddresses()); ok {
		if slices.Contains(repeaters, macAddressStr) {
			mappings := mappingR.Map()
			if value, ok := mappings[ipAddressStr]; ok {
				macAddressStr = value
			}
		}
	}

	self.handleARP(
		macAddressStr,
		ipAddressStr,
		ipNet.IP.String(),
	)
}

func (self *Reader) handleARP(macAddress string, ipAddress, originIPAddress string) {
	hostname, hostnameFromDHCP := self.hostnameCache[macAddress]
	device, adapter, found := findOrCreate(macAddress, ipAddress)

	if !found {
		if vendor, vendorFound := netutil.GetVendor(adapter.MACAddress); vendorFound {
			adapter.Update("vendor", vendor)
		} else if !db.B[m.VendorLookup]("adapter_id = ?", adapter.ID).Exists() {
			db.Create(&m.VendorLookup{AdapterID: adapter.ID})
		}

		if hostname == "" {
			if names, err := net.LookupAddr(ipAddress); err == nil && len(names) > 0 {
				hostname = netutil.UnqualifyHostname(names[0])
			}
		}
	}

	device.Update("last_seen_at", time.Now())
	adapter.Update("last_seen_at", time.Now())

	if ipAddress == originIPAddress {
		db.B[m.Device]().Update("origin", false)
		device.Update("origin", true)
	}

	if hostname != "" {
		device.Update("hostname", hostname)

		if device.Name == "" {
			device.Update("name", hostname)
		}

		if hostnameFromDHCP {
			device.Update("hostname_announced_at", time.Now())
		}
	}

	if !found {
		device = device.Fresh()
		adapter = adapter.Fresh()

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

		self.log.With(
			device.LogAttr(),
			adapter.LogAttr(),
		).Info("new device has joined the network")
	}
}

func findOrCreate(macAddress string, ipAddress string) (*m.Device, *m.Adapter, bool) {
	adapter, found := db.B[m.Adapter]("mac_address = ?", macAddress).First()
	var device *m.Device

	if found {
		device = adapter.Device()
		if adapter.IPAddress != ipAddress {
			adapter.Update("ip_address", ipAddress)
		}
	} else {
		device = db.Create(&m.Device{
			State:          deviceR.StateOnline,
			StateUpdatedAt: time.Now(),
			Icon:           constants.DefaultDeviceIconClass,
			GracePeriod:    constants.DefaultGracePeriod,
			LastSeenAt:     time.Now(),
		})

		adapter = db.Create(&m.Adapter{
			DeviceID:   device.ID,
			MACAddress: macAddress,
			IPAddress:  ipAddress,
		})
	}

	return device, adapter, found
}
