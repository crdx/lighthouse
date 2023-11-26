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
	"crdx.org/lighthouse/m/repo/settingR"
	"crdx.org/lighthouse/pkg/util/dbutil"
	"crdx.org/lighthouse/pkg/util/netutil"
)

func (self *Reader) handleIncoming(macAddress string, ipAddress string, isOrigin bool) {
	macAddress = strings.ToUpper(macAddress)

	if repeaters, ok := netutil.ParseMACList(settingR.SourceMACAddresses()); ok {
		if slices.Contains(repeaters, macAddress) {
			mappings := dbutil.MapBy2[string, string]("IPAddress", "MACAddress", db.B[m.Mapping]().Find())
			if value, ok := mappings[ipAddress]; ok {
				macAddress = value
			}
		}
	}

	if self.macAddressCache.SeenWithinLast(macAddress, 10*time.Second) {
		return
	}

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

	if isOrigin {
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

		if settingR.NotifyOnNewDevice() {
			db.Create(&m.DeviceDiscoveryNotification{
				DeviceID: device.ID,
			})
		}

		if settingR.WatchNew() {
			device.Update("watch", true)
		}

		if settingR.PingNew() {
			device.Update("ping", true)
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
