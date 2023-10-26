package netutil

import (
	"encoding/binary"
	"net"
	"strings"

	"github.com/google/gopacket/macs"
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
