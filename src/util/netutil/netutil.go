package netutil

import (
	"encoding/binary"
	"net"
	"strings"

	"github.com/google/gopacket/macs"
	"github.com/samber/lo"
)

// GetVendor returns the vendor for a mac address, and true if it was found in the gopacket
// database.
func GetVendor(macAddress string) (string, bool) {
	hardwareAddr, err := net.ParseMAC(macAddress)
	if err != nil {
		return "", false
	}

	var prefix [3]byte
	copy(prefix[:], hardwareAddr[:3])

	vendor, found := macs.ValidMACPrefixMap[prefix]
	return vendor, found
}

// UnqualifyHostname returns a hostname without the domain portion and removes any trailing periods,
// if any.
func UnqualifyHostname(hostname string) string {
	hostname = strings.TrimSuffix(hostname, ".")

	index := strings.LastIndex(hostname, ".")
	if index == -1 {
		return hostname
	}

	return hostname[:index]
}

// FindInterface returns the first interface that is considered "up" and is not the loopback
// interface.
func FindInterface() (*net.Interface, bool) {
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

// FindIPNet returns the first IPv4 network for an interface.
func FindIPNet(iface *net.Interface) (*net.IPNet, bool) {
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

// ExpandIPNet takes a network and returns a slice of all IPs that belong to that network.
func ExpandIPNet(ipNet *net.IPNet) []net.IP {
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

// IPNetTooLarge returns true if a network is considered too large to scan.
//
// Anything larger than a /16 is considered too large.
func IPNetTooLarge(ipNet *net.IPNet) bool {
	return ipNet.Mask[0] != 0xff || ipNet.Mask[1] != 0xff
}
