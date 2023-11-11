// Package scanner scans the local network for devices.
package scanner

import (
	"errors"
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
	"crdx.org/lighthouse/pkg/cache"
	"crdx.org/lighthouse/services"
	"crdx.org/lighthouse/util/netutil"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/samber/lo"
)

// Interval between ARP request packets. 50ms means on a /24 a full scan takes about 10 seconds.
const arpPacketInterval = 50 * time.Millisecond

type Scanner struct {
	log             *slog.Logger
	macAddressCache *cache.TemporalCache[string]

	// Since the DHCP handshake happens before the device officially joins the network and starts
	// responding to ARP requests, hostnameCache stores a mapping of MAC addresses to hostnames.
	// This is then used as a lookup for a device hostname when we discover a new device on the
	// network.
	hostnameCache map[string]string
}

func New() *Scanner {
	return &Scanner{
		macAddressCache: cache.NewTemporal[string](),
		hostnameCache:   map[string]string{},
	}
}

func (self *Scanner) Init(args *services.Args) error {
	self.log = args.Logger
	return nil
}

func (self *Scanner) Run() error {
	iface, found := findInterface()
	if !found {
		return errors.New("no interface found")
	}

	ipNet, found := findIPNet(iface)
	if !found {
		return errors.New("no network found")
	} else if netutil.IPNetTooLarge(ipNet) {
		return errors.New("network too large")
	}

	networkMessages := make(chan networkMessage)

	go func() {
		if err := self.scan(iface, ipNet, networkMessages); err != nil {
			// None of the errors generated so far are recoverable, but if one is returned by scan
			// then it's likely something intermittent like a network write failure, so panic here
			// to allow the recovery handler to deal with it.
			panic(err)
		}
	}()

	for {
		networkMessage := <-networkMessages

		switch message := networkMessage.(type) {
		case arpMessage:
			self.handleARPMessage(strings.ToUpper(message.MACAddress), message.IPAddress, ipNet.IP.String())
		case dhcpMessage:
			self.handleDHCPMessage(strings.ToUpper(message.MACAddress), message.Hostname)
		}
	}
}

func (self *Scanner) scan(iface *net.Interface, ipNet *net.IPNet, messages chan<- networkMessage) error {
	handle, err := pcap.OpenLive(iface.Name, 65536, true, pcap.BlockForever)
	if err != nil {
		return err
	}

	defer handle.Close()

	stop := make(chan struct{})
	go self.read(handle, stop, messages)
	defer close(stop)

	for {
		if !settingR.Passive() {
			if err := self.write(handle, iface, ipNet); err != nil {
				return err
			}
		}

		time.Sleep(settingR.ScanInterval())
	}
}

func (*Scanner) write(handle *pcap.Handle, iface *net.Interface, ipNet *net.IPNet) error {
	ethernetLayer := layers.Ethernet{
		SrcMAC:       iface.HardwareAddr,
		DstMAC:       net.HardwareAddr{0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
		EthernetType: layers.EthernetTypeARP,
	}

	arpLayer := layers.ARP{
		AddrType:          layers.LinkTypeEthernet,
		Protocol:          layers.EthernetTypeIPv4,
		HwAddressSize:     6,
		ProtAddressSize:   4,
		Operation:         layers.ARPRequest,
		SourceHwAddress:   []byte(iface.HardwareAddr),
		SourceProtAddress: []byte(ipNet.IP),
		DstHwAddress:      []byte{0, 0, 0, 0, 0, 0},
	}

	buffer := gopacket.NewSerializeBuffer()
	options := gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}

	scan := func() error {
		for _, ip := range netutil.ExpandIPNet(ipNet) {
			arpLayer.DstProtAddress = []byte(ip)
			lo.Must0(gopacket.SerializeLayers(buffer, options, &ethernetLayer, &arpLayer))

			if err := handle.WritePacketData(buffer.Bytes()); err != nil {
				return err
			}

			time.Sleep(arpPacketInterval)
		}
		return nil
	}

	return errors.Join(scan(), scan())
}

func (self *Scanner) read(handle *pcap.Handle, stop chan struct{}, messages chan<- networkMessage) {
	packets := gopacket.NewPacketSource(handle, layers.LayerTypeEthernet).Packets()

	for {
		var packet gopacket.Packet

		select {
		case <-stop:
			return
		case packet = <-packets:
			if dhcpLayer := packet.Layer(layers.LayerTypeDHCPv4); dhcpLayer != nil {
				self.handleDHCPPacket(dhcpLayer.(*layers.DHCPv4), messages)
			} else if arpLayer := packet.Layer(layers.LayerTypeARP); arpLayer != nil {
				self.handleARPPacket(arpLayer.(*layers.ARP), messages)
			}
		}
	}
}

func (*Scanner) handleDHCPPacket(packet *layers.DHCPv4, messages chan<- networkMessage) {
	for _, option := range packet.Options {
		if option.Type == layers.DHCPOptHostname {
			messages <- dhcpMessage{
				MACAddress: packet.ClientHWAddr.String(),
				Hostname:   string(option.Data),
			}
		}
	}
}

func (self *Scanner) handleARPPacket(packet *layers.ARP, messages chan<- networkMessage) {
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

	messages <- arpMessage{
		IPAddress:  ipAddressStr,
		MACAddress: macAddressStr,
	}
}

func (self *Scanner) handleDHCPMessage(macAddress string, hostname string) {
	self.hostnameCache[macAddress] = hostname
	updateHostname(macAddress, hostname)
}

func (self *Scanner) handleARPMessage(macAddress string, ipAddress, originIPAddress string) {
	adapter, adapterFound := adapterR.Upsert(macAddress, ipAddress)

	log := self.log.With(slog.Group(
		"adapter",
		"id", adapter.ID,
		"mac", adapter.MACAddress,
		"ip", adapter.IPAddress,
		"vendor", adapter.Vendor,
	))

	var device *m.Device
	hostname := self.hostnameCache[macAddress]

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

	populateDeviceName(device, hostname)

	if hostname != "" {
		device.Update("hostname", hostname)
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

func updateHostname(macAddress string, hostname string) {
	if adapter, found := adapterR.FindByMACAddress(macAddress); found {
		if device, found := db.First[m.Device](adapter.DeviceID); found {
			device.Update("hostname", hostname)
			device.Update("hostname_announced_at", time.Now())
		}
	}
}

func populateDeviceName(device *m.Device, hostname string) {
	if device.Name != "" {
		return
	}

	if hostname != "" {
		device.Update("name", hostname)
		return
	}
}

// findInterface returns the first interface that is considered "up" and is not the loopback
// interface.
func findInterface() (*net.Interface, bool) {
	interfaces := lo.Must(net.Interfaces())

	for _, iface := range interfaces {
		isRunning := iface.Flags&net.FlagRunning == net.FlagRunning
		hasAssignedIps := lo.Must(iface.Addrs()) != nil

		if iface.Name != "lo" && isRunning && hasAssignedIps {
			return &iface, true
		}
	}

	return nil, false
}

// findIPNet returns the first IPv4 network for an interface.
func findIPNet(iface *net.Interface) (*net.IPNet, bool) {
	addresses := lo.Must(iface.Addrs())

	for _, address := range addresses {
		if network, ok := address.(*net.IPNet); ok {
			if ip := network.IP.To4(); ip != nil {
				return &net.IPNet{
					IP:   ip,
					Mask: network.Mask[len(network.Mask)-4:],
				}, true
			}
		}
	}

	return nil, false
}
