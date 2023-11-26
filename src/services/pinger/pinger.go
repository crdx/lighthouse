package pinger

import (
	"log/slog"
	"net"
	"time"

	"crdx.org/lighthouse/m"
	"crdx.org/lighthouse/m/repo/deviceR"
	"crdx.org/lighthouse/m/repo/settingR"
	"crdx.org/lighthouse/pkg/util/netutil"
	"crdx.org/lighthouse/services"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/samber/lo"
)

type Pinger struct {
	log *slog.Logger
}

func New() *Pinger {
	return &Pinger{}
}

func (self *Pinger) Init(args *services.Args) error {
	self.log = args.Logger
	return nil
}

func (self *Pinger) Run() error {
	iface, ipNet, err := netutil.FindNetwork()
	if err != nil {
		return err
	}

	handle := lo.Must(pcap.OpenLive(iface.Name, 65536, true, pcap.BlockForever))
	defer handle.Close()

	for {
		if settingR.EnableDeviceScan() {
			devices := deviceR.Pingable()
			lo.Shuffle(devices)

			var adapters []*m.Adapter

			for _, device := range devices {
				for _, adapter := range device.Adapters() {
					if adapter.IsOnline() && adapter.IsNotResponding() {
						adapters = append(adapters, adapter)
					}
				}
			}

			if len(adapters) == 0 {
				return nil
			}

			for _, adapter := range adapters {
				lo.Must0(write(
					handle,
					ipNet.IP,
					net.ParseIP(adapter.IPAddress),
					iface.HardwareAddr,
					lo.Must(net.ParseMAC(adapter.MACAddress)),
				))

				self.log.Info(
					"sent ICMP request",
					"device", adapter.Device().DisplayName(),
					"mac", adapter.MACAddress,
					"ip", adapter.IPAddress,
				)
			}
		}

		time.Sleep(settingR.DeviceScanInterval())
	}
}

func write(handle *pcap.Handle, srcIPAddress, dstIPAddress net.IP, srcMACAddress, dstMACAddress net.HardwareAddr) error {
	ethernetLayer := layers.Ethernet{
		SrcMAC:       srcMACAddress,
		DstMAC:       dstMACAddress,
		EthernetType: layers.EthernetTypeIPv4,
	}

	ipLayer := layers.IPv4{
		SrcIP:    srcIPAddress,
		DstIP:    dstIPAddress,
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
