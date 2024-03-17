package reader

import (
	"log/slog"

	"crdx.org/lighthouse/cmd/lighthouse/services"
	"crdx.org/lighthouse/pkg/cache"
	"crdx.org/lighthouse/pkg/util/netutil"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/samber/lo"
)

type Reader struct {
	logger          *slog.Logger
	macAddressCache *cache.TemporalCache[string]

	// Since the DHCP handshake happens before the device officially joins the network and starts
	// responding to ARP requests, hostnameCache stores a mapping of MAC addresses to hostnames.
	// This is then used as a lookup for a device hostname when we discover a new device on the
	// network.
	hostnameCache map[string]string
}

func New() *Reader {
	return &Reader{
		macAddressCache: cache.NewTemporal[string](),
		hostnameCache:   map[string]string{},
	}
}

func (self *Reader) Init(args *services.Args) error {
	self.logger = args.Logger
	return nil
}

func (self *Reader) Run() error {
	iface, ipNet, err := netutil.FindNetwork()
	if err != nil {
		return err
	}

	handle := lo.Must(pcap.OpenLive(iface.Name, 65536, true, pcap.BlockForever))
	defer handle.Close()

	packets := gopacket.NewPacketSource(handle, layers.LayerTypeEthernet).Packets()

	for {
		packet := <-packets

		if layer := packet.Layer(layers.LayerTypeDHCPv4); layer != nil {
			self.handleDHCPPacket(layer.(*layers.DHCPv4))
		} else if layer := packet.Layer(layers.LayerTypeARP); layer != nil {
			self.handleARPPacket(layer.(*layers.ARP), ipNet.IP.String())
		} else if layer := packet.Layer(layers.LayerTypeICMPv4); layer != nil {
			self.handleICMPPacket(layer.(*layers.ICMPv4), packet, ipNet.IP.String())
		}
	}
}
