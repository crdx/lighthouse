package reader

import (
	"net"
	"slices"
	"strings"
	"time"

	"crdx.org/lighthouse/db"
	"crdx.org/lighthouse/db/repo/deviceR"
	"crdx.org/lighthouse/db/repo/settingR"
	"crdx.org/lighthouse/pkg/constants"
	"crdx.org/lighthouse/pkg/util/netutil"
)

func (self *Reader) handleIncoming(macAddress string, ipAddress string, isOrigin bool) {
	macAddress = strings.ToUpper(macAddress)

	if repeaters, ok := netutil.ParseMACList(settingR.SourceMACAddresses()); ok {
		if slices.Contains(repeaters, macAddress) {
			mappings := db.MapBy2[string, string]("IPAddress", "MACAddress", db.FindMappings())
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
			adapter.UpdateVendor(vendor)
		} else if _, found := db.FindVendorLookupByAdapterID(adapter.ID); !found {
			db.CreateVendorLookup(&db.VendorLookup{AdapterID: adapter.ID})
		}

		if hostname == "" {
			if names, err := net.LookupAddr(ipAddress); err == nil && len(names) > 0 {
				hostname = netutil.UnqualifyHostname(names[0])
			}
		}
	}

	ipChanged := adapter.IPAddress != ipAddress

	if ipChanged && found {
		adapter.UpdateIPAddress(ipAddress)
	}
	if ipChanged || !found {
		db.CreateDeviceIPAddressLog(&db.DeviceIPAddressLog{
			DeviceID:  device.ID,
			IPAddress: ipAddress,
		})
	}

	device.UpdateLastSeenAt(db.Now())
	adapter.UpdateLastSeenAt(db.Now())

	if isOrigin {
		db.ResetOriginDevices()
		device.UpdateOrigin(true)
	}

	if hostname != "" {
		device.UpdateHostname(hostname)

		if device.Name == "" {
			device.UpdateName(hostname)
		}

		if hostnameFromDHCP {
			device.UpdateHostnameAnnouncedAt(db.Now())
		}
	}

	if !found {
		db.CreateDeviceStateLog(&db.DeviceStateLog{
			DeviceID:    device.ID,
			State:       deviceR.StateOnline,
			GracePeriod: device.GracePeriod,
		})

		if settingR.NotifyOnNewDevice() {
			db.CreateDeviceDiscoveryNotification(&db.DeviceDiscoveryNotification{
				DeviceID: device.ID,
			})
		}

		if settingR.WatchNew() {
			device.UpdateWatch(true)
		}

		if settingR.PingNew() {
			device.UpdatePing(true)
		}

		self.logger.With(
			device.LogAttr(),
			adapter.LogAttr(),
		).Info("new device has joined the network")
	}
}

func findOrCreate(macAddress string, ipAddress string) (*db.Device, *db.Adapter, bool) {
	adapter, found := db.FindAdapterByMACAddress(macAddress)
	var device *db.Device

	if found {
		device = adapter.Device()
	} else {
		device = db.CreateDevice(&db.Device{
			State:          deviceR.StateOnline,
			StateUpdatedAt: db.Now(),
			Icon:           constants.DefaultDeviceIconClass,
			GracePeriod:    constants.DefaultGracePeriod,
			LastSeenAt:     db.Now(),
		})

		adapter = db.CreateAdapter(&db.Adapter{
			DeviceID:   device.ID,
			MACAddress: macAddress,
			IPAddress:  ipAddress,
		})
	}

	return device, adapter, found
}
