package pinger

import (
	"log/slog"
	"net"

	"crdx.org/lighthouse/cmd/lighthouse/services"
	"crdx.org/lighthouse/db"
	"crdx.org/lighthouse/pkg/util/netutil"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/samber/lo"
	"github.com/samber/lo/mutable"
)

type Pinger struct {
	logger *slog.Logger
}

func New() *Pinger {
	return &Pinger{}
}

func (self *Pinger) Init(args *services.Args) error {
	self.logger = args.Logger
	return nil
}

func (self *Pinger) Run() error {
	iface, ipNet, err := netutil.FindNetwork()
	if err != nil {
		return err
	}

	handle := lo.Must(pcap.OpenLive(iface.Name, 65536, true, pcap.BlockForever))
	defer handle.Close()

	devices := db.FindPingableDevices()
	mutable.Shuffle(devices)

	var adapters []*db.Adapter

	for _, device := range devices {
		for _, adapter := range device.Adapters() {
			if db.IsOnline(device, adapter) && db.IsNotResponding(device, adapter) {
				adapters = append(adapters, adapter)
			}
		}
	}

	for _, adapter := range adapters {
		lo.Must0(write(
			handle,
			ipNet.IP,
			net.ParseIP(adapter.IPAddress),
			iface.HardwareAddr,
			lo.Must(net.ParseMAC(adapter.MACAddress)),
		))

		self.logger.Info(
			"sent ICMP request",
			"device", adapter.Device().DisplayName(),
			"mac", adapter.MACAddress,
			"ip", adapter.IPAddress,
		)
	}

	return nil
}

func write(handle *pcap.Handle, srcIP, dstIP net.IP, srcMAC, dstMAC net.HardwareAddr) error {
	ethernetLayer := layers.Ethernet{
		SrcMAC:       srcMAC,
		DstMAC:       dstMAC,
		EthernetType: layers.EthernetTypeIPv4,
	}

	ipLayer := layers.IPv4{
		SrcIP:    srcIP,
		DstIP:    dstIP,
		Protocol: layers.IPProtocolICMPv4,
		Version:  4,
	}

	icmpLayer := layers.ICMPv4{
		TypeCode: layers.CreateICMPv4TypeCode(layers.ICMPv4TypeEchoRequest, 0),
		Id:       1337,
		Seq:      1,
	}

	options := gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}

	buffer := gopacket.NewSerializeBuffer()
	if err := gopacket.SerializeLayers(buffer, options, &ethernetLayer, &ipLayer, &icmpLayer, gopacket.Payload([]byte("hello"))); err != nil {
		return err
	}

	return handle.WritePacketData(buffer.Bytes())
}
