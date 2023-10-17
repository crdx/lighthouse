// Package scanner scans the local network for devices.
package scanner

import (
	"errors"
	"net"
	"strings"
	"time"

	"log/slog"

	"crdx.org/db"
	"crdx.org/lighthouse/constants"
	"crdx.org/lighthouse/m"
	"crdx.org/lighthouse/pkg/cache"
	"crdx.org/lighthouse/repos/adapterR"
	"crdx.org/lighthouse/repos/deviceR"
	"crdx.org/lighthouse/repos/networkR"
	"crdx.org/lighthouse/services"
	"crdx.org/lighthouse/util/netutil"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/samber/lo"
)

// Interval between ARP request packets. 50ms means on a /24 a full scan takes about 10 seconds.
const arpPacketInterval = 50 * time.Millisecond

// Interval between scans.
const scanInterval = 60 * time.Second

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
	iface, found := netutil.FindInterface()
	if !found {
		return errors.New("no interface found")
	}

	ipNet, found := netutil.FindIPNet(iface)
	if !found {
		return errors.New("no network found")
	} else if netutil.IPNetTooLarge(ipNet) {
		return errors.New("network too large")
	}

	// Convert e.g. 192.168.1.20/24 to 192.168.1.0/24.
	_, generalIPNet := lo.Must2(net.ParseCIDR(ipNet.String()))

	network, _ := networkR.Upsert(generalIPNet.String())

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
			self.handleARPMessage(network, strings.ToUpper(message.MACAddress), message.IPAddress)
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
		if err := self.write(handle, iface, ipNet); err != nil {
			return err
		}

		time.Sleep(scanInterval)
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
	self.log.Info(
		"device has broadcast its hostname",
		"mac", macAddress,
		"hostname", hostname,
	)

	self.hostnameCache[macAddress] = hostname

	updateHostname(macAddress, hostname)
}

func (self *Scanner) handleARPMessage(network *m.Network, macAddress string, ipAddress string) {
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
			NetworkID: network.ID,
			State:     deviceR.StateOnline,
			Icon:      constants.DefaultDeviceIconClass,
		})

		adapter.Update("device_id", device.ID)

		if hostname == "" {
			if names, err := net.LookupAddr(ipAddress); err == nil && len(names) > 0 {
				hostname = netutil.UnqualifyHostname(names[0])
				log.Info("found hostname via DNS", "hostname", hostname)
			}
		}
	}

	log = log.With(slog.Group("device", "id", device.ID, "name", device.Name))

	device.Update("last_seen", time.Now())
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

		log.Info("new device has joined the network")
	}
}

func updateHostname(macAddress string, hostname string) {
	if adapter, found := adapterR.FindByMACAddress(macAddress); found {
		if device, found := db.First[m.Device](adapter.DeviceID); found {
			device.Update("hostname", hostname)
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
