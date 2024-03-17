package writer

import (
	"errors"
	"log/slog"
	"net"
	"time"

	"crdx.org/lighthouse/cmd/lighthouse/services"
	"crdx.org/lighthouse/pkg/util/netutil"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/samber/lo"
)

// Interval between ARP request packets. 50ms means on a /24 a full scan takes about 10 seconds.
const arpPacketInterval = 50 * time.Millisecond

type Writer struct {
	logger *slog.Logger
}

func New() *Writer {
	return &Writer{}
}

func (self *Writer) Init(args *services.Args) error {
	self.logger = args.Logger
	return nil
}

func (self *Writer) Run() error {
	iface, ipNet, err := netutil.FindNetwork()
	if err != nil {
		return err
	}

	handle := lo.Must(pcap.OpenLive(iface.Name, 65536, true, pcap.BlockForever))
	defer handle.Close()

	lo.Must0(self.write(handle, iface, ipNet))
	return nil
}

func (*Writer) write(handle *pcap.Handle, iface *net.Interface, ipNet *net.IPNet) error {
	scan := func() error {
		for _, ip := range netutil.ExpandIPNet(ipNet) {
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
				DstProtAddress:    []byte(ip),
			}

			options := gopacket.SerializeOptions{
				FixLengths:       true,
				ComputeChecksums: true,
			}

			buffer := gopacket.NewSerializeBuffer()
			if err := gopacket.SerializeLayers(buffer, options, &ethernetLayer, &arpLayer); err != nil {
				return err
			}

			if err := handle.WritePacketData(buffer.Bytes()); err != nil {
				return err
			}

			time.Sleep(arpPacketInterval)
		}

		return nil
	}

	return errors.Join(
		scan(),
		scan(),
	)
}
