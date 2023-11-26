package probe

import (
	"errors"
	"io"
	"net"
	"slices"
	"time"

	"crdx.org/lighthouse/pkg/util/netutil"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/samber/lo"
)

type Scanner struct {
	macAddress string
	ipAddress  string
	timeout    time.Duration
}

type Config struct {
	srcIP         net.IP
	dstIP         net.IP
	srcPort       int
	dstPorts      []uint
	srcMACAddress net.HardwareAddr
	dstMACAddress net.HardwareAddr
}

type Result struct {
	Ports []uint
}

func Scan(macAddress string, ipAddress string) (*Result, error) {
	scanner := &Scanner{
		macAddress: macAddress,
		ipAddress:  ipAddress,
		timeout:    10 * time.Second,
	}
	return scanner.run()
}

func (self *Scanner) run() (*Result, error) {
	iface, ipNet, err := netutil.FindNetwork()
	if err != nil {
		return nil, err
	}

	handle, err := pcap.OpenLive(iface.Name, 65535, true, pcap.BlockForever)
	if err != nil {
		return nil, err
	}
	defer handle.Close()

	config := Config{
		dstMACAddress: lo.Must(net.ParseMAC(self.macAddress)),
		srcMACAddress: iface.HardwareAddr,
		dstIP:         net.ParseIP(self.ipAddress),
		srcIP:         ipNet.IP,
		srcPort:       lo.Must(getSourcePort()),
		dstPorts:      Ports(),
	}

	done := make(chan struct{})
	var result Result

	go readPackets(handle, config.srcPort, &result, done)
	sendPackets(handle, &config)
	sendPackets(handle, &config)

	defer time.AfterFunc(self.timeout, func() { handle.Close() }).Stop()

	<-done

	result.Ports = lo.Uniq(result.Ports)
	slices.Sort(result.Ports)

	return &result, nil
}

func sendPackets(handle *pcap.Handle, config *Config) {
	ethernetLayer := layers.Ethernet{
		SrcMAC:       config.srcMACAddress,
		DstMAC:       config.dstMACAddress,
		EthernetType: layers.EthernetTypeIPv4,
	}

	ipLayer := layers.IPv4{
		SrcIP:    config.srcIP,
		DstIP:    config.dstIP,
		Version:  4,
		TTL:      255,
		Protocol: layers.IPProtocolTCP,
	}

	tcpLayer := layers.TCP{
		SrcPort: layers.TCPPort(config.srcPort),
		DstPort: 0,
		SYN:     true,
	}

	lo.Must0(tcpLayer.SetNetworkLayerForChecksum(&ipLayer))

	for _, port := range config.dstPorts {
		tcpLayer.DstPort = layers.TCPPort(port)
		_ = sendPacket(handle, &ethernetLayer, &ipLayer, &tcpLayer)
	}
}

func readPackets(handle *pcap.Handle, srcPort int, result *Result, done chan struct{}) {
	ethernetLayer := &layers.Ethernet{}
	ipLayer := &layers.IPv4{}
	tcpLayer := &layers.TCP{}

	parser := gopacket.NewDecodingLayerParser(layers.LayerTypeEthernet, ethernetLayer, ipLayer, tcpLayer)

	for {
		data, _, err := handle.ReadPacketData()

		if errors.Is(err, pcap.NextErrorTimeoutExpired) || errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			continue
		}

		decodedLayers := []gopacket.LayerType{}
		if err := parser.DecodeLayers(data, &decodedLayers); err != nil {
			continue
		}

		for _, layerType := range decodedLayers {
			if layerType == layers.LayerTypeTCP {
				if tcpLayer.DstPort != layers.TCPPort(srcPort) {
					continue
				} else if tcpLayer.SYN && tcpLayer.ACK {
					result.Ports = append(result.Ports, uint(tcpLayer.SrcPort))
				}
			}
		}
	}

	close(done)
}

func sendPacket(handle *pcap.Handle, layers ...gopacket.SerializableLayer) error {
	buf := gopacket.NewSerializeBuffer()
	if err := gopacket.SerializeLayers(buf, gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}, layers...); err != nil {
		return err
	}
	return handle.WritePacketData(buf.Bytes())
}

func getSourcePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}

	defer listener.Close()
	return listener.Addr().(*net.TCPAddr).Port, nil
}
