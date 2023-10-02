// Package scanner scans the local network for devices.
package scanner

import (
	"errors"
	"net"
	"time"

	"log/slog"

	"crdx.org/lighthouse/m"
	"crdx.org/lighthouse/models/deviceMappingModel"
	"crdx.org/lighthouse/models/deviceModel"
	"crdx.org/lighthouse/models/networkModel"
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

	network, _ := networkModel.Upsert(generalIPNet.String())

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

		util.Sleep(scanInterval)
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

	deviceModel.UpdateHostname(macAddress, hostname)
}

func (self *Scanner) handleARPMessage(network m.Network, macAddress string, ipAddress string) {
	device, _ := deviceModel.Upsert(network.ID, macAddress, self.hostnameCache[macAddress])

	if _, found := deviceMappingModel.Upsert(device.ID, ipAddress); !found {
		self.log.Info(
			"a new device joined the network",
			"device_id", device.ID,
			"name", device.Name,
			"vendor", device.Vendor,
			"mac", macAddress,
			"ip", ipAddress,
		)
	}
}
