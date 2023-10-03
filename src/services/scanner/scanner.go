// Package scanner scans the local network for devices.
package scanner

import (
	"errors"
	"net"
	"time"

	"log/slog"

	"crdx.org/db"
	"crdx.org/lighthouse/m"
	"crdx.org/lighthouse/models/adapterM"
	"crdx.org/lighthouse/models/deviceM"
	"crdx.org/lighthouse/models/networkM"
	"crdx.org/lighthouse/pkg/cache"
	"crdx.org/lighthouse/services"
	"crdx.org/lighthouse/util"

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
	iface, found := findInterface()
	if !found {
		return errors.New("no interface found")
	}

	ipNet, found := findIPNet(iface)
	if !found {
		return errors.New("no network found")
	} else if ipNetTooLarge(ipNet) {
		return errors.New("network too large")
	}

	// Convert e.g. 192.168.1.20/24 to 192.168.1.0/24.
	_, generalIPNet := lo.Must2(net.ParseCIDR(ipNet.String()))

	network, _ := networkM.Upsert(generalIPNet.String())

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
			self.handleARPMessage(network, message.MACAddress, message.IPAddress)
		case dhcpMessage:
			self.handleDHCPMessage(message.MACAddress, message.Hostname)
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

	for _, ip := range expandIPNet(ipNet) {
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
	ipAddressStr := ipAddress.String()

	macAddress := net.HardwareAddr(packet.SourceHwAddress)
	macAddressStr := macAddress.String()

	if self.macAddressCache.SeenWithinLast(macAddressStr, 10*time.Second) {
		return
	}

	if ipAddressStr == "0.0.0.0" {
		return
	}

	messages <- arpMessage{
		IPAddress:  ipAddressStr,
		MACAddress: macAddressStr,
	}
}

func (self *Scanner) handleDHCPMessage(macAddress string, hostname string) {
	self.log.Info(
		"a device has broadcast its hostname",
		"mac", macAddress,
		"hostname", hostname,
	)

	self.hostnameCache[macAddress] = hostname

	updateHostname(macAddress, hostname)
}

func (self *Scanner) handleARPMessage(network *m.Network, macAddress string, ipAddress string) {
	adapter, adapterFound := adapterM.Upsert(macAddress, ipAddress)

	log := self.log.With(
		"adapter_id", adapter.ID,
		"mac", adapter.MACAddress,
		"ip", adapter.IPAddress,
		"vendor", adapter.Vendor,
	)

	var device *m.Device
	hostname := self.hostnameCache[macAddress]

	// If an adapter was found, then we know it must have an attached device.
	if adapterFound {
		var deviceFound bool
		device, deviceFound = deviceM.For(adapter.DeviceID).First()
		if !deviceFound {
			log.Warn("existing adapter found with no associated device")
			return
		}
	} else {
		device = db.Create(&m.Device{
			NetworkID: network.ID,
			State:     deviceM.StateOnline,
		})

		adapter.Update("device_id", device.ID)
	}

	device.Update("last_seen", time.Now())
	populateDeviceName(device, hostname)

	if hostname != "" {
		device.Update("hostname", hostname)
	}

	if !adapterFound {
		log.Info(
			"a new device joined the network",
			"device_id", device.ID,
			"name", device.Name,
		)
	}
}

func updateHostname(macAddress string, hostname string) {
	if adapter, found := db.B(m.Adapter{MACAddress: macAddress}).First(); found {
		if device, found := deviceM.For(adapter.DeviceID).First(); found {
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

	for _, adapter := range device.Adapters() {
		if adapter.Vendor != "" && adapter.Vendor != constants.UnknownVendorLabel {
			device.Update("name", adapter.Vendor)
			return
		}
	}
}
