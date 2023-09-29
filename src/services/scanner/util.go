package scanner

import (
	"encoding/binary"
	"net"

	"github.com/samber/lo"
)

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

// expandIPNet takes a *net.IPNet and returns a slice of all net.IPs that belong to that network.
func expandIPNet(ipNet *net.IPNet) []net.IP {
	var ips []net.IP

	n := binary.BigEndian.Uint32([]byte(ipNet.IP))
	mask := binary.BigEndian.Uint32([]byte(ipNet.Mask))
	ip := n & mask
	broadcast := ip | ^mask

	for ip++; ip < broadcast; ip++ {
		var buffer [4]byte
		binary.BigEndian.PutUint32(buffer[:], ip)
		ips = append(ips, net.IP(buffer[:]))
	}

	return ips
}

func ipNetTooLarge(ipNet *net.IPNet) bool {
	return ipNet.Mask[0] != 0xff || ipNet.Mask[1] != 0xff
}
